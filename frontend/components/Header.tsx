"use client";

import Link from "next/link";
import { useTheme } from "next-themes";
import { Sun, Moon, ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";

interface HeaderProps {
  showBackButton?: boolean;
  backHref?: string;
}

export function Header({ showBackButton, backHref = "/" }: HeaderProps) {
  const { theme, setTheme } = useTheme();

  const toggleTheme = () => {
    setTheme(theme === "dark" ? "light" : "dark");
  };

  return (
    <header className="sticky top-0 z-40 bg-background border-b">
      <div className="max-w-3xl mx-auto px-6 h-14 flex items-center justify-between">
        {/* Left */}
        <div className="flex items-center gap-4">
          {showBackButton && (
            <Link href={backHref}>
              <Button variant="ghost" size="icon">
                <ArrowLeft className="h-4 w-4" />
              </Button>
            </Link>
          )}

          <Link
            href="/"
            className="text-base font-medium text-foreground hover:opacity-70 transition-opacity"
          >
            Daily Brief
          </Link>
        </div>

        {/* Right */}
        <Button
          variant="ghost"
          size="icon"
          onClick={toggleTheme}
          aria-label="Toggle theme"
        >
          {theme === "dark" ? (
            <Sun className="h-4 w-4" />
          ) : (
            <Moon className="h-4 w-4" />
          )}
        </Button>
      </div>
    </header>
  );
}
