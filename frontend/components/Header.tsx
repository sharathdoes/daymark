import Link from 'next/link';

interface HeaderProps {
  showBackButton?: boolean;
  backHref?: string;
}

export function Header({ showBackButton, backHref = '/' }: HeaderProps) {
  return (
    <header className="border-b border-border bg-background sticky top-0 z-40">
      <div className="container flex items-center justify-between h-16 md:h-20">
        <Link
          href="/"
          className="text-2xl md:text-3xl font-serif font-bold text-foreground hover:opacity-75 transition-opacity"
        >
          Daily Brief
        </Link>
        {showBackButton && (
          <Link
            href={backHref}
            className="text-accent underline font-medium hover:opacity-75 transition-opacity"
          >
            Back
          </Link>
        )}
      </div>
    </header>
  );
}
