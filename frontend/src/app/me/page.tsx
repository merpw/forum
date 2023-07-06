import { Metadata } from "next"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"

import UserPage from "@/app/me/user-page"
import checkSession from "@/api/auth/edge"

export const metadata: Metadata = {
  title: "Profile",
}

const Page = async () => {
  const token = cookies().get("forum-token")?.value
  if (!token) {
    return redirect("/login")
  }
  try {
    await checkSession(token)
    return <UserPage />
  } catch (error) {
    return redirect("/login")
  }
}

export default Page
