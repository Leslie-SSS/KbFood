import { useState, useEffect } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Navbar } from "@/components/layout/Navbar";
import { SearchBar } from "@/components/search/SearchBar";
import { ProductFilters } from "@/components/products/ProductFilters";
import { ProductGrid } from "@/components/products/ProductGrid";
import { NotificationModal } from "@/components/modals/NotificationModal";
import { BlockedProductsModal } from "@/components/modals/BlockedProductsModal";
import { SettingsModal } from "@/components/modals/SettingsModal";
import { WelcomeModal } from "@/components/modals/WelcomeModal";
import { ConfirmModal } from "@/components/modals/ConfirmModal";
import { Toast, useToast, showToast } from "@/components/ui/Toast";
import { useProducts, useBlockedProducts, useSystemStatus } from "@/hooks";
import { settingsService } from "@/services/settingsService";
import type { Product, ProductFilters as TProductFilters } from "@/types";

// Create QueryClient instance
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

function AppContent() {
  // Filter state
  const [filters, setFilters] = useState<TProductFilters>({
    keyword: "",
    platform: "",
    region: "",
    salesStatus: "",
    monitorStatus: "",
    recentSevenDays: false,
  });

  // UI state
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  const [isBlockedOpen, setIsBlockedOpen] = useState(false);
  const [isWelcomeOpen, setIsWelcomeOpen] = useState(false);
  const [notificationModal, setNotificationModal] = useState<{
    open: boolean;
    product: Product | null;
    isEdit: boolean;
  }>({ open: false, product: null, isEdit: false });

  // Confirm modal state
  const [confirmModal, setConfirmModal] = useState<{
    open: boolean;
    title: string;
    message: string;
    confirmText?: string;
    onConfirm: () => void;
    variant: "danger" | "warning" | "info";
  }>({
    open: false,
    title: "",
    message: "",
    onConfirm: () => {},
    variant: "warning",
  });

  // Hooks
  const { products, isLoading, refreshProducts } = useProducts(filters);
  const { block } = useBlockedProducts();
  const { toast, hideToast } = useToast();
  const { data: systemStatus } = useSystemStatus();

  // Check if user has Bark Key on mount
  useEffect(() => {
    const barkKey = settingsService.getBarkKey();
    if (!barkKey) {
      // Small delay to let the app render first
      const timer = setTimeout(() => setIsWelcomeOpen(true), 500);
      return () => clearTimeout(timer);
    }
  }, []);

  // Handlers
  const handleFilterChange = (newFilters: Partial<TProductFilters>) => {
    setFilters((prev) => ({ ...prev, ...newFilters }));
  };

  const handleResetFilters = () => {
    setFilters({
      keyword: "",
      platform: "",
      region: "",
      salesStatus: "",
      monitorStatus: "",
      recentSevenDays: false,
    });
  };

  const handleSetNotification = (product: Product, isEdit: boolean) => {
    setNotificationModal({ open: true, product, isEdit });
  };

  const handleDeleteNotification = async (activityId: string) => {
    setConfirmModal({
      open: true,
      title: "取消监控",
      message: "确认要取消该产品的监控吗？",
      confirmText: "确认取消",
      variant: "warning",
      onConfirm: async () => {
        const { notificationService } =
          await import("@/services/notificationService");
        try {
          await notificationService.delete(activityId);
          refreshProducts();
          showToast("已取消监控");
        } catch {
          showToast("操作失败");
        }
      },
    });
  };

  const handleBlockProduct = async (activityId: string) => {
    setConfirmModal({
      open: true,
      title: "屏蔽产品",
      message: "确认屏蔽该产品吗？屏蔽后将不再显示。",
      confirmText: "确认屏蔽",
      variant: "danger",
      onConfirm: async () => {
        try {
          await block(activityId);
          showToast("已屏蔽产品");
        } catch {
          showToast("操作失败");
        }
      },
    });
  };

  const handleNotificationSuccess = (message: string) => {
    showToast(message);
    refreshProducts();
  };

  const handleSettingsSave = (message: string) => {
    showToast(message);
  };

  const handleWelcomeComplete = () => {
    refreshProducts();
    showToast("设置完成，开始监控好价吧！");
  };

  return (
    <div className="min-h-screen bg-slate-50 pb-10">
      <Navbar
        productCount={products.length}
        syncStatus={systemStatus?.sync}
        onOpenSettings={() => setIsSettingsOpen(true)}
        onOpenBlocked={() => setIsBlockedOpen(true)}
      />

      <SearchBar filters={filters} onFilterChange={handleFilterChange} />

      <ProductFilters
        filters={filters}
        onFilterChange={handleFilterChange}
        onReset={handleResetFilters}
      />

      <main className="p-3 max-w-7xl mx-auto">
        <ProductGrid
          products={products}
          isLoading={isLoading}
          onSetNotification={handleSetNotification}
          onDeleteNotification={handleDeleteNotification}
          onBlockProduct={handleBlockProduct}
        />
      </main>

      {/* Modals */}
      <NotificationModal
        product={notificationModal.product}
        isEdit={notificationModal.isEdit}
        open={notificationModal.open}
        onClose={() =>
          setNotificationModal({ open: false, product: null, isEdit: false })
        }
        onSuccess={handleNotificationSuccess}
        onError={showToast}
      />

      <BlockedProductsModal
        open={isBlockedOpen}
        onClose={() => setIsBlockedOpen(false)}
        onSuccess={showToast}
        onError={showToast}
      />

      <SettingsModal
        open={isSettingsOpen}
        onClose={() => setIsSettingsOpen(false)}
        onSave={handleSettingsSave}
      />

      <WelcomeModal
        open={isWelcomeOpen}
        onClose={() => setIsWelcomeOpen(false)}
        onComplete={handleWelcomeComplete}
      />

      <ConfirmModal
        open={confirmModal.open}
        title={confirmModal.title}
        message={confirmModal.message}
        confirmText={confirmModal.confirmText}
        variant={confirmModal.variant}
        onConfirm={confirmModal.onConfirm}
        onCancel={() => setConfirmModal((prev) => ({ ...prev, open: false }))}
      />

      {/* Toast */}
      <Toast
        message={toast.message}
        visible={toast.visible}
        onClose={hideToast}
      />
    </div>
  );
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AppContent />
    </QueryClientProvider>
  );
}
