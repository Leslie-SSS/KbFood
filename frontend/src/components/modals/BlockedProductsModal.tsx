import { useState } from "react";
import { RotateCcw } from "lucide-react";
import { useBlockedProducts } from "@/hooks/useBlockedProducts";
import { ConfirmModal } from "./ConfirmModal";

interface BlockedProductsModalProps {
  open: boolean;
  onClose: () => void;
  onSuccess: (message: string) => void;
  onError: (message: string) => void;
}

export function BlockedProductsModal({
  open,
  onClose,
  onSuccess,
  onError,
}: BlockedProductsModalProps) {
  const { blockedProducts, isLoading, unblock, isUnblocking } =
    useBlockedProducts();
  const [confirmModal, setConfirmModal] = useState<{
    open: boolean;
    activityId: string;
    title: string;
  }>({ open: false, activityId: "", title: "" });

  const handleUnblock = async (activityId: string, title: string) => {
    setConfirmModal({ open: true, activityId, title });
  };

  const handleConfirmUnblock = async () => {
    try {
      await unblock(confirmModal.activityId);
      onSuccess("已恢复显示");
    } catch {
      onError("操作失败");
    }
  };

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      onClose();
    }
  };

  if (!open) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center"
      onClick={() => handleOpenChange(false)}
    >
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" />

      {/* Modal */}
      <div
        className="relative bg-white rounded-2xl w-[90%] max-w-[600px] max-h-[80vh] shadow-modal animate-scale-in
                   flex flex-col"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="px-5 py-4 border-b border-slate-200 shrink-0">
          <h2 className="text-lg font-semibold text-slate-900">
            已屏蔽产品列表
          </h2>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-4 space-y-3">
          {isLoading ? (
            <div className="text-center py-10 text-slate-500">加载中...</div>
          ) : blockedProducts.length === 0 ? (
            <div className="text-center py-10 text-slate-500">
              <div className="text-4xl mb-3">🚫</div>
              <p>暂无屏蔽产品</p>
            </div>
          ) : (
            blockedProducts.map((product) => (
              <div
                key={product.activityId}
                className="p-3 border border-slate-200 rounded-xl space-y-2"
              >
                <div className="font-medium text-sm text-slate-900 line-clamp-1">
                  {product.title}
                </div>
                <div className="flex gap-1.5 text-xs text-slate-500">
                  <span>{product.platform}</span>
                  <span className="text-slate-300">-</span>
                  <span>{product.region}</span>
                  <span className="text-slate-300">-</span>
                  <span className="text-danger font-medium">
                    ¥{product.currentPrice}
                  </span>
                </div>
                <button
                  onClick={() =>
                    handleUnblock(product.activityId, product.title)
                  }
                  disabled={isUnblocking}
                  className="w-full px-3 py-2 text-xs rounded-lg border border-primary-500 text-primary-600
                             hover:bg-primary-50 disabled:opacity-50 disabled:cursor-not-allowed
                             transition-colors flex items-center justify-center gap-1"
                >
                  <RotateCcw className="w-3 h-3" />
                  恢复显示
                </button>
              </div>
            ))
          )}
        </div>

        {/* Footer */}
        <div className="px-5 py-4 border-t border-slate-200 shrink-0">
          <button
            onClick={onClose}
            className="w-full px-4 py-2.5 bg-slate-100 text-slate-700 rounded-lg font-medium
                       hover:bg-slate-200 transition-colors"
          >
            关闭
          </button>
        </div>
      </div>

      {/* Confirm Modal */}
      <ConfirmModal
        open={confirmModal.open}
        title="恢复显示"
        message={`确认恢复「${confirmModal.title.substring(0, 20)}...」的显示吗？`}
        confirmText="确认恢复"
        variant="info"
        onConfirm={handleConfirmUnblock}
        onCancel={() => setConfirmModal((prev) => ({ ...prev, open: false }))}
      />
    </div>
  );
}
