import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import '../globals.css';
import { ThemeProvider } from '@/components/ThemeProvider';
import Navbar from '@/components/Navbar';
import Sidebar from '@/components/Sidebar';
import { CookiesProvider } from 'next-client-cookies/server';

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Dashboard',
  description: 'Dashboard for honeypot system',
}

export default function Layout({
  children,
}: {
  children: React.ReactNode,
}) {
  return (
    <main>
      <ThemeProvider>
        <CookiesProvider>
            <Sidebar />
            <div className="flex min-h-screen w-full flex-col bg-muted">
                <div className="sticky top-0 z-30 flex flex-col bg-background sm:pl-14">
                  <Navbar />
                </div>
                <div className="relative sm:py-4 sm:pl-14">
                  {children}
                </div>
            </div>
        </CookiesProvider>
      </ThemeProvider>
    </main>
  )
}
