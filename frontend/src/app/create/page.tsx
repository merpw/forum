import { Metadata } from "next"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"

import CreatePostPage from "@/app/create/create-page"
import { getCategoriesLocal } from "@/api/posts/edge"
import checkSession from "@/api/auth/edge"

export const metadata: Metadata = {
  title: "Create a new post",
}

const Page = async () => {
  const categories = await getCategoriesLocal()
  const token = cookies().get("forum-token")?.value
  if (!token) {
    return redirect("/login")
  }
  try {
    await checkSession(token)
  } catch (error) {
    return redirect("/login")
  }

  return (
    <CreatePostPage
      categories={categories}
      isAIEnabled={process.env.OPENAI_API_KEY !== undefined}
    />
  )
}

export default Page
