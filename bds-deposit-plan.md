# Hệ Thống Đặt Cọc Bất Động Sản — Tài Liệu Thiết Kế

> **Mô hình:** B2C Platform · **Escrow:** Tài khoản công ty · **Thanh toán:** VNPay / Momo / ZaloPay · **Database:** MySQL

---

## 1. Tổng Quan

### Mục tiêu
- **An toàn cho khách:** Tiền cọc không vào tay môi giới ngay — nằm tại escrow của platform cho đến khi điều kiện thỏa mãn
- **Minh bạch 2 chiều:** Cả khách và môi giới đều phải xác nhận bằng OTP tại chỗ — không bên nào có thể gian lận đơn phương
- **Tự động thông báo:** Hệ thống tự gửi mail nhắc lịch trước 24h, tự động xử lý hoàn tiền theo trạng thái
- **Tranh chấp có trọng tài:** Mọi mâu thuẫn đều được freeze tiền và chuyển admin xử lý với đầy đủ bằng chứng 2 bên

### Các bên tham gia

| Bên | Vai trò |
|---|---|
| **Khách hàng** | Tìm kiếm, đặt lịch xem nhà, thanh toán cọc qua platform |
| **Môi giới** | Đăng BĐS, xác nhận lịch, dẫn khách xem nhà, nhận phí |
| **Platform / Admin** | Giữ escrow, xử lý tranh chấp, gửi thông báo tự động |

---

## 2. Luồng Nghiệp Vụ

### 2.1 Luồng Đặt Cọc

```
1. Khách chọn BĐS → điền form (ngày, khung giờ xem nhà)
        ↓ [Validate: không trùng khung giờ của BĐS đó]
2. Thanh toán qua VNPay / Momo / ZaloPay
        ↓ [Tiền vào tài khoản platform - ESCROW]
3. Webhook xác nhận thanh toán → tạo deposit (status = PENDING)
        ↓
4. Gửi notification + email cho môi giới (có 24h để xác nhận)
```

> ⚠️ **Tiền KHÔNG BAO GIỜ vào tài khoản môi giới ở bước này**

---

### 2.2 Môi Giới Xác Nhận

**Trường hợp Từ chối:**
```
Môi giới từ chối + nhập lý do
        ↓
status = BROKER_REJECTED
        ↓
Tự động hoàn 100% tiền cọc + gửi mail thông báo cho khách
```

**Trường hợp Xác nhận:**
```
Môi giới xác nhận
        ↓
status = BROKER_CONFIRMED
        ↓
Gửi mail xác nhận lịch cho khách (địa chỉ, giờ, tên + SĐT môi giới)
        ↓
Đặt cron job: gửi mail nhắc 24h trước buổi xem
```

> ⚠️ **Nếu môi giới không phản hồi trong 24h → tự động BROKER_REJECTED → hoàn tiền khách**

---

### 2.3 OTP Check-in Chống Gian Lận

```
Đến ngày xem nhà
        ↓
Môi giới mở app → Generate OTP (6 số, hiệu lực 10 phút, dùng 1 lần)
        ↓
Khách nhập OTP → Xác nhận cả 2 có mặt
        ↓
status = CHECKED_IN → Mở luồng xử lý kết quả
```

**Fallback nếu không check-in được (sau 2h từ viewing_start):**
```
Mở cửa sổ 24h để 2 bên tự báo cáo
        ├── Cả 2 đồng thuận → Xử lý tiền theo kết quả
        └── Mâu thuẫn → status = DISPUTE → Admin xử lý
```

---

### 2.4 Xử Lý Tiền Theo Kết Quả

