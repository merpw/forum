import { Metadata } from "next"

import HomePage from "@/app/home-page"
import { getCategoriesLocal, getPostsLocal } from "@/api/posts/edge"

export const metadata: Metadata = {
  title: "Recent Posts - Forum",
}

export const revalidate = 60

const Page = async () => {
  const posts = await getPostsLocal()
  const categories = await getCategoriesLocal()

  return <HomePage posts={posts} categories={categories} />
}

export default Page
