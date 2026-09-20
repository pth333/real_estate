/**
 * RbacService — toàn bộ API quản trị role / permission / gán role cho user.
 * Dùng qua composable useRbacService().
 */
import { BaseService, defineService } from '~/services/api'
import type { ListResult } from '~/services/api'
import type { AdminUserQuery, AdminUserResponse, PermissionResponse, RoleResponse } from '~/types/rbac'

export class RbacService extends BaseService {
  /** Danh sách role kèm permission của từng role */
  getRoles(): Promise<RoleResponse[]> {
    return this.getData<RoleResponse[]>('/admin/roles')
  }

  /** Danh sách toàn bộ permission của hệ thống */
  getPermissions(): Promise<PermissionResponse[]> {
    return this.getData<PermissionResponse[]>('/admin/permissions')
  }

  /** Danh sách user (tìm theo tên/email/số điện thoại) — trả items + tổng số bản ghi */
  getUsers(params: AdminUserQuery): Promise<ListResult<AdminUserResponse>> {
    return this.getList<AdminUserResponse>('/admin/users', params)
  }

  /** Gán lại toàn bộ role cho 1 user; backend trả 400 nếu role không tồn tại hoặc mảng rỗng */
  setUserRoles(userId: number, roles: string[]): Promise<AdminUserResponse> {
    return this.putData<AdminUserResponse>('/admin/users/' + userId + '/roles', { roles })
  }
}

export const useRbacService = defineService(RbacService)