| Trường hợp | Tiền cọc khách | Phí môi giới | Status |
|---|---|---|---|
| Môi giới từ chối | Hoàn 100% | Không có | `BROKER_REJECTED` |
| Môi giới timeout 24h | Hoàn 100% | Không có | `BROKER_REJECTED` |
| Khách đến + Mua nhà | Hoàn 100% | Theo hợp đồng riêng | `VISITED_BOUGHT` |
| Khách đến + Không mua | Hoàn (cọc - phí) | Nhận phí môi giới | `VISITED_NOT_BUY` |
| Khách không đến | Mất toàn bộ | Nhận toàn bộ cọc | `NO_SHOW_CUSTOMER` |
| Môi giới không đến | Hoàn 100% | Bị phạt (trừ bảo lãnh) | `NO_SHOW_BROKER` |
| Tranh chấp | Freeze - chờ admin | Freeze - chờ admin | `DISPUTE` |

---

### 2.5 Xử Lý Tranh Chấp (Dispute)

**Khi nào xảy ra Dispute?**
- Môi giới báo "Khách không đến" nhưng khách claim "Tôi có đến"
- Khách báo "Môi giới không đến" nhưng môi giới claim "Tôi có đến"
- OTP không hoạt động, cả 2 bên đều claim có mặt
- Môi giới confirm "Khách không mua" nhưng khách claim "Tôi đã mua"

**Luồng xử lý:**
```
Trigger DISPUTE → Tiền FREEZE tại escrow
        ↓
Cả 2 bên upload bằng chứng (ảnh, screenshot, v.v.) — deadline 48h
        ↓
Admin review: bằng chứng + lịch sử rating môi giới + lịch sử khách
        ↓ (liên hệ 2 bên nếu cần)
Admin ra quyết định:
  REFUND_CUSTOMER / TRANSFER_BROKER / SPLIT
        ↓
Release tiền + Gửi thông báo kết quả cho 2 bên
```

> 🚨 **Tiền KHÔNG được release tự động khi đang DISPUTE. Phải có quyết định của admin.**

---

## 3. Status Flow Đầy Đủ

```
PENDING
  ├── Môi giới từ chối      → BROKER_REJECTED    → hoàn 100%
  ├── Môi giới timeout 24h  → BROKER_REJECTED    → hoàn 100%
  └── Môi giới xác nhận     → BROKER_CONFIRMED
            └── OTP check-in thành công → CHECKED_IN
                      ├── Khách mua        → VISITED_BOUGHT     → hoàn 100%
                      ├── Khách không mua  → VISITED_NOT_BUY    → hoàn (cọc - phí)
                      └── Mâu thuẫn        → DISPUTE            → admin xử lý
            └── OTP timeout 2h → fallback 24h report
                      ├── Khách không đến    → NO_SHOW_CUSTOMER  → mất cọc
                      ├── Môi giới không đến → NO_SHOW_BROKER    → hoàn 100% + phạt
                      └── Mâu thuẫn          → DISPUTE           → admin xử lý

Mọi kết quả cuối cùng → COMPLETED / REFUNDED
```

---

## 4. Thiết Kế Dữ Liệu

### 4.1 Danh Sách Bảng

| Bảng | Mục đích |
|---|---|
| `customers` | Thông tin khách hàng |
| `brokers` | Thông tin môi giới + tiền bảo lãnh |
| `real_estates` | Thông tin BĐS |
| `deposits` | Core: đặt cọc xem nhà |
| `transactions` | Lịch sử mọi giao dịch tiền |
| `disputes` | Tranh chấp cần admin xử lý |
| `notifications` | Log thông báo email/push |
| `broker_ratings` | Đánh giá môi giới sau buổi xem |

---

### 4.2 Chi Tiết Bảng `deposits` (Core)

