import api from './service.api'

export const menuApi = {
  GetAllCategory() {
    return api.get('/category')
  },
}
