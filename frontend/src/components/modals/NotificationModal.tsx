import { useState, useEffect, useMemo, useRef } from 'react';
import { Bell, TrendingDown, AlertCircle, Check, Loader2, X, Sparkles } from 'lucide-react';
import type { Product } from '@/types';
import { useNotifications } from '@/hooks/useNotifications';

interface NotificationModalProps {
  product: Product | null;
  isEdit: boolean;
  open: boolean;
  onClose: () => void;
  onSuccess: (message: string) => void;
  onError: (message: string) => void;
}

// Price suggestion presets
const DISCOUNT_PRESETS = [
  { label: '9折', discount: 0.1, color: 'bg-blue-50 border-blue-200 text-blue-700' },
  { label: '8折', discount: 0.2, color: 'bg-green-50 border-green-200 text-green-700' },
  { label: '7折', discount: 0.3, color: 'bg-amber-50 border-amber-200 text-amber-700' },
  { label: '半价', discount: 0.5, color: 'bg-red-50 border-red-200 text-red-700' },
];

export function NotificationModal({
  product,
  isEdit,
  open,
  onClose,
  onSuccess,
  onError,
}: NotificationModalProps) {
  const [targetPrice, setTargetPrice] = useState('');
  const [inputError, setInputError] = useState<string | null>(null);
  const [sliderValue, setSliderValue] = useState(100);
  const inputRef = useRef<HTMLInputElement>(null);
  const { create, update, isCreating, isUpdating } = useNotifications();

  const currentPrice = product?.currentPrice || 0;

  // Calculate discount percentage
  const discountPercent = useMemo(() => {
    const target = parseFloat(targetPrice);
    if (!target || target >= currentPrice || currentPrice <= 0) return 0;
    return ((currentPrice - target) / currentPrice) * 100;
  }, [targetPrice, currentPrice]);

  // Calculate savings amount
  const savingsAmount = useMemo(() => {
    const target = parseFloat(targetPrice);
    if (!target || target >= currentPrice || currentPrice <= 0) return 0;
    return currentPrice - target;
  }, [targetPrice, currentPrice]);

  // Validate input
  const validatePrice = (value: string) => {
    const num = parseFloat(value);
    if (!value) {
      setInputError(null);
      return false;
    }
    if (isNaN(num) || num <= 0) {
      setInputError('请输入有效价格');
      return false;
    }
    if (num >= currentPrice) {
      setInputError('目标价格需低于当前价格');
      return false;
    }
    if (num < currentPrice * 0.1) {
      setInputError('目标价格过低');
      return false;
    }
    setInputError(null);
    return true;
  };

  // Reset form when modal opens with product data
  useEffect(() => {
    if (open && product) {
      const initialPrice = isEdit && product.targetPrice ? String(product.targetPrice) : '';
      setTargetPrice(initialPrice);
      setInputError(null);
      if (isEdit && product.targetPrice) {
        setSliderValue(Math.round((product.targetPrice / currentPrice) * 100));
      } else {
        setSliderValue(100);
      }
      // Focus input after animation
      setTimeout(() => inputRef.current?.focus(), 100);
    }
  }, [open, product, isEdit, currentPrice]);

  const handlePriceChange = (value: string) => {
    setTargetPrice(value);
    const num = parseFloat(value);
    if (validatePrice(value) && num > 0 && num < currentPrice) {
      setSliderValue(Math.round((num / currentPrice) * 100));
    }
  };

  const handleSliderChange = (value: number) => {
    setSliderValue(value);
    const newPrice = (currentPrice * value / 100).toFixed(2);
    setTargetPrice(newPrice);
    validatePrice(newPrice);
  };

  const handlePresetClick = (discount: number) => {
    const newPrice = (currentPrice * (1 - discount)).toFixed(2);
    setTargetPrice(newPrice);
    setSliderValue(Math.round((1 - discount) * 100));
    validatePrice(newPrice);
  };

  const handleSubmit = async () => {
    if (!product || !targetPrice || parseFloat(targetPrice) <= 0) {
      onError('请输入有效的目标价格');
      return;
    }

    if (inputError) {
      onError(inputError);
      return;
    }

    try {
      if (isEdit) {
        await update({
          activityId: product.activityId,
          params: { targetPrice: parseFloat(targetPrice) },
        });
        onSuccess('监控设置已更新');
      } else {
        await create({
          activityId: product.activityId,
          targetPrice: parseFloat(targetPrice),
        });
        onSuccess('监控已设置');
      }
      onClose();
    } catch {
      onError('操作失败，请稍后重试');
    }
  };

  const isValid = targetPrice && !inputError && parseFloat(targetPrice) > 0 && parseFloat(targetPrice) < currentPrice;
  const isLoading = isCreating || isUpdating;

  if (!open) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4"
      onClick={onClose}
    >
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" />

      {/* Modal */}
      <div
        className="relative bg-white rounded-2xl w-full max-w-sm shadow-2xl animate-scale-in overflow-hidden"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="relative bg-gradient-to-br from-primary-600 via-primary-500 to-teal-500 px-5 py-5">
          {/* Close button */}
          <button
            onClick={onClose}
            className="absolute top-3 right-3 w-8 h-8 flex items-center justify-center rounded-full bg-white/20 hover:bg-white/30 transition-colors"
          >
            <X className="w-4 h-4 text-white" />
          </button>

          <div className="flex items-center gap-3 pr-8">
            <div className="w-11 h-11 bg-white/20 backdrop-blur-sm rounded-xl flex items-center justify-center shadow-lg">
              <Bell className="w-5 h-5 text-white" />
            </div>
            <div>
              <h2 className="text-lg font-bold text-white">
                {isEdit ? '修改提醒' : '价格提醒'}
              </h2>
              <p className="text-xs text-white/70 mt-0.5">
                降价时自动通知您
              </p>
            </div>
          </div>
        </div>

        {/* Content */}
        <div className="p-5 space-y-4">
          {/* Product Info */}
          <div className="flex items-center gap-3 p-3 bg-gradient-to-r from-slate-50 to-slate-100/50 rounded-xl border border-slate-100">
            <div className="flex-1 min-w-0">
              <p className="text-xs text-slate-400 mb-0.5">商品</p>
              <p className="text-sm font-medium text-slate-800 line-clamp-1">
                {product?.title || '未知商品'}
              </p>
            </div>
            <div className="text-right flex-shrink-0 pl-3 border-l border-slate-200">
              <p className="text-xs text-slate-400 mb-0.5">现价</p>
              <p className="text-xl font-bold text-primary-600 tabular-nums">
                ¥{currentPrice.toFixed(2)}
              </p>
            </div>
          </div>

          {/* Target Price Input */}
          <div className="space-y-2">
            <label className="flex items-center gap-1.5 text-sm font-semibold text-slate-700">
              <TrendingDown className="w-4 h-4 text-primary-500" />
              目标价格
            </label>

            <div className="relative">
              <span className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 font-bold text-lg">
                ¥
              </span>
              <input
                ref={inputRef}
                type="number"
                step="0.01"
                placeholder="输入目标价格"
                value={targetPrice}
                onChange={(e) => handlePriceChange(e.target.value)}
                className={`w-full pl-10 pr-4 py-3.5 border-2 rounded-xl text-xl font-bold tabular-nums
                           transition-all duration-200 outline-none placeholder:text-slate-300
                           ${inputError
                             ? 'border-red-300 bg-red-50/50'
                             : isValid
                               ? 'border-primary-400 bg-primary-50/30'
                               : 'border-slate-200 bg-white hover:border-slate-300'
                           }`}
              />
              {isValid && (
                <span className="absolute right-4 top-1/2 -translate-y-1/2">
                  <div className="w-6 h-6 bg-primary-500 rounded-full flex items-center justify-center">
                    <Check className="w-3.5 h-3.5 text-white" />
                  </div>
                </span>
              )}
            </div>

            {/* Error / Success feedback */}
            {inputError && (
              <div className="flex items-center gap-1.5 text-xs text-red-500 px-1">
                <AlertCircle className="w-3.5 h-3.5" />
                {inputError}
              </div>
            )}

            {isValid && (
              <div className="flex items-center justify-between bg-gradient-to-r from-green-50 to-emerald-50 rounded-lg px-3.5 py-2.5 border border-green-100">
                <div className="flex items-center gap-1.5">
                  <Sparkles className="w-4 h-4 text-green-500" />
                  <span className="text-sm text-green-700">预计省下</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-base font-bold text-green-600 tabular-nums">
                    ¥{savingsAmount.toFixed(2)}
                  </span>
                  <span className="px-1.5 py-0.5 bg-green-500 text-white text-xs font-bold rounded">
                    -{discountPercent.toFixed(0)}%
                  </span>
                </div>
              </div>
            )}
          </div>

          {/* Quick Preset Buttons */}
          <div className="space-y-2">
            <p className="text-xs text-slate-400 px-0.5">快捷折扣</p>
            <div className="grid grid-cols-4 gap-2">
              {DISCOUNT_PRESETS.map((preset) => {
                const presetPrice = (currentPrice * (1 - preset.discount)).toFixed(2);
                const isSelected = Math.abs(parseFloat(targetPrice) - parseFloat(presetPrice)) < 0.01;
                return (
                  <button
                    key={preset.label}
                    type="button"
                    onClick={() => handlePresetClick(preset.discount)}
                    className={`py-2 px-2 rounded-xl text-center transition-all duration-200 border-2
                              ${isSelected
                                ? 'border-primary-400 bg-primary-50 shadow-sm scale-[1.02]'
                                : 'border-slate-100 bg-white hover:border-slate-200 hover:bg-slate-50'
                              }`}
                  >
                    <span className={`block text-sm font-bold ${isSelected ? 'text-primary-600' : 'text-slate-700'}`}>
                      {preset.label}
                    </span>
                    <span className="block text-[10px] text-slate-400 mt-0.5 tabular-nums">
                      ¥{presetPrice}
                    </span>
                  </button>
                );
              })}
            </div>
          </div>

          {/* Price Slider */}
          {currentPrice > 0 && (
            <div className="space-y-2">
              <div className="flex justify-between items-center text-xs px-0.5">
                <span className="text-slate-400">调整目标</span>
                <span className="font-medium text-slate-600 tabular-nums">
                  {sliderValue}% = ¥{(currentPrice * sliderValue / 100).toFixed(2)}
                </span>
              </div>
              <div className="relative px-1">
                <input
                  type="range"
                  min="10"
                  max="95"
                  value={sliderValue}
                  onChange={(e) => handleSliderChange(parseInt(e.target.value))}
                  className="w-full h-2 bg-slate-200 rounded-full appearance-none cursor-pointer
                           [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:w-5
                           [&::-webkit-slider-thumb]:h-5 [&::-webkit-slider-thumb]:bg-primary-500
                           [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:shadow-lg
                           [&::-webkit-slider-thumb]:border-2 [&::-webkit-slider-thumb]:border-white
                           [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:transition-transform
                           [&::-webkit-slider-thumb]:hover:scale-110"
                />
                <div className="flex justify-between text-[10px] text-slate-300 mt-1 px-0.5">
                  <span>¥{(currentPrice * 0.1).toFixed(0)}</span>
                  <span>¥{currentPrice.toFixed(0)}</span>
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="px-5 py-4 bg-slate-50/80 border-t border-slate-100 flex gap-3">
          <button
            type="button"
            onClick={onClose}
            disabled={isLoading}
            className="flex-1 px-4 py-3 bg-white text-slate-600 rounded-xl font-medium
                       border-2 border-slate-200 hover:bg-slate-50 hover:border-slate-300
                       transition-all duration-200 disabled:opacity-50 text-sm"
          >
            取消
          </button>
          <button
            type="button"
            onClick={handleSubmit}
            disabled={isLoading || !isValid}
            className="flex-[2] px-4 py-3 bg-gradient-to-r from-primary-600 to-primary-500
                       text-white rounded-xl font-semibold shadow-lg shadow-primary-500/20
                       hover:from-primary-700 hover:to-primary-600 hover:shadow-xl
                       disabled:opacity-40 disabled:cursor-not-allowed disabled:shadow-none
                       active:scale-[0.98] transition-all duration-200
                       flex items-center justify-center gap-2 text-sm"
          >
            {isLoading ? (
              <>
                <Loader2 className="w-4 h-4 animate-spin" />
                保存中
              </>
            ) : (
              <>
                <Bell className="w-4 h-4" />
                {isEdit ? '更新提醒' : '开启提醒'}
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
}
