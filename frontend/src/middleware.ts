import { NextRequest, NextResponse } from "next/server"

import checkSession from "@/api/auth/edge"

export async function middleware(request: NextRequest) {
  if (process.env.FORUM_IS_PRIVATE === "false") {
    return NextResponse.next()
  }

  const token = request.cookies.get("forum-token")?.value

  if (!token) {
    // TODO: add query param to redirect back after login
    return NextResponse.redirect(new URL("/login", request.url))
  }

  try {
    await checkSession(token)
  } catch (error) {
    console.error("check session error", error)
    const resp = NextResponse.redirect(new URL("/login", request.url))
    resp.cookies.set("forum-token", "", { maxAge: 0 })
    return resp
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - login (login page)
     * - signup (signup page)
     */
    "/((?!api|_next/static|_next/image|favicon.ico|login|signup|avatar).*)",
  ],
}
