import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import { Analytics } from "@vercel/analytics/next";
import { ThemeProvider } from "next-themes";
import Header from "@/components/header";
import RootLayoutClient from "@/components/root-layout-client";
import "../styles/global.css";

const _geist = Geist({ subsets: ["latin"] });
const _geistMono = Geist_Mono({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Daymark - Test Your Knowledge With A Quiz",
  description: "Take interactive quizzes and track your knowledge with Daymark",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className="font-sans antialiased">
        <RootLayoutClient>
          <ThemeProvider attribute="class" defaultTheme="dark" enableSystem>
            <Header />
            {children}
            <Analytics />
          </ThemeProvider>
        </RootLayoutClient>
      </body>
    </html>
  );
}
