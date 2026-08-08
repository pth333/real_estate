# 🏠 Real Estate Slug System — Implementation Plan

> Stack: **Go + GORM** (Backend) · **Nuxt 3** (Frontend) · **MariaDB**  
> Mục tiêu: Implement URL slug system tương tự batdongsan.com.vn

---

## Tổng quan flow

```
User chọn filter
    ↓
Nuxt buildUrl() → /nha-dat-ban-ha-noi/gia-1-den-3-ty?phong-ngu=2
    ↓
Go nhận segments → prefix match category → trim location slug
    ↓
Query filter_ranges theo slug → BETWEEN
    ↓
Query string → filter thêm
    ↓
Trả về listings ✅
```

---

## Phase 1 — Database

### Bước 1.1 — Tạo bảng `categories`

```sql
CREATE TABLE categories (
  id         INT PRIMARY KEY AUTO_INCREMENT,
  name       VARCHAR(100) NOT NULL,
  slug       VARCHAR(100) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Seed data:**
```sql
INSERT INTO categories (name, slug) VALUES
  ('Nhà đất bán',       'nha-dat-ban'),
  ('Nhà đất thuê',      'nha-dat-cho-thue'),
  ('Căn hộ chung cư',   'can-ho-chung-cu'),
  ('Đất nền dự án',     'dat-nen-du-an');
```

---

### Bước 1.2 — Tạo bảng `locations` (3 cấp)

```sql
CREATE TABLE locations (
  id        INT PRIMARY KEY AUTO_INCREMENT,
  name      VARCHAR(100) NOT NULL,
  slug      VARCHAR(100) NOT NULL UNIQUE,
  parent_id INT NULL,
  level     TINYINT NOT NULL COMMENT '1=Tỉnh/TP, 2=Quận/Huyện, 3=Phường/Xã',
  FOREIGN KEY (parent_id) REFERENCES locations(id)
);

-- Index slug để query nhanh
CREATE INDEX idx_locations_slug ON locations(slug);
CREATE INDEX idx_locations_parent ON locations(parent_id);
```

**Seed data mẫu:**
```sql
-- Cấp 1: Tỉnh/TP
INSERT INTO locations (name, slug, parent_id, level) VALUES
  ('Hà Nội',       'ha-noi',    NULL, 1),
  ('TP. Hồ Chí Minh', 'ho-chi-minh', NULL, 1);

-- Cấp 2: Quận/Huyện
INSERT INTO locations (name, slug, parent_id, level) VALUES
  ('Cầu Giấy',  'cau-giay',  1, 2),
  ('Đống Đa',   'dong-da',   1, 2),
  ('Quận 1',    'quan-1',    2, 2),
  ('Quận 7',    'quan-7',    2, 2);
```

---

### Bước 1.3 — Tạo bảng `filter_ranges`

```sql
CREATE TABLE filter_ranges (
  id      INT PRIMARY KEY AUTO_INCREMENT,
  type    ENUM('price', 'area') NOT NULL,
  label   VARCHAR(100) NOT NULL,
  slug    VARCHAR(100) NOT NULL UNIQUE,
  min_val DECIMAL(15,2) NULL COMMENT 'NULL = không giới hạn dưới',
  max_val DECIMAL(15,2) NULL COMMENT 'NULL = không giới hạn trên'
);

CREATE INDEX idx_filter_ranges_slug ON filter_ranges(slug);
```

**Seed data:**
```sql
-- Giá (đơn vị: VNĐ)
INSERT INTO filter_ranges (type, label, slug, min_val, max_val) VALUES
  ('price', 'Dưới 1 tỷ',       'gia-duoi-1-ty',     NULL,           1000000000),
  ('price', 'Từ 1 đến 3 tỷ',   'gia-1-den-3-ty',    1000000000,     3000000000),
  ('price', 'Từ 3 đến 5 tỷ',   'gia-3-den-5-ty',    3000000000,     5000000000),
  ('price', 'Từ 5 đến 10 tỷ',  'gia-5-den-10-ty',   5000000000,     10000000000),
  ('price', 'Trên 10 tỷ',      'gia-tren-10-ty',    10000000000,    NULL);

