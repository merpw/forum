import { Metadata } from "next"

import UserPage from "@/app/me/user-page"

export const metadata: Metadata = {
  title: "Profile",
}

const Page = () => {
  return <UserPage />
}

export default Page
