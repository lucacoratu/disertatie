import { authRoutes, protectedRoutes } from "@/app/constants";
import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(request: NextRequest) {
    const token = request.cookies.get("session")?.value;

    if(request.nextUrl.pathname === "/" && !token) {
      return NextResponse.redirect(new URL("/login", request.url));
    }

    if(request.nextUrl.pathname.includes("/dashboard") && !token) {
      return NextResponse.redirect(new URL("/login", request.url));
    }

    if (authRoutes.includes(request.nextUrl.pathname) && token) {
      return NextResponse.redirect(new URL("/dashboard", request.url));
    }
}