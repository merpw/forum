import { Metadata } from "next"

import SignupPage from "@/app/signup/signup-page"

// TODO: add server-side redirect if logged in

export const metadata: Metadata = {
  title: "Sign Up",
}
const Page = () => {
  return <SignupPage />
}

export default Page