-- Diện tích (đơn vị: m²)
INSERT INTO filter_ranges (type, label, slug, min_val, max_val) VALUES
  ('area', 'Dưới 30m²',        'dien-tich-duoi-30',  NULL,  30),
  ('area', 'Từ 30 đến 50m²',   'dien-tich-30-50',    30,    50),
  ('area', 'Từ 50 đến 100m²',  'dien-tich-50-100',   50,    100),
  ('area', 'Từ 100 đến 200m²', 'dien-tich-100-200',  100,   200),
  ('area', 'Trên 200m²',       'dien-tich-tren-200', 200,   NULL);
```

---

### Bước 1.4 — Tạo bảng `real_estates`

```sql
CREATE TABLE real_estates (
  id           INT PRIMARY KEY AUTO_INCREMENT,
  title        VARCHAR(255) NOT NULL,
  slug         VARCHAR(255) NOT NULL UNIQUE,  -- slug riêng cho từng listing
  category_id  INT NOT NULL,
  location_id  INT NOT NULL,                  -- Quận/Huyện hoặc Phường/Xã
  price        DECIMAL(15,2) NULL,
  area         DECIMAL(10,2) NULL,
  bedroom      TINYINT NULL,
  direction    VARCHAR(50) NULL,
  description  TEXT NULL,
  created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (category_id) REFERENCES categories(id),
  FOREIGN KEY (location_id) REFERENCES locations(id)
);

CREATE INDEX idx_re_category   ON real_estates(category_id);
CREATE INDEX idx_re_location   ON real_estates(location_id);
CREATE INDEX idx_re_price      ON real_estates(price);
CREATE INDEX idx_re_area       ON real_estates(area);
```

---

## Phase 2 — Go Backend

### Bước 2.1 — GORM Models

```go
// models/category.go
type Category struct {
  ID   uint   `gorm:"primaryKey"`
  Name string `gorm:"size:100;not null"`
  Slug string `gorm:"size:100;not null;uniqueIndex"`
}

// models/location.go
type Location struct {
  ID       uint      `gorm:"primaryKey"`
  Name     string    `gorm:"size:100;not null"`
  Slug     string    `gorm:"size:100;not null;uniqueIndex"`
  ParentID *uint     `gorm:"index"`
  Level    int8      `gorm:"not null"`
  Parent   *Location `gorm:"foreignKey:ParentID"`
}

// models/filter_range.go
type FilterRange struct {
  ID     uint    `gorm:"primaryKey"`
  Type   string  `gorm:"type:enum('price','area');not null"`
  Label  string  `gorm:"size:100;not null"`
  Slug   string  `gorm:"size:100;not null;uniqueIndex"`
  MinVal *float64
  MaxVal *float64
}

// models/real_estate.go
type RealEstate struct {
  ID          uint     `gorm:"primaryKey"`
  Title       string   `gorm:"size:255;not null"`
  Slug        string   `gorm:"size:255;not null;uniqueIndex"`
  CategoryID  uint     `gorm:"not null;index"`
  LocationID  uint     `gorm:"not null;index"`
  Price       *float64 `gorm:"index"`
  Area        *float64 `gorm:"index"`
  Bedroom     *int8
  Direction   *string  `gorm:"size:50"`
  Description *string  `gorm:"type:text"`
  Category    Category `gorm:"foreignKey:CategoryID"`
  Location    Location `gorm:"foreignKey:LocationID"`
}
```

---

### Bước 2.2 — Slug Parser Service

```go
// services/slug_parser.go
package services

import (
  "strings"
  "your-project/models"
  "gorm.io/gorm"
)

type ParsedSlug struct {
  Category    *models.Category
  Location    *models.Location
  PriceRange  *models.FilterRange
  AreaRange   *models.FilterRange
}

type SlugParser struct {
  db *gorm.DB
}

func NewSlugParser(db *gorm.DB) *SlugParser {
  return &SlugParser{db: db}
}

func (p *SlugParser) Parse(segments []string) (*ParsedSlug, error) {
  result := &ParsedSlug{}

  if len(segments) == 0 {
    return result, nil
  }

  // --- Parse segment 0: category[-location] ---
  var categories []models.Category
  p.db.Find(&categories)

  seg0 := segments[0]
  for _, cat := range categories {
    if strings.HasPrefix(seg0, cat.Slug) {
      c := cat
      result.Category = &c

      // Tách location slug ra
      locationSlug := strings.TrimPrefix(seg0, cat.Slug)
      locationSlug  = strings.TrimPrefix(locationSlug, "-")

      if locationSlug != "" {
        var loc models.Location
        if err := p.db.Where("slug = ?", locationSlug).First(&loc).Error; err == nil {
          result.Location = &loc
        }
      }
      break
    }
  }

  // --- Parse segment 1+: filter_ranges ---
  for _, seg := range segments[1:] {
    var f models.FilterRange
    if err := p.db.Where("slug = ?", seg).First(&f).Error; err == nil {
      switch f.Type {
      case "price":
        result.PriceRange = &f
      case "area":
        result.AreaRange = &f
      }
    }
  }

  return result, nil
}
```

---

### Bước 2.3 — Listing Handler

```go
// handlers/listing.go
package handlers

