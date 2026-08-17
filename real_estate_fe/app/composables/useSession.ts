/**
 * Sinh UUID v4 đơn giản không cần thư viện ngoài
 */
function generateUUID(): string {
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

const getStoredSessionId = (): string => {
  if (import.meta.client) {
    let id = localStorage.getItem("real_estate_session_id");
    if (!id) {
      id = generateUUID();
      localStorage.setItem("real_estate_session_id", id);
    }
    return id;
  }
  return "";
};

const sessionId = ref<string>(import.meta.client ? getStoredSessionId() : "");

export const useSession = () => {
  return {
    sessionId,
    // getSessionId: getStoredSessionId,
  };
};
