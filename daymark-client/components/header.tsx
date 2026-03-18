'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { useAuth, useTheme } from '@/lib/store'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { Moon, Sun, LogOut, Star, Github } from 'lucide-react'

interface HeaderProps {
  hideAuth?: boolean
}

export default function Header({ hideAuth = false }: HeaderProps) {
  const router = useRouter()
  const pathname = usePathname()
  const { user, isAuthenticated, clearAuth } = useAuth()
  const { isDark, setIsDark } = useTheme()

  const handleLogout = () => {
    clearAuth()
    router.push('/')
  }


  return (
    <header className="sticky top-0 z-40 border-b border-border/40 bg-background/80 backdrop-blur-sm">
      <div className="max-w-5xl mx-auto flex items-center justify-between h-12 px-4">

        {/* Logo */}
        <Link
          href="/"
          className="text-sm font-semibold tracking-tight text-foreground hover:text-primary transition-colors"
        >
          Daymark
        </Link>


        {/* Right actions */}
        <div className="flex items-center gap-1">

          <Link
            href="https://github.com/sharathdoes/daymark"
            target="_blank"
            className="hidden md:flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm text-muted-foreground hover:text-foreground hover:bg-accent/50 transition-colors"
          >
            <Github className="h-3.5 w-3.5" />
            Star
          </Link>

          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={() => setIsDark(!isDark)}
            aria-label="Toggle theme"
          >
            {isDark ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
          </Button>

          {!hideAuth && (
            <>
              {isAuthenticated && user ? (
                <div className="flex items-center gap-1 ml-1">
                  <Link
                    href="/profile"
                    className="flex items-center gap-2 px-2 py-1 rounded-md hover:bg-accent/50 transition-colors"
                  >
                    <Avatar className="h-6 w-6">
                      <AvatarImage
                        src={
                          user.avatar_url ||
                          `https://api.dicebear.com/7.x/identicon/svg?seed=${encodeURIComponent(user.email || 'user')}`
                        }
                        alt={user.name || user.email || 'User avatar'}
                      />
                      <AvatarFallback className="text-xs">
                        {(user.name || user.email || '?').charAt(0).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                    <span className="hidden md:inline text-sm text-muted-foreground">
                      {user.name || user.email}
                    </span>
                  </Link>

                  <Button
                    type="button"
                    variant="ghost"
                    size="icon-sm"
                    onClick={handleLogout}
                    aria-label="Log out"
                  >
                    <LogOut className="h-4 w-4" />
                  </Button>
                </div>
              ) : (
                <div className="flex items-center gap-1 ml-1">
                  <Link href="/login">
                    <Button type="button" size="sm">Sign in</Button>
                  </Link>
                 
                </div>
              )}
            </>
          )}
        </div>

      </div>
    </header>
  )
}