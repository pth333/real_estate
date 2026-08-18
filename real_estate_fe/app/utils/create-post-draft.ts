import type { IInformationRealestate } from "~/types/real_estate";

// Key lưu bản nháp tin đăng trong localStorage
export const CREATE_POST_DRAFT_KEY = "create_post_draft";

/**
 * Lưu bản nháp tin đăng xuống localStorage.
 * Chỉ lưu thông tin text của form — images (id + public_url + thumbnail_url)
 * đã nằm trong form nên khôi phục được preview ảnh/video từ URL S3.
 * Không lưu File binary.
 */
export function saveDraft(form: IInformationRealestate): void {
  if (!import.meta.client) return;
  try {
    localStorage.setItem(
      CREATE_POST_DRAFT_KEY,
      JSON.stringify({ savedAt: Date.now(), form }),
    );
  } catch {
    // localStorage đầy hoặc bị chặn → bỏ qua, không làm crash app
  }
}

/**
 * Đọc bản nháp từ localStorage. Trả null nếu không có hoặc dữ liệu lỗi.
 */
export function loadDraft(): IInformationRealestate | null {
  if (!import.meta.client) return null;
  try {
    const raw = localStorage.getItem(CREATE_POST_DRAFT_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw);
    if (!parsed?.form) return null;
    return parsed.form as IInformationRealestate;
  } catch {
    return null;
  }
}

/**
 * Xoá bản nháp khỏi localStorage.
 */
export function clearDraft(): void {
  if (!import.meta.client) return;
  try {
    localStorage.removeItem(CREATE_POST_DRAFT_KEY);
  } catch {
    // bỏ qua nếu không xoá được
  }
}