| Column | Type | Mô tả |
|---|---|---|
| `id` | CHAR(36) | Primary key (UUID) |
| `customer_id` | CHAR(36) FK | Khách hàng đặt cọc |
| `real_estate_id` | CHAR(36) FK | BĐS được đặt cọc |
| `broker_id` | CHAR(36) FK | Môi giới phụ trách |
| `amount` | DECIMAL(15,2) | Số tiền cọc |
| `broker_fee` | DECIMAL(15,2) | Phí môi giới nếu không mua |
| `viewing_date` | DATE | Ngày xem nhà |
| `viewing_start` | TIME | Giờ bắt đầu xem |
| `viewing_end` | TIME | Giờ kết thúc xem |
| `status` | ENUM | Trạng thái hiện tại |
| `payment_method` | ENUM | VNPAY / MOMO / ZALOPAY |
| `payment_ref` | VARCHAR(100) | Mã giao dịch từ cổng thanh toán |
| `otp_hash` | VARCHAR(255) | OTP check-in đã hash |
| `otp_expires_at` | DATETIME | Thời hạn OTP (10 phút) |
| `broker_checkin` | TINYINT(1) | Môi giới xác nhận khách đến |
| `customer_checkin` | TINYINT(1) | Khách xác nhận môi giới đến |
| `broker_confirmed_at` | DATETIME | Thời điểm môi giới duyệt lịch |
| `reminder_sent_at` | DATETIME | Đã gửi mail nhắc 24h chưa |
| `refund_amount` | DECIMAL(15,2) | Số tiền thực tế hoàn cho khách |
| `penalty_amount` | DECIMAL(15,2) | Phí phạt môi giới nếu bùng |
| `reject_reason` | TEXT | Lý do từ chối (nếu có) |
| `created_at` | DATETIME | Thời điểm tạo |
| `updated_at` | DATETIME | Thời điểm cập nhật cuối |

---

## 5. SQL Migration (MySQL)

> ⚠️ **Khác biệt so với PostgreSQL:**
> - Dùng `CHAR(36)` + `UUID()` thay cho `UUID` + `gen_random_uuid()`
> - `ENUM` khai báo inline trong column, không tạo type riêng
> - `DATETIME` thay cho `TIMESTAMP` để tránh auto timezone convert
> - `TINYINT(1)` thay cho `BOOLEAN`
> - `JSON` thay cho `TEXT[]` (array)
> - Check overlap khung giờ bằng điều kiện thủ công thay cho `OVERLAPS`

---

### Bảng `customers`

```sql
CREATE TABLE customers (
  id           CHAR(36)     PRIMARY KEY DEFAULT (UUID()),
  username     VARCHAR(50)  NOT NULL UNIQUE,
  full_name    VARCHAR(100) NOT NULL,
  phone_number VARCHAR(15),
  email        VARCHAR(255) NOT NULL UNIQUE,
  created_at   DATETIME     DEFAULT CURRENT_TIMESTAMP
);
```

### Bảng `brokers`

```sql
CREATE TABLE brokers (
  id                CHAR(36)      PRIMARY KEY DEFAULT (UUID()),
  full_name         VARCHAR(100)  NOT NULL,
  phone_number      VARCHAR(15),
  email             VARCHAR(255)  NOT NULL UNIQUE,
  bank_account      VARCHAR(50),
  bank_name         VARCHAR(100),
  guarantee_deposit DECIMAL(15,2) DEFAULT 0,  -- tiền bảo lãnh
  rating_avg        DECIMAL(3,2)  DEFAULT 5.0,
  total_reviews     INT           DEFAULT 0,
  is_active         TINYINT(1)    DEFAULT 1,
  created_at        DATETIME      DEFAULT CURRENT_TIMESTAMP
);
```

### Bảng `real_estates`

```sql
CREATE TABLE real_estates (
  id         CHAR(36)     PRIMARY KEY DEFAULT (UUID()),
  broker_id  CHAR(36)     NOT NULL,
  title      VARCHAR(255) NOT NULL,
  address    TEXT,
  price      DECIMAL(15,2),
  is_active  TINYINT(1)   DEFAULT 1,
  created_at DATETIME     DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (broker_id) REFERENCES brokers(id)
);
```

### Bảng `deposits`

