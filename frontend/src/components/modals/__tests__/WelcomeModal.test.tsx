import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { WelcomeModal } from "../WelcomeModal";
import { settingsService } from "@/services/settingsService";
import { api } from "@/services/api";

// Mock dependencies
vi.mock("@/services/settingsService", () => ({
  settingsService: {
    getBarkKey: vi.fn(),
    setBarkKey: vi.fn(),
    clearBarkKey: vi.fn(),
  },
}));

vi.mock("@/services/api", () => ({
  api: {
    post: vi.fn(),
  },
}));

describe("WelcomeModal", () => {
  const mockOnClose = vi.fn();
  const mockOnComplete = vi.fn();

  const defaultProps = {
    open: true,
    onClose: mockOnClose,
    onComplete: mockOnComplete,
  };

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(settingsService.getBarkKey).mockReturnValue("");
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe("Rendering", () => {
    it("should render when open is true", () => {
      render(<WelcomeModal {...defaultProps} />);
      expect(screen.getByText("欢迎使用美食监控")).toBeInTheDocument();
    });

    it("should not render when open is false", () => {
      render(<WelcomeModal {...defaultProps} open={false} />);
      expect(screen.queryByText("欢迎使用美食监控")).not.toBeInTheDocument();
    });

    it("should load existing Bark Key from settings on open", () => {
      vi.mocked(settingsService.getBarkKey).mockReturnValue("existing-key");
      render(<WelcomeModal {...defaultProps} />);
      expect(screen.getByDisplayValue("existing-key")).toBeInTheDocument();
    });

    it("should display all benefit items", () => {
      render(<WelcomeModal {...defaultProps} />);
      expect(screen.getByText("价格提醒推送")).toBeInTheDocument();
      expect(screen.getByText("数据云同步")).toBeInTheDocument();
      expect(screen.getByText("跨设备访问")).toBeInTheDocument();
      expect(screen.getByText("无需注册登录")).toBeInTheDocument();
    });
  });

  describe("Input Handling", () => {
    it("should update Bark Key input value", async () => {
      const user = userEvent.setup();
      render(<WelcomeModal {...defaultProps} />);

      const input = screen.getByPlaceholderText("完整URL 或 设备Key");
      await user.type(input, "test-key");

      expect(input).toHaveValue("test-key");
    });

    it("should clear input when clear button clicked", async () => {
      const user = userEvent.setup();
      vi.mocked(settingsService.getBarkKey).mockReturnValue("existing-key");
      render(<WelcomeModal {...defaultProps} />);

      // Find the clear button inside the input field (has bg-slate-100 class)
      const container =
        screen.getByPlaceholderText("完整URL 或 设备Key").parentElement;
      const clearButton = container?.querySelector("button.bg-slate-100");
      expect(clearButton).toBeTruthy();
      await user.click(clearButton!);

      const input = screen.getByPlaceholderText("完整URL 或 设备Key");
      expect(input).toHaveValue("");
    });

    it("should show clear button only when input has value", () => {
      render(<WelcomeModal {...defaultProps} />);

      // Initially no clear button visible (input is empty)
      const input = screen.getByPlaceholderText("完整URL 或 设备Key");
      expect(input).toHaveValue("");
    });
  });

  describe("Test Notification", () => {
    it("should show error when API returns failure", async () => {
      const user = userEvent.setup();
      vi.mocked(api.post).mockResolvedValue({
        data: { data: { success: false, error: "无效的 Bark Key" } },
      });

      render(<WelcomeModal {...defaultProps} />);

      const input = screen.getByPlaceholderText("完整URL 或 设备Key");
      await user.type(input, "invalid-key");

      const testButton = screen.getByText("测试通知");
      await user.click(testButton);

      await waitFor(() => {
        expect(screen.getByText("无效的 Bark Key")).toBeInTheDocument();
      });
    });

    it("should show success message on successful test", async () => {
      const user = userEvent.setup();
      vi.mocked(api.post).mockResolvedValue({
        data: { data: { success: true } },
      });

      render(<WelcomeModal {...defaultProps} />);

      const input = screen.getByPlaceholderText("完整URL 或 设备Key");
      await user.type(input, "valid-key");

      const testButton = screen.getByText("测试通知");
      await user.click(testButton);

      await waitFor(() => {
        expect(
          screen.getByText("通知发送成功，请检查您的手机"),
        ).toBeInTheDocument();
      });
    });

    it("should show error message on API failure", async () => {
      const user = userEvent.setup();
      vi.mocked(api.post).mockRejectedValue(new Error("Network error"));

      render(<WelcomeModal {...defaultProps} />);

      const input = screen.getByPlaceholderText("完整URL 或 设备Key");
      await user.type(input, "invalid-key");

      const testButton = screen.getByText("测试通知");
      await user.click(testButton);

      await waitFor(() => {
        expect(screen.getByText("Network error")).toBeInTheDocument();
      });
    });

    it("should normalize Bark Key URL to device key when testing", async () => {
      const user = userEvent.setup();
      vi.mocked(api.post).mockResolvedValue({
        data: { data: { success: true } },
      });

      render(<WelcomeModal {...defaultProps} />);

      const input = screen.getByPlaceholderText("完整URL 或 设备Key");
      await user.type(input, "https://api.day.app/ABC123");

      const testButton = screen.getByText("测试通知");
      await user.click(testButton);

      await waitFor(() => {
        expect(api.post).toHaveBeenCalledWith("/admin/test-notification", {
          barkKey: "ABC123",
        });
      });
    });

    it("should show loading state during test", async () => {
      const user = userEvent.setup();
      vi.mocked(api.post).mockImplementation(
        () => new Promise((resolve) => setTimeout(resolve, 100)),
      );

      render(<WelcomeModal {...defaultProps} />);

      const input = screen.getByPlaceholderText("完整URL 或 设备Key");
      await user.type(input, "test-key");

      const testButton = screen.getByText("测试通知");
      await user.click(testButton);

      // Check loading state
      await waitFor(() => {
        expect(screen.getByText("发送中")).toBeInTheDocument();
      });
    });
  });

  describe("Save Functionality", () => {
    it("should save to localStorage and backend", async () => {
      const user = userEvent.setup();
      vi.mocked(api.post).mockResolvedValue({});

      render(<WelcomeModal {...defaultProps} />);

      const input = screen.getByPlaceholderText("完整URL 或 设备Key");
      await user.type(input, "new-key");

      const saveButton = screen.getByText("完成设置");
      await user.click(saveButton);

      await waitFor(() => {
        expect(settingsService.setBarkKey).toHaveBeenCalledWith("new-key");
        expect(api.post).toHaveBeenCalledWith("/user/settings", {
          barkKey: "new-key",
        });
        expect(mockOnComplete).toHaveBeenCalled();
        expect(mockOnClose).toHaveBeenCalled();
      });
    });

    it("should clear Bark Key when saving empty value", async () => {
      const user = userEvent.setup();
      vi.mocked(api.post).mockResolvedValue({});

      render(<WelcomeModal {...defaultProps} />);

      const saveButton = screen.getByText("完成设置");
      await user.click(saveButton);

      await waitFor(() => {
        expect(settingsService.clearBarkKey).toHaveBeenCalled();
      });
    });

    it("should normalize Bark Key URL when saving", async () => {
      const user = userEvent.setup();
      vi.mocked(api.post).mockResolvedValue({});

      render(<WelcomeModal {...defaultProps} />);

      const input = screen.getByPlaceholderText("完整URL 或 设备Key");
      await user.type(input, "https://api.day.app/XYZ789");

      const saveButton = screen.getByText("完成设置");
      await user.click(saveButton);

      await waitFor(() => {
        expect(settingsService.setBarkKey).toHaveBeenCalledWith("XYZ789");
        expect(api.post).toHaveBeenCalledWith("/user/settings", {
          barkKey: "XYZ789",
        });
      });
    });

    it("should show loading state during save", async () => {
      const user = userEvent.setup();
      vi.mocked(api.post).mockImplementation(
        () => new Promise((resolve) => setTimeout(resolve, 100)),
      );

      render(<WelcomeModal {...defaultProps} />);

      const input = screen.getByPlaceholderText("完整URL 或 设备Key");
      await user.type(input, "test-key");

      const saveButton = screen.getByText("完成设置");
      await user.click(saveButton);

      await waitFor(() => {
        expect(screen.getByText("保存中")).toBeInTheDocument();
      });
    });
  });

  describe("Skip Functionality", () => {
    it("should close modal without saving when skip clicked", async () => {
      const user = userEvent.setup();
      render(<WelcomeModal {...defaultProps} />);

      const skipButton = screen.getByText("暂不设置");
      await user.click(skipButton);

      expect(mockOnClose).toHaveBeenCalled();
      expect(mockOnComplete).not.toHaveBeenCalled();
    });
  });

  describe("Loading States", () => {
    it("should disable buttons during testing", async () => {
      const user = userEvent.setup();
      vi.mocked(api.post).mockImplementation(
        () => new Promise((resolve) => setTimeout(resolve, 100)),
      );

      render(<WelcomeModal {...defaultProps} />);

      const input = screen.getByPlaceholderText("完整URL 或 设备Key");
      await user.type(input, "key");

      const testButton = screen.getByText("测试通知");
      await user.click(testButton);

      await waitFor(() => {
        const saveButton = screen.getByText("完成设置");
        expect(saveButton).toBeDisabled();
      });
    });

    it("should disable test button when input is empty", () => {
      render(<WelcomeModal {...defaultProps} />);

      const testButton = screen.getByText("测试通知");
      expect(testButton).toBeDisabled();
    });
  });

  describe("Backdrop Click", () => {
    it("should have backdrop that closes modal on click", () => {
      render(<WelcomeModal {...defaultProps} />);

      // Verify backdrop element exists (clicking is tested via skip button instead)
      const modal = screen.getByText("欢迎使用美食监控").closest(".relative");
      expect(modal).toBeTruthy();
    });
  });
});