import (
  "strings"
  "your-project/models"
  "your-project/services"
  "github.com/gin-gonic/gin"
  "gorm.io/gorm"
)

type ListingHandler struct {
  db     *gorm.DB
  parser *services.SlugParser
}

func NewListingHandler(db *gorm.DB) *ListingHandler {
  return &ListingHandler{
    db:     db,
    parser: services.NewSlugParser(db),
  }
}

func (h *ListingHandler) GetListings(c *gin.Context) {
  // Nhận segments từ Nuxt
  segmentsRaw := c.Query("segments")
  segments := []string{}
  if segmentsRaw != "" {
    segments = strings.Split(segmentsRaw, ",")
  }

  // Parse slug
  parsed, err := h.parser.Parse(segments)
  if err != nil {
    c.JSON(400, gin.H{"error": "invalid slug"})
    return
  }

  // Build query
  query := h.db.Model(&models.RealEstate{})

  if parsed.Category != nil {
    query = query.Where("category_id = ?", parsed.Category.ID)
  }
  if parsed.Location != nil {
    query = query.Where("location_id = ?", parsed.Location.ID)
  }
  if parsed.PriceRange != nil {
    if parsed.PriceRange.MinVal != nil {
      query = query.Where("price >= ?", *parsed.PriceRange.MinVal)
    }
    if parsed.PriceRange.MaxVal != nil {
      query = query.Where("price <= ?", *parsed.PriceRange.MaxVal)
    }
  }
  if parsed.AreaRange != nil {
    if parsed.AreaRange.MinVal != nil {
      query = query.Where("area >= ?", *parsed.AreaRange.MinVal)
    }
    if parsed.AreaRange.MaxVal != nil {
      query = query.Where("area <= ?", *parsed.AreaRange.MaxVal)
    }
  }

  // Query string filters
  if bedroom := c.Query("phong-ngu"); bedroom != "" {
    query = query.Where("bedroom = ?", bedroom)
  }
  if direction := c.Query("huong"); direction != "" {
    query = query.Where("direction = ?", direction)
  }

  // Pagination
  page := c.DefaultQuery("page", "1")
  limit := 20
  query = query.Limit(limit)

  var listings []models.RealEstate
  query.Preload("Category").Preload("Location").Find(&listings)

  c.JSON(200, gin.H{
    "data":     listings,
    "category": parsed.Category,
    "location": parsed.Location,
  })
}
```

---

### Bước 2.4 — Router

```go
// main.go hoặc router.go
listingHandler := handlers.NewListingHandler(db)

r := gin.Default()
r.GET("/api/listings", listingHandler.GetListings)
```

---

## Phase 3 — Nuxt Frontend

### Bước 3.1 — Utility buildUrl

```ts
// utils/buildUrl.ts
interface BuildUrlParams {
  category: { slug: string }
  location?: { slug: string }
  priceSlug?: string
  areaSlug?: string
  query?: Record<string, string>
}

export function buildUrl({
  category,
  location,
  priceSlug,
  areaSlug,
  query = {}
}: BuildUrlParams): string {
  // Segment 1: category[-location]
  const seg1 = location
    ? `${category.slug}-${location.slug}`
    : category.slug

  // Segment 2+: filter ranges
  const filterSegs = [priceSlug, areaSlug].filter(Boolean) as string[]

  const path = `/${[seg1, ...filterSegs].join('/')}`

  const qs = new URLSearchParams(query).toString()
  return qs ? `${path}?${qs}` : path
}
```

---

### Bước 3.2 — Page Route

```
pages/
  [...slug].vue   ← bắt tất cả URL dạng /nha-dat-ban/...
```

```vue
<!-- pages/[...slug].vue -->
<script setup lang="ts">
const route = useRoute()
const config = useRuntimeConfig()

// Parse segments từ URL
const segments = computed(() =>
  (route.params.slug as string[]) || []
)

