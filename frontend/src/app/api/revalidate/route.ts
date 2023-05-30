import { NextResponse } from "next/server"
import { revalidatePath } from "next/cache"

export const POST = async (req: Request) => {
  // Secret token is optional, by default frontend /api/ is not public
  const { searchParams } = new URL(req.url)

  const secret = searchParams.get("secret")

  if (process.env.FRONTEND_REVALIDATE_TOKEN && secret !== process.env.FRONTEND_REVALIDATE_TOKEN) {
    return NextResponse.json({ message: "Invalid secret" }, { status: 401 })
  }

  const url = searchParams.get("url")
  if (!url) {
    return NextResponse.json({ message: "No url provided" }, { status: 400 })
  }

  revalidatePath(url)

  return NextResponse.json({ revalidated: true })
}
