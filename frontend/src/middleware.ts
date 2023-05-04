import { NextRequest, NextResponse } from "next/server"

export async function middleware(request: NextRequest) {
  const token = request.cookies.get("forum-token")?.value

  if (!token) {
    // TODO: add query param to redirect back after login
    return NextResponse.redirect(new URL("/login", request.url))
  }

  const response = await fetch(
    `${process.env.FORUM_BACKEND_PRIVATE_URL}/api/internal/check-session?token=${token}`
  )

  const body = await response.json()
  if (body.error) {
    const resp = NextResponse.rewrite(new URL("/login", request.url))
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
    "/((?!api|_next/static|_next/image|favicon.ico|login|signup).*)",
  ],
}