// Gọi API Go
const { data, pending } = await useFetch('/api/listings', {
  baseURL: config.public.apiBase,
  params: computed(() => ({
    segments: segments.value.join(','),
    ...route.query  // forward query string
  }))
})

// SEO meta từ response
useHead({
  title: computed(() =>
    data.value?.category?.name
      ? `${data.value.category.name}${data.value.location ? ' tại ' + data.value.location.name : ''}`
      : 'Bất động sản'
  )
})
</script>

<template>
  <div>
    <ListingGrid :listings="data?.data" :loading="pending" />
  </div>
</template>
```

---

### Bước 3.3 — Filter Component navigate

```vue
<!-- components/FilterBar.vue -->
<script setup lang="ts">
import { buildUrl } from '~/utils/buildUrl'

const props = defineProps<{
  category: { slug: string }
  location?: { slug: string }
}>()

const selectedPrice = ref('')
const selectedArea = ref('')
const selectedBedroom = ref('')

function applyFilter() {
  const url = buildUrl({
    category: props.category,
    location: props.location,
    priceSlug: selectedPrice.value || undefined,
    areaSlug: selectedArea.value || undefined,
    query: {
      ...(selectedBedroom.value && { 'phong-ngu': selectedBedroom.value })
    }
  })
  navigateTo(url)
}
</script>
```

---

### Bước 3.4 — nuxt.config.ts

```ts
// nuxt.config.ts
export default defineNuxtConfig({
  runtimeConfig: {
    public: {
      apiBase: process.env.API_BASE_URL || 'http://localhost:8080'
    }
  }
})
```

---

## Phase 4 — Kiểm tra & Edge Cases

### Bước 4.1 — Test cases URL

| URL | Category | Location | Price | Area |
|-----|----------|----------|-------|------|
| `/nha-dat-ban` | ✅ nha-dat-ban | ❌ null | ❌ null | ❌ null |
| `/nha-dat-ban-ha-noi` | ✅ nha-dat-ban | ✅ ha-noi | ❌ null | ❌ null |
| `/nha-dat-ban/gia-1-den-3-ty` | ✅ nha-dat-ban | ❌ null | ✅ | ❌ null |
| `/nha-dat-ban-ha-noi/gia-1-den-3-ty` | ✅ | ✅ | ✅ | ❌ null |
| `/nha-dat-ban-ha-noi/gia-1-den-3-ty/dien-tich-50-100` | ✅ | ✅ | ✅ | ✅ |

### Bước 4.2 — Edge case: category slug là prefix của category khác

```
'nha-dat-ban' vs 'nha-dat-ban-cho-thue' → conflict!
```

**Fix:** Sort categories theo độ dài slug **giảm dần** trước khi prefix match:

```go
// Trong SlugParser.Parse()
sort.Slice(categories, func(i, j int) bool {
  return len(categories[i].Slug) > len(categories[j].Slug)
})
```

### Bước 4.3 — Cache categories (tối ưu)

```go
// Không query DB mỗi request — cache in-memory
var categoryCache []models.Category

func (p *SlugParser) getCategories() []models.Category {
  if len(categoryCache) == 0 {
    p.db.Find(&categoryCache)
  }
  return categoryCache
}
```

---

## Phase 5 — Slug trang chi tiết

> URL: `/nha-pho-2-tang-cau-giay-ha-noi-rs123`  
> Cơ chế detect: slug có suffix **`rs` + số** → trang chi tiết, ngược lại → 404

---

### Bước 5.1 — Generate slug cho real_estate

```go
// utils/slug.go
package utils

import (
  "fmt"
  "regexp"
  "strings"
  "golang.org/x/text/unicode/norm"
)

// Chuyển tiếng Việt → slug
func ToSlug(s string) string {
  s = norm.NFD.String(s)
  re := regexp.MustCompile(`\p{Mn}`)
  s = re.ReplaceAllString(s, "")
  s = strings.ToLower(s)
  s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
  s = strings.Trim(s, "-")
  return s
}

// Generate slug cho listing: {title-slug}-rs{id}
func GenerateListingSlug(title string, id uint) string {
  return fmt.Sprintf("%s-rs%d", ToSlug(title), id)
}
```

**Gọi khi tạo listing:**
```go
// Insert trước để có ID
db.Create(&listing)

// Generate slug sau khi có ID
listing.Slug = utils.GenerateListingSlug(listing.Title, listing.ID)
// → "nha-pho-2-tang-cau-giay-ha-noi-rs123"

