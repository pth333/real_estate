import test from "node:test";
import assert from "node:assert/strict";

// Mock localStorage (node --test không có window/localStorage)
function createMockStorage() {
  const store = new Map();
  return {
    getItem: (key) => (store.has(key) ? store.get(key) : null),
    setItem: (key, value) => void store.set(key, String(value)),
    removeItem: (key) => void store.delete(key),
    _store: store,
  };
}

globalThis.localStorage = createMockStorage();

// import sau khi mock localStorage
const draftModule = await import("../app/utils/create-post-draft.ts");
const { saveDraft, loadDraft, clearDraft, CREATE_POST_DRAFT_KEY } = draftModule;

const sampleForm = {
  listingType: "sell",
  province: "79",
  ward: "76049",
  detail_address: "123 Lê Lợi",
  real_estate_type: "12-CanHo",
  area: 80,
  price: 100_000_000,
  unit: "vnd",
  legal_docs: "Sổ hồng",
  interior: "Cao cấp",
  bathroom_count: 2,
  bedroom_count: 3,
  house_direction: "Đông",
  balcony_direction: "Nam",
  move_in_time: "Ngay",
  price_electricity: 3500,
  price_water: 20000,
  price_internet: 200000,
  amenities: ["Hồ bơi", "Thang máy"],
  image_ids: [1, 2, 3],
  images: [
    {
      imageId: 1,
      publicUrl: "https://pub-5eb4e976c2fe4062ba3cdabce48568cc.r2.dev/uploads/6e7d20302f54d99b0caf51c9822e08c3.png",
      thumbnailUrl: "https://pub-5eb4e976c2fe4062ba3cdabce48568cc.r2.dev/uploads/6e7d20302f54d99b0caf51c9822e08c3-thumb.jpg",
      fileType: "image",
    },
    {
      imageId: 2,
      publicUrl: "https://pub-5eb4e976c2fe4062ba3cdabce48568cc.r2.dev/uploads/b3362dc90cd41e89fb181c6c18b5857c.mp4",
      thumbnailUrl: "https://pub-5eb4e976c2fe4062ba3cdabce48568cc.r2.dev/uploads/b3362dc90cd41e89fb181c6c18b5857c-thumb.jpg",
      fileType: "video",
    },
  ],
  contact_name: "Nguyễn Văn A",
  contact_email: "a@example.com",
  contact_phone: "0901234567",
  title: "Bán căn hộ chung cư đẹp tại Quận 1",
  description: "Căn hộ 2PN 2WC, view đẹp, nội thất cao cấp, sổ hồng riêng.",
  tab: "information",
};

test("saveDraft lưu form vào localStorage dưới đúng key", () => {
  clearDraft();
  saveDraft(sampleForm);
  const raw = localStorage.getItem(CREATE_POST_DRAFT_KEY);
  assert.ok(raw, "draft phải được lưu vào localStorage");
  const parsed = JSON.parse(raw);
  assert.equal(parsed.form.title, sampleForm.title);
  assert.equal(parsed.form.province, "79");
  assert.deepEqual(parsed.form.images, sampleForm.images);
  assert.ok(parsed.savedAt, "phải lưu kèm mốc thời gian savedAt");
});

test("loadDraft trả về form đúng dữ liệu đã lưu (gồm images với public_url)", () => {
  clearDraft();
  saveDraft(sampleForm);
  const loaded = loadDraft();
  assert.ok(loaded, "loadDraft phải trả về form");
  assert.equal(loaded.title, sampleForm.title);
  assert.deepEqual(loaded.amenities, ["Hồ bơi", "Thang máy"]);
  assert.equal(loaded.images.length, 2);
  assert.equal(loaded.images[1].fileType, "video");
  assert.equal(
    loaded.images[1].publicUrl,
    "https://pub-5eb4e976c2fe4062ba3cdabce48568cc.r2.dev/uploads/b3362dc90cd41e89fb181c6c18b5857c.mp4",
  );
  assert.equal(loaded.tab, "information");
});

test("loadDraft trả null khi chưa lưu gì", () => {
  clearDraft();
  assert.equal(loadDraft(), null);
});

test("loadDraft trả null khi dữ liệu trong localStorage bị hỏng", () => {
  localStorage.setItem(CREATE_POST_DRAFT_KEY, "{{{ không phải json");
  assert.equal(loadDraft(), null);
});

test("loadDraft trả null khi thiếu trường form", () => {
  localStorage.setItem(CREATE_POST_DRAFT_KEY, JSON.stringify({ savedAt: 1 }));
  assert.equal(loadDraft(), null);
});

test("clearDraft xoá hẳn bản nháp", () => {
  saveDraft(sampleForm);
  assert.ok(loadDraft(), "có draft trước khi clear");
  clearDraft();
  assert.equal(loadDraft(), null);
});

test("saveDraft không crash khi localStorage bị chặn", () => {
  globalThis.localStorage.setItem = () => {
    throw new Error("QuotaExceededError");
  };
  saveDraft(sampleForm); // không được throw
  assert.ok(true);
});
