"use client"

import Link from "next/link"
import {
  Home,
  PanelLeft,
  Search,
  Server,
  VenetianMask,
  ScrollText
} from "lucide-react"


import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";

import { Button } from "@/components/ui/button";

import { Input } from "@/components/ui/input";

import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet"

import ThemeButton from "./ThemeButton"
import { constants } from "@/app/constants"
import NotificationButton from "./NotificationButton"
import ProfileButton from "./ProfileButton";

import { usePathname } from "next/navigation";


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
        },
        {
            href:"/dashboard/logs",
            label: "Logs",
            icon: <ScrollText className="w-5 h-5" />,
        }
    ];

    //Get the current path from the URL
    const pathname = usePathname();
    const pathParts = pathname.split("/").slice(1);

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
                    {pathParts.map((part, i) => {
                        //Get the parts until the current one
                        const refToCurrent: string = "/" + pathParts.slice(0, i + 1).join("/");
                        part = part.charAt(0).toUpperCase() + part.slice(1);
                        return (
                            <>
                                <BreadcrumbItem >
                                    <BreadcrumbLink asChild>
                                        <Link href={refToCurrent}>{part}</Link>
                                    </BreadcrumbLink>
                                </BreadcrumbItem>
                                { i !== pathParts.length - 1 ? <BreadcrumbSeparator/> : <></> }
                            </>
                        );
                    })}
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