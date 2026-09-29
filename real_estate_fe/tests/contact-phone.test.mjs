import test from "node:test";
import assert from "node:assert/strict";

const { isVietnamesePhone, normalizeVietnamesePhone } = await import("../app/utils/deposit.ts");

test("Chuẩn hoá SĐT: bỏ khoảng trắng, dấu chấm, gạch nối và đổi +84 thành 0", () => {
  assert.equal(normalizeVietnamesePhone("0901234567"), "0901234567");
  assert.equal(normalizeVietnamesePhone(" 090 123 4567 "), "0901234567");
  assert.equal(normalizeVietnamesePhone("090.123.4567"), "0901234567");
  assert.equal(normalizeVietnamesePhone("090-123-4567"), "0901234567");
  assert.equal(normalizeVietnamesePhone("+84901234567"), "0901234567");
  assert.equal(normalizeVietnamesePhone("(090) 1234567"), "0901234567");
});

test("SĐT hợp lệ: 10 số bắt đầu bằng 0", () => {
  assert.equal(isVietnamesePhone("0901234567"), true);
  assert.equal(isVietnamesePhone("+84 901 234 567"), true);
});

test("SĐT không hợp lệ bị chặn", () => {
  const invalid = ["", "090123456", "09012345678", "19001234567", "09012345a7", "+8490123456", "84 901234567"];
  for (const phone of invalid) {
    assert.equal(isVietnamesePhone(phone), false, `mong đợi SĐT không hợp lệ: ${phone}`);
  }
});
