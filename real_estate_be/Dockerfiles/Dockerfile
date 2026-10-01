FROM golang:1.26.0-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest

# tzdata để mọi thứ trong container hiểu múi giờ tên (mặc định alpine KHÔNG có)
# + ép TZ = giờ VN. Code Go vẫn tự ép UTC+7 (pkg/vntime) nên không phụ thuộc bước này.
RUN apk add --no-cache tzdata
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/config ./config

EXPOSE 8000

CMD ["./server"]