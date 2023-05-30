import { Metadata } from "next"

import CreatePostPage from "@/app/create/create-page"
import { getCategoriesLocal } from "@/api/posts/edge"

export const metadata: Metadata = {
  title: "Create a new post",
}

const Page = async () => {
  const categories = await getCategoriesLocal()

  return (
    <CreatePostPage
      categories={categories}
      isAIEnabled={process.env.OPENAI_API_KEY !== undefined}
    />
  )
}

export default Page
