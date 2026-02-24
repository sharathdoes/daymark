'use client';

import { useEffect, useState } from 'react';

const LOADING_MESSAGES = [
  'Generating your quiz...',
  'Curating the best questions...',
  'Preparing your challenge...',
  'Almost ready...',
  'One more moment...',
];

interface LoadingOverlayProps {
  isVisible: boolean;
  message?: string;
}

export function LoadingOverlay({ isVisible, message }: LoadingOverlayProps) {
  const [displayMessage, setDisplayMessage] = useState(message || LOADING_MESSAGES[0]);
  const [messageIndex, setMessageIndex] = useState(0);

  useEffect(() => {
    if (!isVisible) return;

    const interval = setInterval(() => {
      setMessageIndex((prev) => (prev + 1) % LOADING_MESSAGES.length);
      setDisplayMessage(LOADING_MESSAGES[messageIndex]);
    }, 5000);

    return () => clearInterval(interval);
  }, [isVisible, messageIndex]);

  useEffect(() => {
    if (message) {
      setDisplayMessage(message);
    }
  }, [message]);

  if (!isVisible) return null;

  return (
    <div className="loading-overlay" role="status" aria-live="polite" aria-label="Loading">
      <div className="flex flex-col items-center gap-4">
        <div className="flex gap-1">
          <div className="w-2 h-8 bg-accent rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
          <div className="w-2 h-8 bg-accent rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
          <div className="w-2 h-8 bg-accent rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
        </div>
        <p className="text-accent-foreground text-center text-lg font-serif">{displayMessage}</p>
      </div>
    </div>
  );
}
