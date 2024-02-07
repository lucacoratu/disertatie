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

import { Menu, Bell, Home, VenetianMask } from "lucide-react";
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
            icon: <Home/>,
        },
        {
            href: "/machines",
            label: "Machines",
            icon: <VenetianMask/>,
        },
        {
            href: "/agents",
            label: "Agents",
            icon: <VenetianMask/>,
        }
    ];

    return (
        <header className="sm:flex sm:justify-start py-3 border-b dark:border-darksurface-400 dark:bg-darksurface-100">
            <Container>
                <div className="relative px-4 sm:px-6 lg:px-8 flex h-8 justify-between items-center w-full">
                    <div className="flex justify-start items-center gap-4">
                        <Sheet>
                            <SheetTrigger>
                                <Menu className="h-6 w-6"/>
                            </SheetTrigger>
                            <SheetContent side="left" className="w-[250px] flex flex-col justify-between items-center dark:bg-darksurface-100">
                                <SheetHeader>
                                    <SheetTitle className="text-center">Navigation Menu</SheetTitle>
                                    {/* <SheetDescription className="text-center">
                                        Select the page you want to go to
                                    </SheetDescription> */}
                                    <DropdownMenuSeparator/>
                                    <div className="flex flex-col">
                                        {routes.map((route, i) => {
                                            //console.log(i);
                                            return (
                                                <Button key={i} asChild variant="ghost" className="flex flex-col dark:hover:bg-primary light:hover:bg-secondary">
                                                    <Link key={i}
                                                        href={route.href}
                                                        className="text-sm font-medium transition-colors">
                                                            <div className="flex flex-row items-start gap-4">
                                                                {/* {route.icon} */}
                                                                {route.label}
                                                            </div>
                                                    </Link>
                                                </Button>
                                            )
                                        })}
                                    </div>
                                </SheetHeader>
                                <SheetFooter>
                                    <DropdownMenuSeparator/>
                                    <p className="text-center justify-center">v{constants.appVersion}</p>
                                </SheetFooter>
                            </SheetContent>
                        </Sheet>
                        <Link href="/" className="ml-4 lg:ml-0">
                            <h1 className="text-xl font-bold">
                               {constants.appName}
                            </h1>
                        </Link>
                    </div>
                    
                    <div className="mx-6 flex flex-row items-center space-x-4 lg:space-x-6">
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