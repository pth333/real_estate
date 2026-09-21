import test from "node:test";
import assert from "node:assert/strict";

import { NotificationStream } from "../app/services/notification-stream.ts";

const encoder = new TextEncoder();
const STREAM_URL = "http://localhost:8000/api/2026/notifications/stream";

/** Tạo ReadableStream từ danh sách chunk (string hoặc Uint8Array), giữ controller để đóng thủ công */
function makeStream(chunks, { autoClose = true } = {}) {
  let controller;
  const body = new ReadableStream({
    start(c) {
      controller = c;
      // ReadableStream chỉ chứa byte — string phải encode trước, nếu không
      // TextDecoder.decode() ở phía client sẽ throw ERR_INVALID_ARG_TYPE
      for (const chunk of chunks) {
        c.enqueue(typeof chunk === "string" ? encoder.encode(chunk) : chunk);
      }
      if (autoClose) c.close();
    },
  });
  return { body, close: () => controller?.close() };
}

/** Các client đang chạy — đảm bảo mọi test đều dừng lại, tránh vòng lặp kết nối rò rỉ */
const activeClients = new Set();

test.afterEach(() => {
  for (const client of activeClients) client.stream.stop();
  activeClients.clear();
});

/** Tạo client + ghi lại mọi thứ cần assert */
function makeClient({ responses = [], token = "token-abc", refreshResult = true } = {}) {
  const received = [];
  const events = [];
  const fetchCalls = [];
  const refreshCalls = { count: 0 };

  const stream = new NotificationStream({
    url: STREAM_URL,
    getToken: () => token,
    refreshToken: async () => {
      refreshCalls.count++;
      return refreshResult;
    },
    reconnectDelayMs: 10,
    handlers: {
      onMessage: (payload) => {
        received.push(payload);
        events.push("message");
      },
      onOpen: () => events.push("open"),
      onClose: () => events.push("close"),
    },
    fetchImpl: async (url, init) => {
      fetchCalls.push({ url, authorization: init?.headers?.Authorization });
      // responses[i] là hàm trả về response cho lần gọi thứ i (1-based)
      const buildResponse = responses[fetchCalls.length - 1];
      return buildResponse ? buildResponse() : { ok: false, status: 500, body: null };
    },
  });

  const client = { stream, received, events, fetchCalls, refreshCalls };
  activeClients.add(client);
  return client;
}

function okChunks(chunks) {
  return () => {
    const { body } = makeStream(chunks);
    return { ok: true, status: 200, body };
  };
}

function unauthorized() {
  return () => ({ ok: false, status: 401, body: null });
}

function waitFor(predicate, { timeoutMs = 1500 } = {}) {
  const started = Date.now();
  return new Promise((resolve, reject) => {
    const tick = () => {
      if (predicate()) return resolve();
      if (Date.now() - started > timeoutMs) return reject(new Error("waitFor timeout"));
      setTimeout(tick, 5);
    };
    tick();
  });
}

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

test("gửi Authorization header — điều EventSource không làm được", async () => {
  const client = makeClient({ responses: [okChunks([])] });
  client.stream.start();
  await waitFor(() => client.fetchCalls.length === 1);

  assert.equal(client.fetchCalls[0].url, STREAM_URL);
  assert.equal(client.fetchCalls[0].authorization, "Bearer token-abc");
  client.stream.stop();
});

test("parse frame data và bỏ qua dòng retry:", async () => {
  const payload = { title: "Nhà 3 tầng", price: 3_500_000_000, slug: "nha-3-tang-rs1" };
  const client = makeClient({
    responses: [okChunks([`retry: 5000\n\n`, `data: ${JSON.stringify(payload)}\n\n`])],
  });

  client.stream.start();
  await waitFor(() => client.received.length === 1);

  assert.deepEqual(client.received[0], payload);
  assert.ok(client.events.includes("open"), "phải báo onOpen khi kết nối thành công");
  client.stream.stop();
});

test("ghép đúng frame khi chunk bị cắt từng byte (kể cả giữa ký tự UTF-8)", async () => {
  const payload = { title: "Căn hộ cao cấp Trung Hoà", price: 2_000_000_000, slug: "can-ho-rs2" };
  const bytes = encoder.encode(`data: ${JSON.stringify(payload)}\n\n`);
  const chunks = [];
  for (const byte of bytes) chunks.push(new Uint8Array([byte]));

  const client = makeClient({ responses: [okChunks(chunks)] });
  client.stream.start();
  await waitFor(() => client.received.length === 1);

  assert.deepEqual(client.received[0], payload, "frame bị cắt nhỏ vẫn phải parse đúng");
  client.stream.stop();
});