db.Save(&listing)
```

---

### Bước 5.2 — Go detect chi tiết vs filter

```go
// services/slug_parser.go
import (
  "regexp"
  "strconv"
)

var detailPattern = regexp.MustCompile(`-rs(\d+)$`)

func IsDetailSlug(segment string) (id uint, isDetail bool) {
  match := detailPattern.FindStringSubmatch(segment)
  if match == nil {
    return 0, false
  }
  n, _ := strconv.ParseUint(match[1], 10, 64)
  return uint(n), true
}
```

**Handler tổng:**
```go
func (h *ListingHandler) Handle(c *gin.Context) {
  segmentsRaw := c.Query("segments")
  segments := strings.Split(segmentsRaw, ",")

  // Chỉ có 1 segment + có pattern rs → trang chi tiết
  if len(segments) == 1 {
    if id, isDetail := IsDetailSlug(segments[0]); isDetail {
      h.getDetail(c, id)
      return
    }
  }

  // Còn lại → trang filter
  h.getListings(c, segments)
}

func (h *ListingHandler) getDetail(c *gin.Context, id uint) {
  var listing models.RealEstate
  if err := h.db.
    Preload("Category").
    Preload("Location").
    First(&listing, id).Error; err != nil {
    c.JSON(404, gin.H{"error": "not found"})
    return
  }
  c.JSON(200, gin.H{"data": listing, "type": "detail"})
}
```

---

### Bước 5.3 — Nuxt route

```
pages/
  [...slug].vue   ← bắt cả filter lẫn chi tiết
```

```vue
<!-- pages/[...slug].vue -->
<script setup lang="ts">
const route = useRoute()
const segments = (route.params.slug as string[]) || []

const { data } = await useFetch('/api/listings', {
  params: {
    segments: segments.join(','),
    ...route.query
  }
})

const isDetail = computed(() => data.value?.type === 'detail')

useHead({
  title: computed(() =>
    isDetail.value
      ? data.value?.data?.title         // 'Nhà phố 2 tầng Cầu Giấy'
      : data.value?.category?.name      // 'Nhà đất bán tại Hà Nội'
  )
})
</script>

<template>
  <ListingDetail v-if="isDetail" :listing="data?.data" />
  <ListingGrid   v-else          :listings="data?.data" />
</template>
```

---

### Bước 5.4 — buildDetailUrl utility

```ts
// utils/buildUrl.ts
export function buildDetailUrl(listingSlug: string): string {
  return `/${listingSlug}`
  // → /nha-pho-2-tang-cau-giay-ha-noi-rs123
}
```

---

### Bước 5.5 — Test cases

| URL | Detect | Kết quả |
|-----|--------|---------|
| `/nha-pho-2-tang-rs123` | `rs123` match → detail | ✅ Chi tiết ID 123 |
| `/nha-dat-ban-ha-noi/gia-1-den-3-ty` | không có `rs` → filter | ✅ Trang filter |
| `/nha-dat-ban/dien-tich-50-100` | `100` không có `rs` → filter | ✅ Không conflict |
| `/nha-pho-dep-rs456` | `rs456` match → detail | ✅ Chi tiết ID 456 |

---

## Checklist tổng

- [ ] **Phase 1.1** Tạo bảng `categories` + seed data
- [ ] **Phase 1.2** Tạo bảng `locations` + seed data
- [ ] **Phase 1.3** Tạo bảng `filter_ranges` + seed data
- [ ] **Phase 1.4** Tạo bảng `real_estates`
- [ ] **Phase 2.1** GORM models
- [ ] **Phase 2.2** SlugParser service
- [ ] **Phase 2.3** Listing handler
- [ ] **Phase 2.4** Router setup
- [ ] **Phase 3.1** `buildUrl` utility
- [ ] **Phase 3.2** `[...slug].vue` page
- [ ] **Phase 3.3** FilterBar component
- [ ] **Phase 3.4** nuxt.config.ts
- [ ] **Phase 4.1** Test all URL cases
- [ ] **Phase 4.2** Fix prefix conflict
- [ ] **Phase 4.3** Cache categories
- [ ] **Phase 5.1** `ToSlug` + `GenerateListingSlug` utility
- [ ] **Phase 5.2** `IsDetailSlug` detect trong Go handler
- [ ] **Phase 5.3** `[...slug].vue` render detail vs listing
- [ ] **Phase 5.4** `buildDetailUrl` utility
- [ ] **Phase 5.5** Test edge case filter slug có số cuối