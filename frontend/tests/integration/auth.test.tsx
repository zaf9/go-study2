import "@testing-library/jest-dom";
import React, { useContext, useEffect } from "react";
import { render, waitFor, act } from "@testing-library/react";
import { AuthContext, AuthProvider } from "@/contexts/AuthContext";
import AuthGuard from "@/components/auth/AuthGuard";

jest.mock("@/lib/auth", () => ({
  loginWithPassword: jest.fn(),
  registerAccount: jest.fn(),
  fetchProfile: jest.fn(),
  logoutAccount: jest.fn(),
  refreshAccessToken: jest.fn(),
  getAccessToken: jest.fn(() => null),
  clearTokens: jest.fn(),
  setAccessToken: jest.fn(),
  setRememberMe: jest.fn(),
  changePassword: jest.fn(),
}));

jest.mock("next/navigation", () => ({
  useRouter: jest.fn(),
  usePathname: jest.fn(),
}));

const authMock = jest.requireMock("@/lib/auth");
const navigationMock = jest.requireMock("next/navigation");

function renderWithAuth(consumer: (ctx: any) => void) {
  const Wrapper = () => {
    const ctx = useContext(AuthContext)!;
    useEffect(() => consumer(ctx), [ctx]);
    return null;
  };
  render(
    <AuthProvider>
      <Wrapper />
    </AuthProvider>,
  );
}

describe("AuthProvider 集成", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("login 会调用后端并填充用户信息", async () => {
    const profile = {
      id: 1,
      username: "tester",
      isAdmin: false,
      mustChangePassword: false,
    };
    (authMock.loginWithPassword as jest.Mock).mockResolvedValue({
      accessToken: "token",
    });
    (authMock.fetchProfile as jest.Mock).mockResolvedValue(profile);

    let latest: any;
    renderWithAuth((ctx) => {
      latest = ctx;
    });

    await act(async () => {
      await latest.login("tester", "Passw0rd", true);
    });

    await waitFor(() => expect(latest.user?.username).toBe("tester"));
    expect(authMock.loginWithPassword).toHaveBeenCalledWith(
      expect.objectContaining({
        username: "tester",
        password: "Passw0rd",
        remember: true,
      }),
    );
  });

  it("logout 会清空用户状态", async () => {
    const profile = {
      id: 2,
      username: "logout_user",
      isAdmin: false,
      mustChangePassword: false,
    };
    (authMock.loginWithPassword as jest.Mock).mockResolvedValue({
      accessToken: "token2",
    });
    (authMock.fetchProfile as jest.Mock).mockResolvedValue(profile);

    let latest: any;
    renderWithAuth((ctx) => {
      latest = ctx;
    });

    await act(async () => {
      await latest.login("logout_user", "Passw0rd", true);
    });
    await waitFor(() => expect(latest.user?.username).toBe("logout_user"));

    await act(async () => {
      await latest.logout();
    });
    await waitFor(() => expect(latest.user).toBeNull());
    expect(authMock.logoutAccount).toHaveBeenCalled();
  });

  it("AuthGuard 会在需改密时重定向改密页", async () => {
    const replace = jest.fn();
    (navigationMock.useRouter as jest.Mock).mockReturnValue({
      replace,
    });
    (navigationMock.usePathname as jest.Mock).mockReturnValue("/topics");

    render(
      <AuthContext.Provider
        value={{
          user: {
            id: 1,
            username: "admin",
            isAdmin: true,
            mustChangePassword: true,
          },
          loading: false,
          login: jest.fn(),
          register: jest.fn(),
          changePassword: jest.fn(),
          logout: jest.fn(),
          refreshProfile: jest.fn(),
        }}
      >
        <AuthGuard>
          <div>protected</div>
        </AuthGuard>
      </AuthContext.Provider>,
    );

    await waitFor(() =>
      expect(replace).toHaveBeenCalledWith("/change-password"),
    );
  });
});
