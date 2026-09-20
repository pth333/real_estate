/**
 * AiService — API gọi AI sinh tiêu đề & mô tả cho tin đăng.
 * Dùng qua composable useAiService().
 */
import { BaseService, defineService } from '~/services/api'

/** Nội dung AI trả về */
export interface AIContent {
  title: string
  description: string
}

/** Payload gửi lên AI: thông tin BĐS đang nhập dở + văn phong mong muốn */
export interface AIContentPayload {
  /** Văn phong: "lich_su" hoặc "tre_trung" */
  tone: 'lich_su' | 'tre_trung'
  /** Loại tin: "sell" hoặc "rent" */
  listing_type: 'sell' | 'rent'
  real_estate_type?: string
  province?: string
  ward?: string
  address?: string
  area?: number
  price_per_m2?: number
  unit?: string
  bedrooms?: number
  bathrooms?: number
  legal_docs?: string
  interior?: string
  house_direction?: string
  balcony_direction?: string
  contact_name?: string
  contact_phone?: string
}

export class AiService extends BaseService {
  /** Sinh tiêu đề + mô tả dựa trên thông tin BĐS */
  generateContent(payload: AIContentPayload): Promise<AIContent> {
    return this.postData<AIContent>('/ai/generate-content', payload)
  }
}

export const useAiService = defineService(AiService)
