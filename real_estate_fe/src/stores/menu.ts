import { defineStore } from 'pinia'
import { menuApi } from '@/api/general.api'

export const useMenuStore = defineStore('menu', () => {
  const menuItems = ref<string[]>([])

  const fetchMenuItems = async () => {
    try {
      const res = await menuApi.GetAllCategory()
      const data = res.data?.data || []
      menuItems.value = data.map((item: any) => item.name)
    } catch (error) {
      console.error('Error fetching menu items:', error)
    }
  }

  return {
    menuItems,
    fetchMenuItems,
  }
})
