import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import '../globals.css'
import { ThemeProvider } from '@/components/ThemeProvider'
import Header from "@/components/Header"
import { Toaster } from "@/components/ui/sonner"
import Navbar from '@/components/Navbar'
import Sidebar from '@/components/Sidebar'

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
    <html lang="en">
      <body className={inter.className}>
          {/* <Header/> */}
          <Sidebar />
          <div className="flex min-h-screen w-full flex-col bg-muted">
              <div className="flex flex-col sm:gap-4 sm:pl-14">
                  <Navbar />
                  {children}
              </div>
          </div>
      </body>
    </html>
  )
}
