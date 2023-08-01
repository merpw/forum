import { Metadata } from "next"
import { redirect } from "next/navigation"
import { cookies } from "next/headers"

import SignupForm from "@/components/forms/signup/SignUpForm"

// TODO: add server-side redirect if logged in

export const metadata: Metadata = {
  title: "Sign Up",
}
const Page = () => {
  const token = cookies().get("forum-token")?.value
  if (token) {
    return redirect("/me")
  }
  return <SignupForm />
}

export default Page
