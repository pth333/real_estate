/**
 * UploadService — toàn bộ API upload media: presign R2, upload ảnh qua backend, confirm.
 * Dùng qua composable useUploadService().
 */
import { BaseService, defineService } from '~/services/api'
import type { ConfirmResponse, PresignResponse } from '~/types/uploadmedia'

export class UploadService extends BaseService {
  /** Xin presigned URL để client PUT file trực tiếp lên R2 */
  async presign(
    filename: string,
    contentType: string,
  ): Promise<{ upload_url: string; key: string; expires_at: string }> {
    const res = await this.post<PresignResponse>('/upload/presign', {
      filename,
      content_type: contentType,
    })

    if (!res.success || !res.data) {
      throw new Error(res.message || 'Không thể lấy presigned URL')
    }

    return res.data
  }

  /**
   * Upload ảnh qua backend (multipart/form-data).
   * @param kind "project" → lưu vào bảng image_projects (ảnh dự án)
   */
  async uploadImage(
    file: File,
    kind?: 'project',
  ): Promise<{ image_id: number; public_url: string; thumbnail_url?: string }> {
    const formData = new FormData()
    formData.append('file', file)
    if (kind) formData.append('kind', kind)

    const res = await this.post<ConfirmResponse>('/upload/image', formData)

    if (!res.success || !res.data) {
      throw new Error(res.message || 'Upload ảnh thất bại')
    }

    return res.data
  }

  /**
   * Xác nhận đã PUT file lên R2 xong, backend ghi nhận vào DB.
   * @param kind "project" → lưu vào bảng image_projects (ảnh dự án)
   */
  async confirm(
    key: string,
    kind?: 'project',
  ): Promise<{ image_id: number; public_url: string; thumbnail_url?: string }> {
    const res = await this.post<ConfirmResponse>('/upload/confirm', {
      key,
      ...(kind ? { kind } : {}),
    })

    if (!res.success || !res.data) {
      throw new Error(res.message || 'Xác nhận upload thất bại')
    }

    return res.data
  }
}

export const useUploadService = defineService(UploadService)
