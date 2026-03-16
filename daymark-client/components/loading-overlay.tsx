'use client'

import { useEffect, useState } from 'react'

const MESSAGES = [
  'Gathering news articles...',
  'Crafting questions...',
  'Shuffling options...',
  'Preparing quiz...',
  'Almost ready...',
]

export default function LoadingOverlay() {
  const [messageIndex, setMessageIndex] = useState(0)

  useEffect(() => {
    const interval = setInterval(() => {
      setMessageIndex(prev => (prev + 1) % MESSAGES.length)
    }, 2000)

    return () => clearInterval(interval)
  }, [])

  return (
    <div className="fixed inset-0 bg-black/95 flex flex-col items-center justify-center z-50">
      <div className="text-center">
        <h2 className="text-3xl font-bold text-white mb-8 tracking-tight">Daymark</h2>
        <div className="text-[#666666] text-sm min-h-6">
          {MESSAGES[messageIndex]}
          <span className="animate-pulse">|</span>
        </div>
      </div>
    </div>
  )
}
