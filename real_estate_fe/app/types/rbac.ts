/** Role trong hệ thống RBAC (GET /admin/roles) */
export interface RoleResponse {
  id: number
  /** Mã role: CUSTOMER / BROKER / ADMIN */
  code: string
  name: string
  description: string
  is_active: boolean
  /** Danh sách permission code thuộc role này */
  permissions: string[]
}

/** Permission trong hệ thống RBAC (GET /admin/permissions) */
export interface PermissionResponse {
  id: number
  /** Mã permission, ví dụ: user.read */
  code: string
  name: string
  /** Nhóm chức năng mà permission thuộc về */
  module: string
  description: string
}

/** Người dùng ở góc nhìn quản trị (GET /admin/users) */
export interface AdminUserResponse {
  id: number
  name: string
  email: string
  phone: string
  is_active: boolean
  /** Danh sách role code đang gán cho user */
  roles: string[]
  created_at: string
}

/**
 * Tham số truy vấn danh sách user cho admin.
 * Dùng type alias (không dùng interface) để gán được vào QueryParams của BaseService.
 */
export type AdminUserQuery = {
  search?: string
  page?: number
  size?: number
}
