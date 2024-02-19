"use client"

import Container from "@/components/ui/container";
import Link from "next/link";
import { useTheme } from "next-themes";
import { Button } from "@/components/ui/button";

import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetFooter,
    SheetHeader,
    SheetTitle,
    SheetTrigger,
  } from "@/components/ui/sheet"

import { Menu, Bell, Server, Home, VenetianMask } from "lucide-react";
import ProfileButton from "./ProfileButton";
import ThemeButton from "./ThemeButton";
import { DropdownMenuSeparator } from "./ui/dropdown-menu";
import { constants } from "@/app/constants";
import NotificationButton from "./NotificationButton";

const Header = () => {
    const {theme, setTheme} = useTheme();

    const routes = [
        {
            href: "/",
            label: "Home",
            icon: <Home className="w-4 h-4"/>,
        },
        {
            href: "/machines",
            label: "Machines",
            icon: <Server className="w-4 h-4"/>,
        },
        {
            href: "/agents",
            label: "Agents",
            icon: <VenetianMask className="w-4 h-4"/>,
        }
    ];

    return (
        <header className="flex justify-start py-3 px-4 border-b dark:border-darksurface-400 dark:bg-darksurface-100">
            <Container>
                <div className="flex h-8 justify-between items-center w-full">
                    <div className="flex justify-start items-center gap-3">
                        <Sheet>
                            <SheetTrigger>
                                <Menu className="h-6 w-6"/>
                            </SheetTrigger>
                            <SheetContent side="left" className="w-[250px] flex flex-col justify-between dark:bg-darksurface-100">
                                <SheetHeader>
                                    <SheetTitle className="text-center">Navigation Menu</SheetTitle>
                                    {/* <SheetDescription className="text-center">
                                        Select the page you want to go to
                                    </SheetDescription> */}
                                    <DropdownMenuSeparator/>
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
                                </SheetHeader>
                                <SheetFooter className="self-center">
                                    <DropdownMenuSeparator/>
                                    <p className="text-center">v{constants.appVersion}</p>
                                </SheetFooter>
                            </SheetContent>
                        </Sheet>
                        <Link href="/" className="ml-4 lg:ml-0">
                            <h1 className="text-xl font-bold">
                               {constants.appName}
                            </h1>
                        </Link>
                    </div>
                    
                    <div className="flex flex-row items-center gap-6">
                        <NotificationButton />
                        <ThemeButton />   
                        <ProfileButton/>
                    </div>
                </div>
            </Container>
        </header>
    )
}

export default Header;