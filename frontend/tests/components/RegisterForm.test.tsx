import "@testing-library/jest-dom";
import { act, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import RegisterForm from "@/components/auth/RegisterForm";

const pushMock = jest.fn();
const registerMock = jest.fn().mockResolvedValue({ id: 1, username: "tester" });
const messageMock = (jest.requireMock("antd") as any).message as {
  success: jest.Mock;
  error: jest.Mock;
};

jest.mock("next/navigation", () => ({
  useRouter: () => ({ push: pushMock, replace: jest.fn() }),
}));

jest.mock("@/hooks/useAuth", () => ({
  __esModule: true,
  default: () => ({
    register: registerMock,
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

describe("RegisterForm", () => {
  beforeEach(() => {
    pushMock.mockReset();
    registerMock.mockClear();
    messageMock.success.mockClear();
    messageMock.error.mockClear();
  });

  it("提交注册表单后调用 register 并跳转", async () => {
    render(<RegisterForm />);

    await act(async () => {
      await userEvent.type(screen.getByLabelText("用户名"), "tester");
      await userEvent.type(screen.getByLabelText("密码"), "Password123!");
      await userEvent.type(screen.getByLabelText("确认密码"), "Password123!");
      await userEvent.click(screen.getByRole("button", { name: /注册并登录/ }));
    });

    await waitFor(() => {
      expect(registerMock).toHaveBeenCalledWith("tester", "Password123!", true);
    });
    expect(messageMock.success).toHaveBeenCalledWith("注册成功，已创建新用户");
    expect(pushMock).not.toHaveBeenCalled();
  });

  it("弱口令时阻止提交并提示策略要求", async () => {
    render(<RegisterForm />);

    await act(async () => {
      await userEvent.type(screen.getByLabelText("用户名"), "tester");
      await userEvent.type(screen.getByLabelText("密码"), "Password123");
      await userEvent.type(screen.getByLabelText("确认密码"), "Password123");
      await userEvent.click(screen.getByRole("button", { name: /注册并登录/ }));
    });

    await waitFor(() => {
      expect(
        screen.getByText("需包含大写、小写、数字和特殊字符"),
      ).toBeInTheDocument();
    });
    expect(registerMock).not.toHaveBeenCalled();
    expect(pushMock).not.toHaveBeenCalled();
  });
});
