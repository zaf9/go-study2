import "@testing-library/jest-dom";
import { act, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import LoginForm from "@/components/auth/LoginForm";

jest.setTimeout(15000);

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

describe("LoginForm", () => {
  beforeEach(() => {
    pushMock.mockReset();
    loginMock.mockClear();
  });

  it("提交登录表单时调用 login 并跳转", async () => {
    render(<LoginForm />);

    await act(async () => {
      await userEvent.type(
        screen.getByPlaceholderText("请输入用户名"),
        "tester",
      );
      await userEvent.type(screen.getByPlaceholderText("请输入密码"), "Password123!");
      await userEvent.click(screen.getByRole("button", { name: /登\s*录/ }));
    });

    await waitFor(() => {
      expect(loginMock).toHaveBeenCalledWith("tester", "Password123!", true);
    });
    expect(pushMock).toHaveBeenCalledWith("/topics");
  });

  it("弱口令时展示校验提示并不触发登录", async () => {
    render(<LoginForm />);

    await act(async () => {
      await userEvent.type(
        screen.getByPlaceholderText("请输入用户名"),
        "tester",
      );
      await userEvent.type(screen.getByPlaceholderText("请输入密码"), "Password123");
      await userEvent.click(screen.getByRole("button", { name: /登\s*录/ }));
    });

    await waitFor(() => {
      expect(
        screen.getByText("需包含大写、小写、数字和特殊字符"),
      ).toBeInTheDocument();
    });
    expect(loginMock).not.toHaveBeenCalled();
    expect(pushMock).not.toHaveBeenCalled();
  });
});