test("nhiều frame trong cùng 1 chunk đều được xử lý", async () => {
  const one = { title: "Tin 1", price: 1_000_000_000, slug: "tin-1" };
  const two = { title: "Tin 2", price: 2_000_000_000, slug: "tin-2" };
  const client = makeClient({
    responses: [okChunks([`data: ${JSON.stringify(one)}\n\ndata: ${JSON.stringify(two)}\n\n`])],
  });

  client.stream.start();
  await waitFor(() => client.received.length === 2);

  assert.deepEqual(client.received, [one, two]);
  client.stream.stop();
});

test("JSON hỏng không làm chết vòng đọc", async () => {
  const good = { title: "Tin tốt", price: 1_000_000_000, slug: "tin-good" };
  const client = makeClient({
    responses: [okChunks([`data: {khong-phai-json\n\n`, `data: ${JSON.stringify(good)}\n\n`])],
  });

  client.stream.start();
  await waitFor(() => client.received.length === 1);

  assert.deepEqual(client.received[0], good, "frame sau frame lỗi vẫn phải nhận được");
  client.stream.stop();
});

test("401 → làm mới token rồi kết nối lại thành công", async () => {
  const payload = { title: "Sau khi refresh", price: 1_500_000_000, slug: "sau-refresh" };
  const client = makeClient({
    responses: [unauthorized(), okChunks([`data: ${JSON.stringify(payload)}\n\n`])],
  });

  client.stream.start();
  await waitFor(() => client.received.length === 1);

  assert.deepEqual(client.received[0], payload);
  assert.equal(client.refreshCalls.count, 1, "phải gọi làm mới token đúng 1 lần");
  assert.equal(client.fetchCalls.length, 2, "phải thử kết nối lại đúng 1 lần");
  client.stream.stop();
});

test("401 và không làm mới được token → dừng hẳn, không lặp vô hạn", async () => {
  const client = makeClient({ responses: [unauthorized()], refreshResult: false });

  client.stream.start();
  await waitFor(() => client.fetchCalls.length === 1);
  await sleep(80);

  assert.equal(client.fetchCalls.length, 1, "không được kết nối lại khi refresh thất bại");
  assert.equal(client.stream.isRunning, false, "phải dừng hẳn");
});

test("chưa đăng nhập thì không gọi API thông báo", async () => {
  const client = makeClient({ token: null });

  client.stream.start();
  await sleep(40);

  assert.equal(client.fetchCalls.length, 0, "khách vãng lai không được gọi endpoint thông báo");
  assert.equal(client.stream.isRunning, false);
  assert.ok(client.events.includes("close"), "phải báo đóng kết nối");
});

test("stop() ngắt kết nối và không kết nối lại", async () => {
  // Stream chỉ đóng khi gọi close() → mô phỏng kết nối SSE đang mở
  const { body, close } = makeStream([], { autoClose: false });
  const client = makeClient({
    responses: [() => ({ ok: true, status: 200, body })],
  });

  client.stream.start();
  await waitFor(() => client.fetchCalls.length === 1);
  assert.equal(client.stream.isRunning, true);

  client.stream.stop();
  // Đóng stream để không để lại read() treo (chỉ ảnh hưởng tiến trình test)
  close();
  await sleep(80);

  assert.equal(client.stream.isRunning, false);
  assert.equal(client.fetchCalls.length, 1, "sau stop() không được kết nối lại");
});

test("server đóng stream → tự kết nối lại", async () => {
  const payload = { title: "Lần 2", price: 900_000_000, slug: "lan-2" };
  const client = makeClient({
    responses: [okChunks([]), okChunks([`data: ${JSON.stringify(payload)}\n\n`])],
  });

  client.stream.start();
  await waitFor(() => client.received.length === 1);

  assert.deepEqual(client.received[0], payload);
  assert.ok(client.fetchCalls.length >= 2, "stream đóng phải kích hoạt kết nối lại");
  client.stream.stop();
});
