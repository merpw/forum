import { Metadata } from "next"

import HomePage from "@/app/(explore-posts)/home-page"
import { getPostsLocal } from "@/api/posts/edge"

export const metadata: Metadata = {
  title: "Recent Posts - Forum",
}

export const revalidate = 60

const Page = async () => {
  const posts = await getPostsLocal()

  return <HomePage posts={posts} />
}

export default Page
