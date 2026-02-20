import { useState, useEffect, useCallback } from 'react';

interface ToastState {
  message: string;
  visible: boolean;
}

let showToastFn: (message: string) => void;

export function useToast() {
  const [toast, setToast] = useState<ToastState>({ message: '', visible: false });

  const showToast = useCallback((message: string) => {
    setToast({ message, visible: true });
  }, []);

  const hideToast = useCallback(() => {
    setToast((prev) => ({ ...prev, visible: false }));
  }, []);

  // Register the show function globally
  useEffect(() => {
    showToastFn = showToast;
  }, [showToast]);

  return { toast, hideToast };
}

// Global function to show toast from anywhere
export function showToast(message: string) {
  if (showToastFn) {
    showToastFn(message);
  }
}

interface ToastProps {
  message: string;
  visible: boolean;
  onClose: () => void;
}

export function Toast({ message, visible, onClose }: ToastProps) {
  useEffect(() => {
    if (visible) {
      const timer = setTimeout(onClose, 2000);
      return () => clearTimeout(timer);
    }
  }, [visible, onClose]);

  if (!visible) return null;

  return (
    <div
      className="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-[9999]
                 bg-slate-900/90 text-white px-6 py-3 rounded-lg text-sm
                 animate-scale-in"
    >
      {message}
    </div>
  );
}
