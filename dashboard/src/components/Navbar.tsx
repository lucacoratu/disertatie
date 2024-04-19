"use client"

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
  VenetianMask
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
import ThemeButton from "./ThemeButton"
import { constants } from "@/app/constants"
import NotificationButton from "./NotificationButton"
import ProfileButton from "./ProfileButton"

export default function Navbar() {
    const routes = [
        {
            href: "/dashboard",
            label: "Home",
            icon: <Home className="w-4 h-4"/>,
        },
        {
            href: "/dashboard/machines",
            label: "Machines",
            icon: <Server className="w-4 h-4"/>,
        },
        {
            href: "/dashboard/agents",
            label: "Agents",
            icon: <VenetianMask className="w-4 h-4"/>,
        }
    ];

    return (
        <header className="sticky top-0 z-30 flex h-14 sm:py-3 sm:border-b items-center gap-4 border-b bg-background px-4 sm:static sm:h-auto sm:border-0 sm:bg-background sm:px-6">
            <Sheet>
            <SheetTrigger asChild>
                <Button size="icon" variant="outline" className="sm:hidden">
                <PanelLeft className="h-5 w-5" />
                <span className="sr-only">Toggle Menu</span>
                </Button>
            </SheetTrigger>
            <SheetContent side="left" className="sm:max-w-xs">
                <nav className="grid gap-6 text-lg font-medium">
                    <div className="flex flex-col">
                        {routes.map((route, i) => {
                            return (
                                <Button key={i} asChild variant="ghost" className="justify-start dark:hover:bg-primary light:hover:bg-secondary">
                                    <Link key={i}
                                        href={route.href}
                                        className="transition-colors flex flex-row justify-start gap-6">
                                            {route.icon}
                                            <p className="text-base font-medium">{route.label}</p>
                                    </Link>
                                </Button>
                            )
                        })}
                    </div>
                </nav>
            </SheetContent>
            </Sheet>
            <Link href="/dashboard" className="ml-4 lg:ml-0">
                <h1 className="text-xl font-bold">
                    {constants.appName}
                </h1>
            </Link>
            <Breadcrumb className="hidden md:flex">
            <BreadcrumbList>
                <BreadcrumbItem>
                <BreadcrumbLink asChild>
                    <Link href="#">Dashboard</Link>
                </BreadcrumbLink>
                </BreadcrumbItem>
                <BreadcrumbSeparator />
                <BreadcrumbItem>
                <BreadcrumbLink asChild>
                    <Link href="#">Machines</Link>
                </BreadcrumbLink>
                </BreadcrumbItem>
            </BreadcrumbList>
            </Breadcrumb>
            <div className="relative ml-auto flex-1 md:grow-0">
            <Search className="absolute left-2.5 top-[12px] h-4 w-4 text-muted-foreground" />
            <Input
                type="search"
                placeholder="Search..."
                className="w-full rounded-lg bg-background pl-8 md:w-[200px] lg:w-[336px]"
            />
            </div>
            
            <NotificationButton />
            <ThemeButton />   
            <ProfileButton/>
        </header>
    );
}