```sql
CREATE TABLE deposits (
  id              CHAR(36)      PRIMARY KEY DEFAULT (UUID()),
  customer_id     CHAR(36)      NOT NULL,
  real_estate_id  CHAR(36)      NOT NULL,
  broker_id       CHAR(36)      NOT NULL,

  -- Tiền
  amount          DECIMAL(15,2) NOT NULL,
  broker_fee      DECIMAL(15,2) NOT NULL,
  refund_amount   DECIMAL(15,2),
  penalty_amount  DECIMAL(15,2),

  -- Lịch xem nhà
  viewing_date    DATE          NOT NULL,
  viewing_start   TIME          NOT NULL,
  viewing_end     TIME          NOT NULL,

  -- Trạng thái
  status          ENUM(
                    'PENDING',
                    'BROKER_REJECTED',
                    'BROKER_CONFIRMED',
                    'CHECKED_IN',
                    'VISITED_BOUGHT',
                    'VISITED_NOT_BUY',
                    'NO_SHOW_CUSTOMER',
                    'NO_SHOW_BROKER',
                    'DISPUTE',
                    'REFUNDED',
                    'COMPLETED'
                  ) DEFAULT 'PENDING',

  -- Thanh toán
  payment_method  ENUM('VNPAY', 'MOMO', 'ZALOPAY'),
  payment_ref     VARCHAR(100),

  -- OTP check-in
  otp_hash        VARCHAR(255),
  otp_expires_at  DATETIME,

  -- Xác nhận 2 bên
  broker_checkin    TINYINT(1),
  customer_checkin  TINYINT(1),

  -- Timestamps
  broker_confirmed_at DATETIME,
  reminder_sent_at    DATETIME,
  reject_reason       TEXT,
  created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  FOREIGN KEY (customer_id)    REFERENCES customers(id),
  FOREIGN KEY (real_estate_id) REFERENCES real_estates(id),
  FOREIGN KEY (broker_id)      REFERENCES brokers(id),

  -- Index check trùng khung giờ
  INDEX idx_deposit_viewing (real_estate_id, viewing_date, viewing_start, viewing_end)
);
```

### Check trùng khung giờ (Application layer)

```sql
-- Chạy trước khi INSERT deposit mới
SELECT COUNT(*) FROM deposits
WHERE real_estate_id = ?
  AND viewing_date   = ?
  AND status NOT IN ('BROKER_REJECTED', 'REFUNDED', 'COMPLETED')
  AND viewing_start  < ?   -- viewing_end của slot mới
  AND viewing_end    > ?;  -- viewing_start của slot mới
-- Nếu COUNT > 0 → báo lỗi trùng khung giờ
```

### Bảng `transactions`

```sql
CREATE TABLE transactions (
  id         CHAR(36)      PRIMARY KEY DEFAULT (UUID()),
  deposit_id CHAR(36)      NOT NULL,
  type       ENUM(
               'DEPOSIT',
               'REFUND_FULL',
               'REFUND_PARTIAL',
               'TRANSFER_TO_BROKER',
               'PENALTY_BROKER'
             ) NOT NULL,
  amount     DECIMAL(15,2) NOT NULL,
  note       TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (deposit_id) REFERENCES deposits(id)
);
```

### Bảng `disputes`

```sql
CREATE TABLE disputes (
  id            CHAR(36) PRIMARY KEY DEFAULT (UUID()),
  deposit_id    CHAR(36) NOT NULL,
  raised_by     ENUM('CUSTOMER', 'BROKER') NOT NULL,
  reason        TEXT,
  evidence_urls JSON,      -- thay cho TEXT[] của PostgreSQL
  status        ENUM('OPEN', 'REVIEWING', 'RESOLVED') DEFAULT 'OPEN',
  resolution    ENUM('REFUND_CUSTOMER', 'TRANSFER_BROKER', 'SPLIT'),
  resolved_by   CHAR(36),  -- admin id
  resolved_at   DATETIME,
  created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (deposit_id) REFERENCES deposits(id)
);
```

