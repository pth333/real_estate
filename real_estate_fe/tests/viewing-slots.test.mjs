import test from "node:test";
import assert from "node:assert/strict";

const { buildViewingSlots } = await import("../app/utils/deposit.ts");

// Ngày tương lai: đủ khung 08:00 → 18:00, mỗi khung đúng 1 tiếng
test("Ngày tương lai có đủ khung giờ 1 tiếng trong giờ làm việc", () => {
  const slots = buildViewingSlots("2099-01-01", new Date(2099, 0, 1, 7, 0));

  assert.equal(slots.length, 10);
  assert.equal(slots[0].start, "08:00");
  assert.equal(slots[0].end, "09:00");
  assert.equal(slots[0].label, "08:00 - 09:00");
  assert.equal(slots.at(-1).start, "17:00");
  assert.equal(slots.at(-1).end, "18:00");
});

// Mỗi khung luôn dao động đúng 1 tiếng, không có khung nào lệch
test("Mọi khung đều dao động đúng 1 tiếng và nối tiếp nhau", () => {
  const slots = buildViewingSlots("2099-01-01", new Date(2099, 0, 1, 7, 0));

  for (const slot of slots) {
    const [startHour, startMinute] = slot.start.split(":").map(Number);
    const [endHour, endMinute] = slot.end.split(":").map(Number);
    assert.equal(startMinute, 0);
    assert.equal(endMinute, 0);
    assert.equal(endHour - startHour, 1);
  }
});

// Hôm nay: chỉ còn khung bắt đầu sau giờ hiện tại
test("Ngày hôm nay chỉ còn khung giờ sau thời điểm hiện tại", () => {
  const now = new Date(2099, 0, 1, 9, 30);
  const slots = buildViewingSlots("2099-01-01", now);

  assert.equal(slots[0].start, "10:00");
  assert.equal(slots[0].end, "11:00");
  assert.ok(slots.every((slot) => slot.start > "09:30"));
  assert.ok(!slots.some((slot) => slot.start === "09:00"));
});

// Đúng đầu khung thì khung đó bị loại (khung đã bắt đầu)
test("Đúng mốc giờ bắt đầu thì khung đó không còn chọn được", () => {
  const slots = buildViewingSlots("2099-01-01", new Date(2099, 0, 1, 8, 0));

  assert.equal(slots[0].start, "09:00");
});

// Quá giờ làm việc: hết khung của hôm nay
test("Sau giờ làm việc thì hôm nay không còn khung nào", () => {
  assert.deepEqual(buildViewingSlots("2099-01-01", new Date(2099, 0, 1, 18, 0)), []);
  assert.deepEqual(buildViewingSlots("2099-01-01", new Date(2099, 0, 1, 22, 15)), []);
});

// Chưa chọn ngày thì trả về danh sách đầy đủ (select đang bị disable ở UI)
test("Chưa chọn ngày thì trả về danh sách khung đầy đủ", () => {
  const slots = buildViewingSlots(null, new Date(2099, 0, 1, 9, 30));

  assert.equal(slots.length, 10);
  assert.equal(slots[0].start, "08:00");
});
