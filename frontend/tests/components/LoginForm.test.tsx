import "@testing-library/jest-dom";
import { act, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import LoginForm from "@/components/auth/LoginForm";

const pushMock = jest.fn();
const loginMock = jest.fn().mockResolvedValue({ id: 1, username: "tester" });

jest.mock("next/navigation", () => ({
  useRouter: () => ({ push: pushMock, replace: jest.fn() }),
}));

jest.mock("@/hooks/useAuth", () => ({
  __esModule: true,
  default: () => ({
    login: loginMock,
  }),
}));

jest.mock("antd", () => {
  const actual = jest.requireActual("antd");
  return {
    ...actual,
    message: {
      success: jest.fn(),
      error: jest.fn(),
      warning: jest.fn(),
      info: jest.fn(),
      open: jest.fn(),
      destroy: jest.fn(),
    },
  };
});

describe("LoginForm", () => {
  beforeEach(() => {
    pushMock.mockReset();
    loginMock.mockClear();
  });

  it("提交登录表单时调用 login 并跳转", async () => {
    render(<LoginForm />);

    await act(async () => {
      await userEvent.type(screen.getByLabelText("用户名"), "tester");
      await userEvent.type(screen.getByLabelText("密码"), "Password123");
      await userEvent.click(screen.getByRole("button", { name: /登\s*录/ }));
    });

    await waitFor(() => {
      expect(loginMock).toHaveBeenCalledWith("tester", "Password123", true);
    });
    expect(pushMock).toHaveBeenCalledWith("/topics");
  });
});


