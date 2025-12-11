import "@testing-library/jest-dom";

const TOKEN_KEY = "go-study2.access_token";

function mockAxios() {
  const post = jest.fn();
  const get = jest.fn();
  const create = jest.fn(() => ({ post, get }));
  jest.doMock("axios", () => ({
    __esModule: true,
    default: Object.assign(create, { create, post, get }),
  }));
  return { post, get, create };
}

describe("auth storage helpers", () => {
  beforeEach(() => {
    jest.resetModules();
    localStorage.clear();
  });

  it("记住我为 true 时会持久化 token", async () => {
    mockAxios();
    const { setRememberMe, setAccessToken, getAccessToken, clearTokens } =
      await import("@/lib/auth");

    setRememberMe(true);
    setAccessToken("abc");

    expect(localStorage.getItem(TOKEN_KEY)).toBe("abc");
    expect(getAccessToken()).toBe("abc");

    clearTokens();
  });

  it("未勾选记住我时仅使用内存中的 token", async () => {
    mockAxios();
    const { setRememberMe, setAccessToken, getAccessToken, clearTokens } =
      await import("@/lib/auth");

    setRememberMe(false);
    setAccessToken("temp-token");

    expect(localStorage.getItem(TOKEN_KEY)).toBeNull();
    expect(getAccessToken()).toBe("temp-token");

    clearTokens();
    expect(getAccessToken()).toBeNull();
  });

  it("刷新 access token 成功时更新缓存", async () => {
    const axiosMock = mockAxios();
    const { setRememberMe, refreshAccessToken, getAccessToken, clearTokens } =
      await import("@/lib/auth");

    setRememberMe(true);
    axiosMock.post.mockResolvedValue({
      data: {
        code: 20000,
        message: "ok",
        data: { accessToken: "new-token", expiresIn: 3600 },
      },
    });

    const token = await refreshAccessToken();

    expect(axiosMock.post).toHaveBeenCalled();
    expect(token).toBe("new-token");
    expect(getAccessToken()).toBe("new-token");

    clearTokens();
  });

  it("刷新失败时返回 null 且清理 token", async () => {
    const axiosMock = mockAxios();
    const {
      setRememberMe,
      setAccessToken,
      refreshAccessToken,
      getAccessToken,
      clearTokens,
    } = await import("@/lib/auth");

    setRememberMe(true);
    setAccessToken("old-token");
    axiosMock.post.mockRejectedValue(new Error("network"));

    const token = await refreshAccessToken();

    expect(axiosMock.post).toHaveBeenCalled();
    expect(token).toBeNull();
    expect(getAccessToken()).toBeNull();
    expect(localStorage.getItem(TOKEN_KEY)).toBeNull();

    clearTokens();
  });
});
