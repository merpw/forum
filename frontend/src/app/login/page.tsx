import { Metadata } from "next"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"

import LoginPage from "@/app/login/login-page"

// TODO: add server-side redirect if logged in

export const metadata: Metadata = {
  title: "Login",
}
const Page = async () => {
  const token = cookies().get("forum-token")?.value
  if (token) {
    return redirect("/me")
  }

  return <LoginPage />
}

export default Page
