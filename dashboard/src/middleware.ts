import { authRoutes, protectedRoutes } from "@/app/constants";
import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { cookies } from 'next/headers';
import { constants } from "@/app/constants";

async function checkToken(token: string): Promise<boolean> {
  //Check the token using the API
  const url = constants.apiBaseURL + "/auth/check-token";
  const res = await fetch(url, { cache: 'no-store' });
  return res.ok;
}

export async function middleware(request: NextRequest) {
    const token = request.cookies.get("session")?.value;

    if(request.nextUrl.pathname === "/" && !token) {
      return NextResponse.redirect(new URL("/login", request.url));
    }

    if(request.nextUrl.pathname.includes("/dashboard") && !token) {
      return NextResponse.redirect(new URL("/login", request.url));
    }

    if (authRoutes.includes(request.nextUrl.pathname) && token) {
      //Check if the token is still valid
      const isTokenValid = await checkToken(token);
      //If the token is not valid redirect the user back to login
      if(!isTokenValid) {
        //Clear the token from the cookies
        cookies().delete('token');
        return NextResponse.redirect(new URL("/login", request.url));
      }
      //If the token is valid redirect the user to the dashboard
      return NextResponse.redirect(new URL("/dashboard", request.url));
    }
}