### Bảng `broker_ratings`

```sql
CREATE TABLE broker_ratings (
  id         CHAR(36) PRIMARY KEY DEFAULT (UUID()),
  deposit_id CHAR(36) NOT NULL,
  broker_id  CHAR(36) NOT NULL,
  customer_id CHAR(36) NOT NULL,
  rating     TINYINT  NOT NULL CHECK (rating BETWEEN 1 AND 5),
  comment    TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (deposit_id)  REFERENCES deposits(id),
  FOREIGN KEY (broker_id)   REFERENCES brokers(id),
  FOREIGN KEY (customer_id) REFERENCES customers(id)
);
```

---

## 6. Hệ Thống Thông Báo

### Email Trigger

| Trigger | Gửi cho | Nội dung |
|---|---|---|
| Đặt cọc thành công | Khách + Môi giới | Thông tin lịch, số tiền cọc |
| Môi giới từ chối | Khách | Lý do từ chối, thông báo hoàn tiền |
| Môi giới xác nhận | Khách | Xác nhận lịch, địa chỉ, SĐT môi giới |
| Nhắc trước 24h | Khách + Môi giới | Reminder lịch xem nhà ngày mai |
| Sau xem nhà | Khách + Môi giới | Yêu cầu confirm kết quả trong 24h |
| Hoàn tiền | Khách | Số tiền hoàn, thời gian xử lý |
| Chuyển phí môi giới | Môi giới | Số tiền nhận, thông tin giao dịch |
| Dispute mở | Cả 2 + Admin | Yêu cầu upload bằng chứng trong 48h |
| Dispute resolved | Cả 2 | Kết quả admin, xử lý tiền |

### Cron Jobs

```sql
-- Chạy mỗi giờ: gửi mail nhắc 24h trước
SELECT * FROM deposits
WHERE status = 'BROKER_CONFIRMED'
  AND viewing_date = DATE_ADD(CURDATE(), INTERVAL 1 DAY)
  AND reminder_sent_at IS NULL;
-- → Gửi mail + UPDATE reminder_sent_at = NOW()

-- Chạy mỗi giờ: auto reject nếu môi giới không phản hồi 24h
SELECT * FROM deposits
WHERE status = 'PENDING'
  AND created_at < DATE_SUB(NOW(), INTERVAL 24 HOUR);
-- → UPDATE status = 'BROKER_REJECTED' → trigger hoàn tiền
```

---

## 7. Escrow & Thanh Toán

### Mô hình Escrow

```
Khách thanh toán → Tiền vào tài khoản PLATFORM (không phải môi giới)
        ↓
Tiền FREEZE tại escrow_balance (không ai rút được)
        ↓
Sau kết quả xem nhà → Release theo logic bảng mục 2.4
```

> ✅ Platform tự làm escrow: tiền nằm trong tài khoản ngân hàng của công ty, platform chủ động release theo trạng thái.

### Lưu ý pháp lý

> ⚠️ Platform cần có **pháp nhân (công ty TNHH/CP)** để hợp pháp "giữ tiền hộ" khách hàng. Đăng ký merchant với các cổng thanh toán yêu cầu giấy phép kinh doanh.

---

## 8. Admin Panel

### Tính năng cần có

| Module | Mô tả |
|---|---|
| **Dispute Management** | Xem danh sách dispute, bằng chứng 2 bên, ra quyết định và release tiền |
| **Escrow Dashboard** | Tổng tiền đang giữ, tiền cần release hôm nay, lịch sử giao dịch |
| **Broker Management** | Duyệt môi giới mới, xem rating, lịch sử vi phạm, khóa tài khoản |
| **Deposit Overview** | Tất cả deposits theo status, filter theo ngày/BĐS/môi giới |

### SLA xử lý Dispute

- Nhận alert → bắt đầu review: **< 24h**
- Hoàn tất xử lý: **< 72h**
- Tiền release sau quyết định: **< 1 ngày làm việc**
