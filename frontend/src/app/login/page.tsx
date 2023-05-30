import { Metadata } from "next"

import LoginPage from "@/app/login/login-page"

// TODO: add server-side redirect if logged in

export const metadata: Metadata = {
  title: "Login",
}
const Page = () => {
  return <LoginPage />
}

export default Page
