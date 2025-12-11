import "@testing-library/jest-dom";
import api from "@/lib/api";

jest.mock("antd", () => ({
  message: {
    error: jest.fn(),
  },
}));

jest.mock("@/lib/auth", () => ({
  getAccessToken: jest.fn(() => null),
  refreshAccessToken: jest.fn(),
  clearTokens: jest.fn(),
}));

const authMock = jest.requireMock("@/lib/auth");

describe("api 拦截器行为", () => {
  it("业务成功时返回 data", async () => {
    const handler = (api as any).interceptors.response.handlers[0].fulfilled;
    const res = await handler({ data: { code: 20000, data: { ok: true } } });
    expect(res).toEqual({ ok: true });
  });

  it("业务错误时抛出异常", async () => {
    const handler = (api as any).interceptors.response.handlers[0].fulfilled;
    await expect(
      handler({ data: { code: 40004, message: "bad" } }),
    ).rejects.toThrow("bad");
  });

  it("401 时尝试刷新并重试原请求", async () => {
    const handler = (api as any).interceptors.response.handlers[0].rejected;
    (authMock.refreshAccessToken as jest.Mock).mockResolvedValue("new-token");
    const adapter = jest.fn(async (config) => ({
      data: { code: 20000, data: "retried" },
      status: 200,
      statusText: "OK",
      headers: {},
      config,
    }));
    const error = {
      response: { status: 401 },
      config: { headers: {}, adapter },
    };
    const result = await handler(error);
    expect(authMock.refreshAccessToken).toHaveBeenCalled();
    expect(adapter).toHaveBeenCalled();
    expect(result).toBe("retried");
  });
});
