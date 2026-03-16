'use client'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useAuth, useTheme } from '@/lib/store'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { Moon, Sun } from 'lucide-react'
import { LogOut } from 'lucide-react'

interface HeaderProps {
  hideAuth?: boolean
}

export default function Header({ hideAuth = false }: HeaderProps) {
  const router = useRouter()
  const { user, isAuthenticated, clearAuth } = useAuth()
  const { isDark, setIsDark } = useTheme()

  const handleLogout = () => {
    clearAuth()
    router.push('/')
  }

  return (
    <header className="border-b border-border/60 bg-background/80 backdrop-blur sticky top-0 z-40">
      <div className="max-w-5xl mx-auto flex items-center justify-between px-4 py-3">
        {/* Logo */}
        <Link
          href="/"
          className="text-base md:text-lg font-semibold tracking-tight text-foreground hover:text-primary transition-colors"
        >
          Daymark
        </Link>

        {/* Right side */}
        <div className="flex items-center gap-3">
          {/* Theme Toggle */}
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={() => setIsDark(!isDark)}
            aria-label="Toggle theme"
          >
            <span className="text-lg" aria-hidden="true">
              {isDark ? (
                <Sun className="h-5 w-5" />
              ) : (
                <Moon className="h-5 w-5" />
              )}
            </span>
          </Button>

          {/* Auth Actions */}
          {!hideAuth && (
            <div className="flex items-center gap-2">
              {isAuthenticated && user ? (
                <>
                  <div className="flex items-center gap-2">
                    <Link href="/profile" className="flex items-center gap-2 hover:opacity-80 transition-opacity">
                      <Avatar>
                        <AvatarImage
                          src={
                            user.avatar_url ||
                            `https://api.dicebear.com/7.x/identicon/svg?seed=${encodeURIComponent(
                              user.email || 'user',
                            )}`
                          }
                          alt={user.name || user.email || 'User avatar'}
                        />
                        <AvatarFallback>
                          {(user.name || user.email || '?').charAt(0).toUpperCase()}
                        </AvatarFallback>
                      </Avatar>
                      <span className="text-sm text-muted-foreground inline-flex">
                        {user.name || user.email}
                      </span>
                    </Link>
                  </div>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={handleLogout}
                  >
                    <LogOut className="h-5 w-5" />
                  </Button>
                </>
              ) : (
                <>
                  <Link href="/login">
                    <Button type="button" variant="ghost" size="sm">
                      Sign in
                    </Button>
                  </Link>
                  <Link href="/signup">
                    <Button type="button" size="sm">
                      Sign up
                    </Button>
                  </Link>
                </>
              )}
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
