'use client'

import { useState, useEffect } from 'react'
import { AuthProvider } from '@/lib/store'

export default function RootLayoutClient({
  children,
}: {
  children: React.ReactNode
}) {
  const [isDark, setIsDark] = useState(true)
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
    const theme = localStorage.getItem('daymark_theme') || 'dark'
    setIsDark(theme === 'dark')
    updateTheme(theme === 'dark')
  }, [])

  useEffect(() => {
    if (!mounted) return
    const html = document.documentElement
    if (isDark) {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
    }
    localStorage.setItem('daymark_theme', isDark ? 'dark' : 'light')
    updateTheme(isDark)
  }, [isDark, mounted])

  const updateTheme = (dark: boolean) => {
    const html = document.documentElement
    if (dark) {
      html.classList.add('dark')
      html.style.colorScheme = 'dark'
    } else {
      html.classList.remove('dark')
      html.style.colorScheme = 'light'
    }
  }

  if (!mounted) {
    return null
  }

  return (
    <AuthProvider themeToggle={{ isDark, setIsDark }}>
      {children}
    </AuthProvider>
  )
}
