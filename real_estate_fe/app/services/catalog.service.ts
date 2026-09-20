/**
 * CatalogService — các API danh mục dùng chung: menu, tỉnh/thành, phường/xã,
 * dự án, loại bất động sản. Dùng qua composable useCatalogService().
 */
import { BaseService, defineService } from '~/services/api'
import type { QueryParams } from '~/services/api'
import type {
  CityOption,
  OptionTypeRealestate,
  ProjectOption,
  WardOption,
} from '~/types/real_estate'
import type { MenuSettings } from '~/types/window'

export class CatalogService extends BaseService {
  /** Cấu hình menu + user_id của người dùng hiện tại */
  getMenu(): Promise<MenuSettings> {
    return this.getData<MenuSettings>('/category')
  }

  /** Danh sách tỉnh/thành phố */
  getCities(): Promise<CityOption[]> {
    return this.getData<CityOption[]>('/real-estate/list/city')
  }

  /** Danh sách phường/xã theo mã tỉnh (backend nhận param `code`) */
  getWards(provinceCode: string): Promise<WardOption[]> {
    return this.getData<WardOption[]>('/real-estate/list/ward', { code: provinceCode })
  }

  /** Danh sách dự án lọc theo tỉnh/phường (backend nhận `province` và `ward`) */
  getProjects(provinceCode?: string, wardCode?: string): Promise<ProjectOption[]> {
    // Chỉ gửi param có giá trị — giữ đúng cách dựng params của code cũ
    const params: QueryParams = {}
    if (provinceCode) {
      params.province = provinceCode
    }
    if (wardCode) {
      params.ward = wardCode
    }
    return this.getData<ProjectOption[]>('/real-estate/list/project', params)
  }

  /** Danh sách loại bất động sản (bán / cho thuê) */
  getRealEstateTypes(): Promise<OptionTypeRealestate[]> {
    return this.getData<OptionTypeRealestate[]>('/real-estate/list/types')
  }
}

export const useCatalogService = defineService(CatalogService)
