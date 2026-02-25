'use client';


interface ErrorBannerProps {
  message: string;
  onDismiss?: () => void;
  onRetry?: () => void;
}

export function ErrorBanner({ message, onDismiss, onRetry }: ErrorBannerProps) {
  return (
    <div className="error-banner flex items-start justify-between gap-4">
      <div>
        <p className="font-medium">Something went wrong</p>
        <p className="text-sm opacity-90 mt-1">{message}</p>
      </div>
      <div className="flex gap-2">
        {onRetry && (
          <button
            onClick={onRetry}
            className="text-sm font-medium underline hover:opacity-75 transition-opacity whitespace-nowrap"
            aria-label="Retry"
          >
            Retry
          </button>
        )}
        {onDismiss && (
          <button
            onClick={onDismiss}
            className="text-sm font-medium underline hover:opacity-75 transition-opacity whitespace-nowrap"
            aria-label="Dismiss"
          >
            Dismiss
          </button>
        )}
      </div>
    </div>
  );
}
