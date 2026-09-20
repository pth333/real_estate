/**
 * ProjectService — toàn bộ API của domain dự án bất động sản.
 * Dùng qua composable useProjectService().
 */
import { BaseService, defineService } from '~/services/api'
import type { ProjectDetail, ProjectSummary } from '~/types/project'
import type { RealEstateResponse } from '~/types/real_estate'

export class ProjectService extends BaseService {
  /** Danh sách dự án theo slug danh mục dự án */
  getByCategorySlug(categorySlug: string): Promise<ProjectSummary[]> {
    return this.getData<ProjectSummary[]>(`/real-estate/project-category/${categorySlug}`)
  }

  /** Danh sách dự án nổi bật (backend nhận param `limit`) */
  getFeatured(limit: number): Promise<ProjectSummary[]> {
    return this.getData<ProjectSummary[]>('/real-estate/project/featured', { limit })
  }

  /** Tin đăng đang bán thuộc 1 dự án */
  getListings(projectId: number): Promise<RealEstateResponse[]> {
    return this.getData<RealEstateResponse[]>(`/real-estate/project/${projectId}/listings`)
  }

  /** Chi tiết dự án theo id (hoặc slug) để hiển thị trang chi tiết */
  getDetail(slugOrId: string | number): Promise<ProjectDetail> {
    return this.getData<ProjectDetail>(`/real-estate/project/detail/${slugOrId}`)
  }

  /** Tăng lượt xem dự án (endpoint không cần đọc dữ liệu trả về) */
  async incrementView(id: number): Promise<void> {
    await this.post(`/real-estate/project/view/${id}`)
  }
}

export const useProjectService = defineService(ProjectService)