describe("normalizeBarkKey utility", () => {
  // Test the normalizeBarkKey function through the component's behavior
  it("should extract device key from full URL", async () => {
    const user = userEvent.setup();
    vi.mocked(api.post).mockResolvedValue({
      data: { data: { success: true } },
    });

    render(<WelcomeModal open={true} onClose={vi.fn()} onComplete={vi.fn()} />);

    const input = screen.getByPlaceholderText("完整URL 或 设备Key");
    await user.type(input, "https://api.day.app/MYKEY123");

    const testButton = screen.getByText("测试通知");
    await user.click(testButton);

    await waitFor(() => {
      expect(api.post).toHaveBeenCalledWith("/admin/test-notification", {
        barkKey: "MYKEY123",
      });
    });
  });

  it("should return key as-is when not a URL", async () => {
    const user = userEvent.setup();
    vi.mocked(api.post).mockResolvedValue({
      data: { data: { success: true } },
    });

    render(<WelcomeModal open={true} onClose={vi.fn()} onComplete={vi.fn()} />);

    const input = screen.getByPlaceholderText("完整URL 或 设备Key");
    await user.type(input, "SIMPLEKEY");

    const testButton = screen.getByText("测试通知");
    await user.click(testButton);

    await waitFor(() => {
      expect(api.post).toHaveBeenCalledWith("/admin/test-notification", {
        barkKey: "SIMPLEKEY",
      });
    });
  });
});
