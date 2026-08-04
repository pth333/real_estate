# BDS Recommendation System — Full Technical Plan

> Stack: Vue (FE) · Go (Backend) · Python (ML Service) · PostgreSQL · Redis · gRPC

---

## 1. Tổng quan Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        FRONTEND (Vue)                       │
│   Tracking Events: view_duration | search | location        │
└──────────────────────────┬──────────────────────────────────┘
                           │ HTTP REST
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    GO BACKEND (API Gateway)                  │
│                                                             │
│   /tracking/*   →  Ghi events vào PostgreSQL               │
│   /recommend/*  →  Check Redis → gRPC call Python          │
└───────┬──────────────────────────────────┬──────────────────┘
        │ SQL                              │ gRPC
        ▼                                 ▼
┌───────────────┐                ┌─────────────────────────┐
│  PostgreSQL   │◄───────────────│   PYTHON ML SERVICE     │
│               │   Query data   │                         │
│  view_history │                │  Content-Based Filter   │
│  search_hist  │                │  Behavioral Filter      │
│  properties   │                │  Scoring & Ranking      │
└───────────────┘                └────────────┬────────────┘
                                              │ Write precomputed
        ┌─────────────────────────────────────▼────────────┐
        │                   REDIS CACHE                     │
        │   rec:user:{id}   →  [property_ids]  TTL=1h      │
        │   rec:prop:{id}   →  [property_ids]  TTL=2h      │
        └──────────────────────────────────────────────────┘
```

---

## 2. Frontend Tracking Layer (Vue)

### 2.1 Track thời gian xem (View Duration)

```
User vào trang /property/:id
    │
    ├─ onMounted() → ghi startTime = Date.now()
    │
    ├─ onBeforeUnmount() → tính duration = Date.now() - startTime
    │                    → gọi POST /tracking/view
    │
    └─ visibilitychange event → pause timer khi tab bị ẩn
```

**Payload gửi lên:**
```json
{
  "property_id": "uuid",
  "duration_seconds": 120,
  "session_id": "uuid"   // dùng nếu chưa login
}
```

**Rule lọc nhiễu (xử lý bên Go):**
- duration < 5s → bỏ qua (vào nhầm)
- duration > 3600s → cap lại 3600 (để tab rồi bỏ đó)

---

### 2.2 Track Search Query

```
User gõ vào search bar
    │
    ├─ Debounce 500ms
    │
    └─ Khi user submit search → gọi POST /tracking/search
```

**Payload:**
```json
{
  "query": "căn hộ 2 phòng ngủ quận 1",
  "filters": {
    "province": "Hà Nội",
    "district": "Cầu Giấy",
    "price_min": 2000000000,
    "price_max": 5000000000,
    "property_type": "apartment"
  }
}
```

---

### 2.3 Track Vị trí

```
User bật "Tìm quanh đây" hoặc filter theo vị trí
    │
    ├─ Hỏi permission GPS (Geolocation API)
    │
    ├─ Có permission → lấy lat/lng thật
    │
    └─ Không có → fallback lấy từ IP (bên Go xử lý)
         → gọi POST /tracking/location
```

**Payload:**
```json
{
  "latitude": 21.0285,
  "longitude": 105.8542,
  "source": "gps" | "ip_fallback" | "user_selected"
}
```

---

### 2.4 Fetch Recommendation

```
User vào trang /property/:id
    │
    ├─ Gọi GET /recommend?property_id=xxx&user_id=xxx
    │
    └─ Render danh sách "Bạn có thể thích"
```

---

## 3. Go Backend

### 3.1 Các Endpoint Tracking

| Method | Path | Mô tả |
|--------|------|--------|
| POST | /tracking/view | Lưu lịch sử xem |
| POST | /tracking/search | Lưu lịch sử tìm kiếm |
| POST | /tracking/location | Cập nhật vị trí user |
| GET | /recommend | Lấy danh sách gợi ý |

---

### 3.2 Recommendation Flow

```
GET /recommend?property_id=xxx&user_id=xxx
    │
    ├─ [1] Check Redis: GET rec:user:{user_id}
    │       │
    │       ├─ HIT  → trả về ngay ✓
    │       │
    │       └─ MISS → [2] gọi gRPC sang Python
    │                       │
    │                       └─ Nhận [property_ids]
    │                           → SET Redis rec:user:{id} TTL=1h
    │                           → Trả về client ✓
```

---

### 3.3 gRPC Client trong Go

**Cấu trúc call:**
```
Lấy từ DB: view_history của user (top 20 gần nhất)
         + location hiện tại của user
         + property_id đang xem

Build RecommendRequest → gRPC call → nhận RecommendResponse
```

---

### 3.4 Middleware cần có

- **Auth middleware** — lấy user_id từ JWT, fallback sang session_id
- **Rate limit** — chống spam tracking endpoint
- **Validation** — kiểm tra duration, lat/lng hợp lệ

---

## 4. PostgreSQL Schema

### 4.1 Bảng Tracking

```sql
-- Lịch sử xem BDS
CREATE TABLE view_history (
    id              BIGSERIAL PRIMARY KEY,
    user_id         UUID,                          -- NULL nếu chưa login
    session_id      UUID NOT NULL,
    property_id     UUID NOT NULL REFERENCES properties(id),
    duration        INTEGER NOT NULL,              -- seconds, đã cap 5-3600
    viewed_at       TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_view_history_user    ON view_history(user_id, viewed_at DESC);
CREATE INDEX idx_view_history_session ON view_history(session_id, viewed_at DESC);

-- Lịch sử tìm kiếm
CREATE TABLE search_history (
    id              BIGSERIAL PRIMARY KEY,
    user_id         UUID,
    session_id      UUID NOT NULL,
    query           TEXT,
    province        VARCHAR(100),
    district        VARCHAR(100),
    property_type   VARCHAR(50),
    price_min       BIGINT,
    price_max       BIGINT,
    searched_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_search_history_user ON search_history(user_id, searched_at DESC);

-- Vị trí user (1 row per user, upsert)
CREATE TABLE user_location (
    user_id         UUID PRIMARY KEY,
    latitude        DOUBLE PRECISION NOT NULL,
    longitude       DOUBLE PRECISION NOT NULL,
    source          VARCHAR(20),                   -- gps | ip_fallback | user_selected
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
```

### 4.2 Bảng Properties (tóm tắt các cột liên quan)

```sql
CREATE TABLE properties (
    id              UUID PRIMARY KEY,
    title           TEXT,
    property_type   VARCHAR(50),                   -- apartment | house | land | ...
    province        VARCHAR(100),
    district        VARCHAR(100),
    ward            VARCHAR(100),
    latitude        DOUBLE PRECISION,
    longitude       DOUBLE PRECISION,
    price           BIGINT,
    area            FLOAT,
    bedrooms        INTEGER,
    status          VARCHAR(20),                   -- active | sold | hidden
    project_id      UUID,                          -- thuộc dự án nào
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Index cho query recommendation
CREATE INDEX idx_properties_location ON properties(province, district);
CREATE INDEX idx_properties_type     ON properties(property_type, status);
CREATE INDEX idx_properties_project  ON properties(project_id);
CREATE INDEX idx_properties_latlon   ON properties USING GIST (
    point(longitude, latitude)                     -- spatial query quanh khu vực
);
```

---

## 5. Redis Cache Strategy

| Key Pattern | Value | TTL | Ghi chú |
|-------------|-------|-----|---------|
| `rec:user:{user_id}` | JSON array of property_ids | 1h | Recommendation cho logged-in user |
| `rec:session:{session_id}` | JSON array of property_ids | 30m | Recommendation cho guest |
| `rec:prop:{property_id}` | JSON array of property_ids | 2h | "Tương tự BDS này" |

**Invalidation:**
- User xem BDS mới → xóa `rec:user:{user_id}` → lần sau sẽ recompute
- BDS bị ẩn/bán → xóa tất cả cache liên quan đến property đó

---

## 6. gRPC Contract (Proto)

```protobuf
syntax = "proto3";
package recommendation;

service RecommendationService {
    // Lấy gợi ý cho user
    rpc GetRecommendations (RecommendRequest) returns (RecommendResponse);

    // Precompute batch (gọi từ background job)
    rpc PrecomputeForUser (PrecomputeRequest) returns (PrecomputeResponse);
}

// === Request / Response ===

message RecommendRequest {
    string user_id         = 1;  // empty nếu guest
    string session_id      = 2;
    string property_id     = 3;  // BDS đang xem
    double latitude        = 4;
    double longitude       = 5;
    int32  limit           = 6;  // số lượng gợi ý, default 10
}

message RecommendResponse {
    repeated string property_ids = 1;
    string          strategy     = 2;  // "content_based" | "behavioral" | "hybrid" | "fallback"
}

message ViewedProperty {
    string property_id     = 1;
    int32  duration        = 2;  // seconds
    int64  viewed_at_unix  = 3;
}

message PrecomputeRequest {
    string                  user_id         = 1;
    repeated ViewedProperty view_history    = 2;
    double                  latitude        = 3;
    double                  longitude       = 4;
}

message PrecomputeResponse {
    bool   success         = 1;
    int32  count           = 2;  // số gợi ý đã tính
}
```

---

## 7. Python ML Service

### 7.1 Cấu trúc Service

```
ml-service/
├── main.py                  # gRPC server entrypoint
├── servicer.py              # implement RecommendationServicer
├── strategies/
│   ├── content_based.py     # lọc theo đặc điểm BDS
│   ├── behavioral.py        # lọc theo lịch sử xem
│   └── hybrid.py            # merge + scoring
├── db.py                    # kết nối PostgreSQL (read-only)
├── proto/                   # generated từ .proto file
└── scheduler.py             # APScheduler batch job
```

---

### 7.2 Recommendation Logic

**Content-Based Filter** (khi có property_id đang xem):
```
Lấy thông tin property_id hiện tại
    │
    ├─ Query BDS cùng quận/huyện
    ├─ Lọc: cùng loại hình (apartment/house/land)
    ├─ Lọc: giá trong khoảng ±30%
    ├─ Lọc: diện tích ±40%
    ├─ Ưu tiên: cùng dự án (project_id)
    └─ Spatial query: bán kính 3km nếu còn thiếu
```

**Behavioral Filter** (khi có user_id + view_history):
```
Lấy top 10 BDS xem nhiều nhất (weight theo duration)
    │
    ├─ Tìm pattern: loại hình thường xem, khu vực thường xem, range giá
    └─ Query BDS phù hợp với pattern đó
```

**Scoring & Merge:**
```
Score mỗi BDS = (content_score * 0.6) + (behavioral_score * 0.4)

Fallback nếu thiếu data:
    - Không có user_id → content_based 100%
    - Không có property_id → behavioral 100%
    - Không có cả hai → trending theo khu vực
```

---

### 7.3 Background Batch Job (APScheduler)

```
Chạy mỗi 30 phút:
    │
    ├─ Query PostgreSQL: users có activity trong 24h
    │
    ├─ Với mỗi user:
    │   ├─ Lấy view_history (20 bản ghi gần nhất)
    │   ├─ Lấy location
    │   ├─ Chạy hybrid recommendation
    │   └─ Push kết quả vào Redis: rec:user:{id} TTL=2h
    │
    └─ Log: bao nhiêu user đã precompute, thời gian chạy
```

---

## 8. Sequence Diagram — Full Flow

### Khi user xem BDS (không có cache)

```
FE          Go Backend       Redis       Go→Python(gRPC)   PostgreSQL
│                │              │               │               │
│ vào /prop/:id  │              │               │               │
│───────────────►│              │               │               │
│                │ GET rec:user │               │               │
│                │─────────────►│               │               │
│                │   MISS       │               │               │
│                │◄─────────────│               │               │
│                │              │  GetRecs RPC  │               │
│                │──────────────│──────────────►│               │
│                │              │               │ query props   │
│                │              │               │──────────────►│
│                │              │               │◄──────────────│
│                │              │               │ [prop_ids]    │
│                │◄─────────────│───────────────│               │
│                │ SET rec:user TTL=1h           │               │
│                │─────────────►│               │               │
│ [prop_ids]     │              │               │               │
│◄───────────────│              │               │               │

Đồng thời (async, không block response):
│ onBeforeUnmount│              │               │               │
│───────────────►│ POST /tracking/view          │               │
│                │──────────────│──────────────►│──────────────►│
│                │              │               │   INSERT      │
```

---

## 9. Checklist Implement theo thứ tự

### Phase 1 — Foundation
- [ ] Tạo PostgreSQL schema (view_history, search_history, user_location)
- [ ] Go: implement 3 tracking endpoints
- [ ] Frontend: implement tracking composable (useTracking.ts)

### Phase 2 — gRPC Setup
- [ ] Viết file `.proto`
- [ ] Generate code cho Go (protoc-gen-go) và Python (grpcio-tools)
- [ ] Go: tạo gRPC client
- [ ] Python: tạo gRPC server skeleton

### Phase 3 — ML Logic
- [ ] Python: content_based.py
- [ ] Python: behavioral.py
- [ ] Python: hybrid scoring
- [ ] Test với data mẫu

### Phase 4 — Cache & Batch
- [ ] Go: tích hợp Redis check/set
- [ ] Python: APScheduler batch job
- [ ] Test invalidation flow

### Phase 5 — Integration
- [ ] Kết nối end-to-end: FE → Go → Redis/gRPC → Python → PostgreSQL
- [ ] Test full flow với data thật

---

## 10. Tech Stack Summary

| Layer | Tech |
|-------|------|
| Frontend | Vue 3 + Composition API |
| Backend | Go (Gin/Fiber) |
| ML Service | Python 3.11 + gRPC |
| Database | PostgreSQL 15 |
| Cache | Redis 7 |
| Communication | gRPC + Protobuf |
| Background Job | APScheduler (Python) |
| Spatial Query | PostGIS hoặc Haversine formula |
