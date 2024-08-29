"use client";

import { usePathname } from "next/navigation";
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";

import Link from "next/link";
import React from "react";

export default function BreadcrumbNavbar() {
    //Get the current path from the URL
    const pathname = usePathname();
    const pathParts = pathname.split("/").slice(1);

    return (
        <Breadcrumb className="hidden md:flex">
            <BreadcrumbList>
                {pathParts.map((part, i) => {
                    //Get the parts until the current one
                    const refToCurrent: string = "/" + pathParts.slice(0, i + 1).join("/");
                    part = part.charAt(0).toUpperCase() + part.slice(1);
                    return (
                        <React.Fragment key={i}>
                            <BreadcrumbItem>
                                <BreadcrumbLink asChild>
                                    <Link href={refToCurrent}>{part}</Link>
                                </BreadcrumbLink>
                            </BreadcrumbItem>
                            { i !== pathParts.length - 1 ? <BreadcrumbSeparator /> : <></> }
                        </React.Fragment>
                    );
                })}
            </BreadcrumbList>
        </Breadcrumb>
    );
}