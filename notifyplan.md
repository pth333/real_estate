# Plan: Luồng thông báo BĐS mới cho tất cả user (SSE + Kafka)

---

## Tổng quan luồng

```
Admin đăng BĐS mới (instance bất kỳ)
  → Lưu DB thành công
  → Produce message → Kafka Topic "new-listing"
  → Tất cả instance consume topic đó
  → Mỗi instance broadcast SSE tới tất cả user đang connect vào mình
  → Frontend nhận event → hiển thị thông báo
```

---

## Phần 1 — Chuẩn bị dữ liệu

### 1.1 Bảng `notifications`
- Lưu thông báo để user offline vẫn nhận được khi mở lại app
- Đây là thông báo global — 1 BĐS mới = 1 bản ghi duy nhất
- Các trường: `id`, `listing_id` (unique), `type`, `payload` (JSON: title, slug, price), `created_at`
- `listing_id` là unique key — dùng `ON DUPLICATE KEY IGNORE` để tránh INSERT trùng khi nhiều instance chạy

### 1.2 API
- `GET /api/notifications` — lấy 50 thông báo gần nhất, không cần auth

---

## Phần 2 — SSE Hub (Backend)

### 2.1 Cấu trúc Hub
- Hub là singleton chạy suốt vòng đời ứng dụng
- Lưu tất cả connection dạng `map[connID]chan string` — connID là UUID sinh lúc connect
- Dùng `sync.RWMutex` tránh race condition

### 2.2 Vòng đời một connection
1. User mở trang web → frontend gọi `GET /api/notifications/stream`
2. Hub sinh UUID cho connection, tạo channel, thêm vào map
3. Connection giữ mở, lắng nghe event từ channel
4. Khi user đóng tab / mất mạng → `request.Context()` bị cancel
5. Backend detect cancel → xoá channel khỏi Hub theo UUID → đóng channel

### 2.3 Cơ chế Hub.Broadcast
1. Lock read, loop qua tất cả channel trong map
2. Với mỗi channel: thử push message vào
3. Nếu channel đầy (buffer full) → bỏ qua, không block

### 2.4 Header SSE bắt buộc
- `Content-Type: text/event-stream`
- `Cache-Control: no-cache`
- `Connection: keep-alive`
- `X-Accel-Buffering: no` — bắt buộc nếu chạy sau Nginx

---

## Phần 3 — Luồng Admin đăng BĐS

### 3.1 Các bước xử lý trong Service
1. Validate input từ Admin
2. Lưu bản ghi BĐS mới vào DB
3. Nếu lưu thành công → Produce message vào Kafka topic `new-listing`
   - Payload: `listing_id`, `title`, `slug`, `price`, `created_at`
   - Partition key: `listing_id`
4. Trả response thành công cho Admin ngay lập tức — không chờ Kafka xác nhận

---

## Phần 4 — Kafka

### 4.1 Topic `new-listing`
- Số partition: bằng số instance dự kiến chạy đồng thời
- Retention: 24h
- `acks=1` — đủ đảm bảo, không cần wait all replicas

### 4.2 Consumer (chạy trên mỗi instance)
1. Mỗi instance khởi động → spawn 1 goroutine chạy vòng lặp poll Kafka liên tục
2. Mỗi instance dùng **Consumer Group riêng** (ví dụ: `notification-{instanceID}`)
   - Lý do: cần tất cả instance đều nhận message để broadcast cho user của mình
3. Khi nhận được message:
   - Parse payload
   - INSERT vào bảng `notifications` với `ON DUPLICATE KEY IGNORE`
   - Gọi `Hub.Broadcast(payload)` → push SSE tới tất cả user đang connect vào instance này
4. Commit offset sau khi xử lý xong

### 4.3 Xử lý lỗi
| Tình huống | Cách xử lý |
|---|---|
| Instance crash giữa chừng | Chưa commit offset → Kafka re-deliver khi instance restart |
| DB timeout khi INSERT | Không commit offset → retry ở lần poll tiếp theo |
| Kafka broker down | Producer buffer local, tự retry khi broker lên lại |
| Message duplicate | `ON DUPLICATE KEY IGNORE` với `listing_id` unique |

---

## Phần 5 — Frontend (Vue/Nuxt)

### 5.1 Khởi tạo kết nối SSE
1. User mở trang web → gọi `new EventSource('/api/notifications/stream')`
2. Kết nối giữ mở liên tục trong suốt session
3. Đặt trong `app.vue` hoặc layout chính

### 5.2 Xử lý khi nhận event
1. Parse JSON từ `event.data`
2. Prepend vào danh sách `notifications` (mới nhất lên đầu)
3. Tăng badge counter trên icon chuông
4. Hiển thị toast: tên BĐS + giá, có nút "Xem ngay" → navigate tới trang chi tiết

### 5.3 Xử lý reconnect
1. `EventSource.onerror` bắn → đóng connection hiện tại
2. Chờ 5 giây → reconnect
3. Lặp lại cho đến khi thành công

### 5.4 Load thông báo cũ khi mở app
1. Gọi `GET /api/notifications` lấy 50 bản ghi gần nhất
2. Lưu timestamp lần cuối đọc vào `localStorage`
3. Badge counter = số bản ghi có `created_at` > timestamp đó

### 5.5 Disconnect
1. Khi component unmount → đóng `EventSource`

---

## Phần 6 — Edge case

| Tình huống | Xử lý |
|---|---|
| User mở nhiều tab | Hub broadcast tới tất cả connection → mọi tab đều nhận |
| Server restart | Frontend tự reconnect sau 5s |
| Nginx đứng trước | Header `X-Accel-Buffering: no` |
| User offline khi có BĐS mới | Load từ `GET /api/notifications` khi mở lại app |
| Nhiều user online cùng lúc | Go xử lý tốt hàng nghìn goroutine concurrent |
| Scale nhiều instance | Mỗi instance Consumer Group riêng → tất cả đều nhận và broadcast |

---

## Thứ tự implement

1. Tạo bảng `notifications` + migration (`listing_id` unique)
2. Tạo Kafka topic `new-listing`
3. Dựng SSE Hub với `Broadcast`
4. Endpoint `GET /api/notifications/stream`
5. Kafka Producer trong `RealEstateService.Create`
6. Kafka Consumer — spawn goroutine khi app khởi động, xử lý broadcast + INSERT
7. API `GET /api/notifications`
8. Frontend composable `useNotificationSSE` — connect, reconnect, parse event
9. Toast + badge + đọc trạng thái từ `localStorage`
10. Test multi-instance: chạy 2 instance, đăng BĐS → cả 2 đều broadcast