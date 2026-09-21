/**
 * NotificationStream — client SSE dùng `fetch` + ReadableStream.
 *
 * VÌ SAO KHÔNG DÙNG `EventSource`:
 * `EventSource` của trình duyệt KHÔNG cho phép gửi header, nên không thể gắn
 * `Authorization: Bearer <token>`. Route `/notifications/stream` yêu cầu đăng nhập
 * → EventSource luôn nhận 401. Cách dùng fetch giữ nguyên cơ chế auth bằng header
 * (không phải đưa token lên query string nên không lộ token vào access log /
 * lịch sử trình duyệt / Referer), đổi lại phải tự đọc stream và tự kết nối lại.
 *
 * Class thuần, không gọi composable của Nuxt → test được bằng node:test.
 */

export interface NotificationStreamHandlers {
  /** Nhận payload đã parse từ 1 frame `data:` */
  onMessage: (payload: unknown) => void
  /** Kết nối đã mở */
  onOpen: () => void
  /** Kết nối đã đóng (trước khi chờ kết nối lại) */
  onClose: () => void
}

export interface NotificationStreamOptions {
  /** URL đầy đủ của endpoint SSE */
  url: string
  /** Lấy access token hiện tại (null = chưa đăng nhập) */
  getToken: () => string | null
  /** Làm mới token khi gặp 401; trả false nếu không làm mới được */
  refreshToken: () => Promise<boolean>
  handlers: NotificationStreamHandlers
  /** Thời gian chờ trước khi kết nối lại (ms) */
  reconnectDelayMs?: number
  /** Cho phép bơm fetch khác khi test */
  fetchImpl?: typeof fetch
}

const DEFAULT_RECONNECT_DELAY_MS = 5000

export class NotificationStream {
  private controller: AbortController | null = null
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null
  private running = false

  private readonly options: NotificationStreamOptions
  private readonly fetchImpl: typeof fetch
  private readonly reconnectDelayMs: number

  // Gán tường minh (không dùng parameter property) để file này chạy được cả với
  // chế độ xoá type của Node khi test — parameter property không phải cú pháp xoá được.
  constructor(options: NotificationStreamOptions) {
    this.options = options
    this.fetchImpl = options.fetchImpl ?? fetch.bind(globalThis)
    this.reconnectDelayMs = options.reconnectDelayMs ?? DEFAULT_RECONNECT_DELAY_MS
  }

  /** Bắt đầu kết nối. Gọi nhiều lần cũng không tạo kết nối trùng. */
  start(): void {
    if (this.running) return
    this.running = true
    void this.connect()
  }

  /** Ngắt kết nối và dừng mọi lần kết nối lại */
  stop(): void {
    if (!this.running && !this.controller) {
      this.options.handlers.onClose()
      return
    }
    this.running = false
    this.clearReconnectTimer()
    this.controller?.abort()
    this.controller = null
    this.options.handlers.onClose()
  }

  get isRunning(): boolean {
    return this.running
  }

  private async connect(): Promise<void> {
    if (!this.running) return

    const token = this.options.getToken()
    if (!token) {
      // Chưa đăng nhập thì không kết nối (khách vãng lai không gọi API thông báo)
      this.stop()
      return
    }

    this.controller = new AbortController()

    try {
      const response = await this.fetchImpl(this.options.url, {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${token}`,
          Accept: 'text/event-stream',
        },
        credentials: 'include',
        signal: this.controller.signal,
      })

      // Token hết hạn → làm mới 1 lần rồi kết nối lại; không làm mới được thì dừng hẳn
      if (response.status === 401) {
        const refreshed = await this.options.refreshToken()
        if (refreshed) {
          this.scheduleReconnect()
        } else {
          this.stop()
        }
        return
      }

      if (!response.ok || !response.body) {
        throw new Error(`Kết nối thông báo thất bại (HTTP ${response.status})`)
      }

      this.options.handlers.onOpen()
      await this.readStream(response.body)
      // Server đóng stream → kết nối lại
      this.scheduleReconnect()
    } catch (error) {
      // stop() chủ động abort thì không phải lỗi, cũng không kết nối lại
      if (error instanceof Error && error.name === 'AbortError') return
      console.error('SSE error:', error)
      this.scheduleReconnect()
    }
  }

  /** Đọc stream và cắt frame theo dòng trống (chuẩn SSE) */
  private async readStream(body: ReadableStream<Uint8Array>): Promise<void> {
    const reader = body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    try {
      while (this.running) {
        const { done, value } = await reader.read()
        if (done) break

        // stream: true để ký tự UTF-8 bị cắt giữa 2 chunk vẫn ghép đúng
        buffer += decoder.decode(value, { stream: true })

        let boundary = buffer.indexOf('\n\n')
        while (boundary !== -1) {
          this.handleFrame(buffer.slice(0, boundary))
          buffer = buffer.slice(boundary + 2)
          boundary = buffer.indexOf('\n\n')
        }
      }
    } finally {
      reader.releaseLock()
    }
  }

  /** Xử lý 1 frame: chỉ lấy các dòng `data:`, bỏ qua `retry:`/comment */
  private handleFrame(frame: string): void {
    const data = frame
      .split('\n')
      .filter((line) => line.startsWith('data:'))
      .map((line) => line.slice('data:'.length).trimStart())
      .join('\n')

    if (!data) return

    try {
      this.options.handlers.onMessage(JSON.parse(data))
    } catch (error) {
      console.error('SSE parse error:', error)
    }
  }

  private scheduleReconnect(): void {
    this.controller = null
    if (!this.running) return

    this.options.handlers.onClose()
    this.clearReconnectTimer()
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null
      void this.connect()
    }, this.reconnectDelayMs)
  }

  private clearReconnectTimer(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
  }
}
