import type { FileItem } from "~/types/uploadmedia";
import { useUploadService } from "~/services/upload.service";
interface ValidationResult {
  valid: boolean;
  message: string;
}
export function useValidate() {
  async function validateImage(file: File): Promise<ValidationResult> {
    const allowedFormats = [
      "image/png",
      "image/jpeg",
      "image/jpg",
      "image/gif",
    ];
    if (!allowedFormats.includes(file.type)) {
      return {
        valid: false,
        message: `Định dạng không hợp lệ. Chỉ hỗ trợ PNG, JPG, JPEG, GIF`,
      };
    }
    const maxSize = 10 * 1024 * 1024; // 10MB
    if (file.size > maxSize) {
      return { valid: false, message: `Dung lượng tối đa 10MB` };
    }

    const imageUrl = URL.createObjectURL(file);
    try {
      const dimensions = await new Promise<{ width: number; height: number }>(
        (resolve, reject) => {
          const image = new Image();
          image.onload = () =>
            resolve({
              width: image.naturalWidth,
              height: image.naturalHeight,
            });
          image.onerror = () => reject(new Error("Không thể đọc ảnh"));
          image.src = imageUrl;
        },
      );

      if (dimensions.width < 1400 || dimensions.height < 1050) {
        return {
          valid: false,
          message: `Kích thước tối thiểu 1600x1200px, ảnh hiện tại ${dimensions.width}x${dimensions.height}px`,
        };
      }

      const ratio = dimensions.width / dimensions.height;
      const isSquare = dimensions.width === dimensions.height;
      const isFourByThree = Math.abs(ratio - 4 / 3) <= 0.02;
      console.log(dimensions.width, dimensions.height);
      console.log(isFourByThree);
      if (!isSquare && !isFourByThree) {
        return {
          valid: false,
          message: `Tỷ lệ ảnh chỉ được là 4:3 hoặc hình vuông, ảnh hiện tại ${dimensions.width}x${dimensions.height}px`,
        };
      }

      return { valid: true, message: "" };
    } catch {
      return { valid: false, message: "File không phải ảnh hợp lệ" };
    } finally {
      URL.revokeObjectURL(imageUrl);
    }
  }

  function validateVideo(file: File): ValidationResult {
    const allowedFormats = ["video/mp4", "video/quicktime"];
    if (!allowedFormats.includes(file.type)) {
      return {
        valid: false,
        message: `Định dạng không hợp lệ. Chỉ hỗ trợ MP4, MOV`,
      };
    }
    const maxSize = 200 * 1024 * 1024; // 200MB
    if (file.size > maxSize) {
      return { valid: false, message: `Dung lượng tối đa 200MB` };
    }
    return { valid: true, message: "" };
  }

  return { validateImage, validateVideo };
}

/**
 * Tạo presigned URL từ backend
 */
export async function getPresignedUrl(
  filename: string,
  contentType: string,
): Promise<{ upload_url: string; key: string; expires_at: string }> {
  const uploadService = useUploadService();
  return uploadService.presign(filename, contentType);
}

export async function uploadImageToBackend(
  file: File,
  kind: "project" | undefined,
  onProgress?: (pct: number) => void,
): Promise<{ image_id: number; public_url: string; thumbnail_url?: string }> {
  const uploadService = useUploadService();
  onProgress?.(10);

  const data = await uploadService.uploadImage(file, kind);
  onProgress?.(100);

  return data;
}

/**
 * Upload file lên R2 bằng presigned URL (XHR, có progress)
 */
export function uploadToR2(
  file: File,
  url: string,
  onProgress?: (pct: number) => void,
): Promise<void> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open("PUT", url, true);
    xhr.setRequestHeader("Content-Type", file.type);

    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable && onProgress) {
        onProgress(Math.round((e.loaded / e.total) * 100));
      }
    };

    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        onProgress?.(100);
        resolve();
      } else {
        reject(new Error(`Upload R2 thất bại: ${xhr.status}`));
      }
    };

    xhr.onerror = () => reject(new Error("Lỗi mạng khi upload lên R2"));
    xhr.ontimeout = () => reject(new Error("Upload lên R2 đã hết thời gian"));
    xhr.send(file);
  });
}

/**
 * Xác nhận upload với backend
 * @param kind "project" → lưu vào bảng image_projects (ảnh dự án)
 */
export async function confirmUpload(
  key: string,
  kind?: "project",
): Promise<{ image_id: number; public_url: string; thumbnail_url?: string }> {
  const uploadService = useUploadService();
  return uploadService.confirm(key, kind);
}

/**
 * Upload hoàn chỉnh: presign → R2 → confirm, cập nhật status & progress vào item
 * @param kind "project" → confirm lưu vào bảng image_projects
 */
export async function uploadFile(
  item: FileItem,
  kind?: "project",
): Promise<void> {
  try {
    if (item.fileType === "image") {
      item.status = "uploading";
      const result = await uploadImageToBackend(item.file, kind, (pct) => {
        item.progress = pct;
      });

      Object.assign(item, {
        imageId: result.image_id,
        publicUrl: result.public_url,
        thumbnailUrl: result.thumbnail_url,
        status: "done",
        progress: 100,
      });
      return;
    }

    item.status = "gettingPresign";
    const { upload_url, key, expires_at } = await getPresignedUrl(
      item.file.name,
      item.file.type,
    );

    Object.assign(item, {
      key,
      uploadUrl: upload_url,
      expiresAt: new Date(expires_at).getTime(),
    });

    item.status = "uploading";
    await uploadToR2(item.file, upload_url, (pct) => {
      item.progress = pct;
    });

    // Step 5: Confirm với backend
    item.status = "confirming";
    const result = await confirmUpload(key, kind);

    Object.assign(item, {
      imageId: result.image_id,
      publicUrl: result.public_url,
      thumbnailUrl: result.thumbnail_url,
      status: "done",
      progress: 100,
    });
  } catch (err: unknown) {
    item.status = "error";
    item.errorMessage = err instanceof Error ? err.message : "Upload thất bại";
    console.error("Upload error:", err);
  }
}
