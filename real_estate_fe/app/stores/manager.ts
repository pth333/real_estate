/**
 * Store quản lý (manager) — cache dữ liệu các danh mục (bài viết, dự án,
 * yêu thích) để khi navigate giữa các trang con không phải gọi lại API.
 *
 * Mỗi danh mục giữ: data + total + page/size/search hiện tại + cờ loaded.
 * fetchXxx chỉ gọi API khi chưa có dữ liệu hoặc tham số (page/search) thay đổi.
 * Sau nghiệp vụ thêm/sửa/xóa → gọi invalidateXxx() để lần sau fetch lại.
 */

import type { IManagerPostItem, IManagerPostListResponse, ManagerProject } from '~/types/manager'
import type { RealEstateResponse } from '~/types/real_estate'

export const useManagerStore = defineStore('manager', () => {
  const { $api } = useNuxtApp()

  // ── Bài viết ──
  const posts = ref<IManagerPostItem[]>([])
  const postsTotal = ref(0)
  const postsPage = ref(1)
  const postsSize = ref(10)
  const postsSearch = ref('')
  const postsLoaded = ref(false)

  async function fetchPosts(opts: { page: number; size: number; search: string }) {
    // Đã có dữ liệu cho đúng trang + từ khoá hiện tại → không gọi lại API
    if (postsLoaded.value && postsPage.value === opts.page && postsSearch.value === opts.search) {
      return
    }
    postsPage.value = opts.page
    postsSize.value = opts.size
    postsSearch.value = opts.search
    try {
      const res = await $api.get<{ data: IManagerPostListResponse }>('/manager/posts', {
        params: { search: opts.search, page: opts.page, size: opts.size },
      })
      posts.value = res?.data?.posts || []
      postsTotal.value = res?.data?.total || 0
      postsLoaded.value = true
    } catch (error: any) {
      window.message?.error('Lỗi khi tải danh sách bài viết: ' + (error?.message || 'Lỗi máy chủ'))
      throw error
    }
  }

  function invalidatePosts() {
    postsLoaded.value = false
  }

  // ── Dự án ──
  const projects = ref<ManagerProject[]>([])
  const projectsTotal = ref(0)
  const projectsPage = ref(1)
  const projectsSize = ref(10)
  const projectsSearch = ref('')
  const projectsLoaded = ref(false)

  async function fetchProjects(opts: { page: number; size: number; search: string }) {
    if (projectsLoaded.value && projectsPage.value === opts.page && projectsSearch.value === opts.search) {
      return
    }
    projectsPage.value = opts.page
    projectsSize.value = opts.size
    projectsSearch.value = opts.search
    try {
      const res = await $api.get<{ data: ManagerProject[]; total: number }>('/manager/projects', {
        params: { search: opts.search, page: opts.page, size: opts.size },
      })
      projects.value = res?.data || []
      projectsTotal.value = res?.total || 0
      projectsLoaded.value = true
    } catch (error: any) {
      window.message?.error('Lỗi khi tải danh sách dự án: ' + (error?.message || 'Lỗi máy chủ'))
      throw error
    }
  }

  function invalidateProjects() {
    projectsLoaded.value = false
  }

  // ── Yêu thích ──
  const favorites = ref<RealEstateResponse[]>([])
  const favoritesTotal = ref(0)
  const favoritesPage = ref(1)
  const favoritesSize = ref(12)
  const favoritesLoaded = ref(false)

  async function fetchFavorites(opts: { page: number; size: number }) {
    if (favoritesLoaded.value && favoritesPage.value === opts.page) {
      return
    }
    favoritesPage.value = opts.page
    favoritesSize.value = opts.size
    try {
      const res = await $api.get<{ data: RealEstateResponse[]; total: number }>('/real-estate/favorites', {
        params: { page: opts.page, size: opts.size },
      })
      favorites.value = res?.data || []
      favoritesTotal.value = res?.total || 0
      favoritesLoaded.value = true
    } catch (error: any) {
      window.message?.error('Lỗi khi tải danh mục yêu thích: ' + (error?.message || 'Lỗi máy chủ'))
      throw error
    }
  }

  function invalidateFavorites() {
    favoritesLoaded.value = false
  }

  return {
    // bài viết
    posts, postsTotal, postsPage, postsSize, postsSearch, postsLoaded,
    fetchPosts, invalidatePosts,
    // dự án
    projects, projectsTotal, projectsPage, projectsSize, projectsSearch, projectsLoaded,
    fetchProjects, invalidateProjects,
    // yêu thích
    favorites, favoritesTotal, favoritesPage, favoritesSize, favoritesLoaded,
    fetchFavorites, invalidateFavorites,
  }
})
