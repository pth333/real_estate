import type { FileItem, PresignResponse, ConfirmResponse } from '~/types/uploadmedia'

/**
 * Tạo presigned URL từ backend
 */
export async function getPresignedUrl(
  filename: string,
  contentType: string,
): Promise<{ upload_url: string; key: string; expires_at: string }> {
  const { $api } = useNuxtApp()
  const res = await $api.post<PresignResponse>('/upload/presign', {
    filename,
    content_type: contentType,
  })

  if (!res.success || !res.data) {
    throw new Error(res.message || 'Không thể lấy presigned URL')
  }

  return res.data
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
    const xhr = new XMLHttpRequest()
    xhr.open('PUT', url, true)
    xhr.setRequestHeader('Content-Type', file.type)

    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable && onProgress) {
        onProgress(Math.round((e.loaded / e.total) * 100))
      }
    }

    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        onProgress?.(100)
        resolve()
      } else {
        reject(new Error(`Upload R2 thất bại: ${xhr.status}`))
      }
    }

    xhr.onerror = () => reject(new Error('Lỗi mạng khi upload lên R2'))
    xhr.ontimeout = () => reject(new Error('Upload lên R2 đã hết thời gian'))
    xhr.send(file)
  })
}

/**
 * Xác nhận upload với backend
 */
export async function confirmUpload(
  key: string,
): Promise<{ image_id: number; public_url: string; thumbnail_url?: string }> {
  const { $api } = useNuxtApp()
  const res = await $api.post<ConfirmResponse>('/upload/confirm', { key })

  if (!res.success || !res.data) {
    throw new Error(res.message || 'Xác nhận upload thất bại')
  }

  return res.data
}

/**
 * Upload hoàn chỉnh: presign → R2 → confirm, cập nhật status & progress vào item
 */
export async function uploadFile(item: FileItem): Promise<void> {
  try {
    // Step 1: Lấy presigned URL
    item.status = 'gettingPresign'
    const { upload_url, key, expires_at } = await getPresignedUrl(
      item.file.name,
      item.file.type,
    )

    Object.assign(item, {
      key,
      uploadUrl: upload_url,
      expiresAt: new Date(expires_at).getTime(),
    })

    // Step 4: Upload trực tiếp lên R2
    item.status = 'uploading'
    await uploadToR2(item.file, upload_url, (pct) => {
      item.progress = pct
    })

    // Step 5: Confirm với backend
    item.status = 'confirming'
    const result = await confirmUpload(key)

    Object.assign(item, {
      imageId: result.image_id,
      publicUrl: result.public_url,
      thumbnailUrl: result.thumbnail_url,
      status: 'done',
      progress: 100,
    })
  } catch (err: unknown) {
    item.status = 'error'
    item.errorMessage = err instanceof Error ? err.message : 'Upload thất bại'
    console.error('Upload error:', err)
  }
}
