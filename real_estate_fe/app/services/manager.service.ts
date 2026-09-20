/**
 * ManagerService — toàn bộ API khu vực quản lý: bài viết và dự án.
 * Dùng qua composable useManagerService().
 *
 * Lưu ý shape response (đối chiếu internal/response/response.go + handler backend):
 * - GET /manager/posts        → envelope { success, data: { posts, total } }
 * - GET /manager/projects     → { total, data } ở TOP-LEVEL (không envelope)
 * - GET /manager/projects/:id → envelope { success, data: <chi tiết dự án> }
 * - POST /manager/create-post → envelope { success, message, data: { id } }
 * - PUT /manager/update-post  → envelope { success, message, data }
 * - POST /manager/create-project → envelope { success, message, data: { id } }
 * - PUT /manager/update-project  → envelope { success, message, data }
 *
 * Vì vậy create/update post + project trả NGUYÊN envelope (không bóc `.data`)
 * để trang gọi giữ được điều kiện `if (res.success)` như trước.
 */
import { BaseService, defineService } from '~/services/api'
import type { ApiEnvelope, ListResult } from '~/services/api'
import type {
  CreatePostPayload,
  CreateProjectPayload,
  IManagerPostListResponse,
  ManagerProject,
  ManagerProjectDetailResponse,
  UpdatePostPayload,
  UpdateProjectPayload,
} from '~/types/manager'
import type { CreatePostResponse, UpdatePostResponse } from '~/types/real_estate'

/** Kết quả POST /manager/create-project — id dự án vừa tạo nằm trong `data.id` */
export type CreateProjectResult = ApiEnvelope<{ id: number }>

/** Kết quả PUT /manager/update-project/:id — backend đặt message vào `data` */
export type UpdateProjectResult = ApiEnvelope<unknown>

export class ManagerService extends BaseService {
  // ── Bài viết ─────────────────────────────────────────

  /** Danh sách bài viết của manager (envelope → bóc `.data` = { posts, total }) */
  getPosts(params: { search: string; page: number; size: number }): Promise<IManagerPostListResponse> {
    return this.getData<IManagerPostListResponse>('/manager/posts', params)
  }

  /** Tạo tin đăng mới */
  createPost(payload: CreatePostPayload): Promise<CreatePostResponse> {
    return this.post<CreatePostResponse>('/manager/create-post', payload)
  }

  /** Cập nhật tin đăng theo id */
  updatePost(id: number, payload: UpdatePostPayload): Promise<UpdatePostResponse> {
    return this.put<UpdatePostResponse>(`/manager/update-post/${id}`, payload)
  }

  /** Xóa tin đăng theo id */
  async deletePost(id: number): Promise<void> {
    await this.remove(`/manager/posts/${id}`)
  }

  // ── Dự án ────────────────────────────────────────────

  /**
   * Danh sách dự án của manager.
   * Endpoint trả `{ total, data }` top-level (KHÔNG envelope) → map sang { items, total }.
   */
  async getProjects(params: {
    search: string
    page: number
    size: number
  }): Promise<ListResult<ManagerProject>> {
    const res = await this.get<{ data: ManagerProject[]; total: number }>('/manager/projects', params)
    return { items: res.data ?? [], total: res.total ?? 0 }
  }

  /** Chi tiết 1 dự án để điền form chỉnh sửa (envelope → bóc `.data`) */
  getProject(id: number): Promise<ManagerProjectDetailResponse> {
    return this.getData<ManagerProjectDetailResponse>(`/manager/projects/${id}`)
  }

  /** Tạo dự án mới — trả nguyên envelope để trang đọc `res.success` */
  createProject(payload: CreateProjectPayload): Promise<CreateProjectResult> {
    return this.post<CreateProjectResult>('/manager/create-project', payload)
  }

  /** Cập nhật dự án theo id — trả nguyên envelope để trang đọc `res.success` */
  updateProject(id: number, payload: UpdateProjectPayload): Promise<UpdateProjectResult> {
    return this.put<UpdateProjectResult>(`/manager/update-project/${id}`, payload)
  }
}

export const useManagerService = defineService(ManagerService)
