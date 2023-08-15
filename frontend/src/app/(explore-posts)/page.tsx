import { Metadata } from "next"

import HomePage from "@/app/(explore-posts)/home-page"

export const metadata: Metadata = {
  title: "Recent Posts - Forum",
}

export const revalidate = 60

const Page = async () => {
  return <HomePage />
}

export default Page
