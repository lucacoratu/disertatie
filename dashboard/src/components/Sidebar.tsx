import Image from "next/image"
import Link from "next/link"
import {
  CreditCard,
  Home,
  LineChart,
  Ghost,
  PanelLeft,
  Search,
  Settings,
  Server,
  Users2,
  Users,
  DollarSign,
  Activity,
  ArrowUpRight,
  VenetianMask,
  ScrollText,
} from "lucide-react"

import { Badge } from "@/components/ui/badge"

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Input } from "@/components/ui/input"

import {
    Avatar,
    AvatarFallback,
    AvatarImage,
  } from "@/components/ui/avatar"

import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
  TooltipProvider,
} from "@/components/ui/tooltip"

export default function Sidebar() {
  const routes = [
      {
          href: "/dashboard",
          label: "Home",
          icon: <Home className="w-5 h-5"/>,
      },
      {
          href: "/dashboard/machines",
          label: "Machines",
          icon: <Server className="w-5 h-5"/>,
      },
      {
          href: "/dashboard/agents",
          label: "Agents",
          icon: <VenetianMask className="w-5 h-5"/>,
      },
      {
          href:"/dashboard/logs",
          label: "Logs",
          icon: <ScrollText className="w-5 h-5" />,
      }
  ];

    return (
        <aside className="fixed z-40 inset-y-0 left-0 z-10 hidden w-14 flex-col border-r bg-background sm:flex">
          <nav className="flex flex-col items-center gap-4 px-2 sm:py-5">
          <Link
              href="#"
              className="group flex h-9 w-9 shrink-0 items-center justify-center gap-2 rounded-full bg-primary text-lg font-semibold text-primary-foreground md:h-8 md:w-8 md:text-base"
            >
              <Ghost className="h-5 w-5 transition-all group-hover:scale-110" />
              <span className="sr-only">Dashboard</span>
            </Link>
          <div className="flex flex-col">
            {routes.map((route, i) => {
                return (
                  <TooltipProvider key={i}>
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button key={i} asChild variant="ghost" className="justify-start dark:hover:bg-primary light:hover:bg-secondary">
                          <Link key={i}
                              href={route.href}
                              className="transition-colors flex flex-row justify-start gap-6">
                                  {route.icon}
                          </Link>
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent side="right">{route.label}</TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                )
            })}
        </div>
          </nav>
          <nav className="mt-auto flex flex-col items-center gap-2 px-2 sm:py-5">
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Link href="#" className="flex h-9 w-9 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:text-foreground md:h-8 md:w-8">
                    <Settings className="h-5 w-5" />
                    <span className="sr-only">Settings</span>
                  </Link>
                </TooltipTrigger>
                <TooltipContent side="right">Settings</TooltipContent>
              </Tooltip>
            </TooltipProvider>
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <label className="text-sm text-muted-foreground">v1.0</label>
              </TooltipTrigger>
              <TooltipContent side="right">Version</TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </nav>
      </aside>
    )
}