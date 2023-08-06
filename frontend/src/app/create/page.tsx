import { Metadata } from "next"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"

import CreatePostForm from "@/components/forms/create-post/CreatePostForm"
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
    <CreatePostForm categories={categories} isAIEnabled={Boolean(process.env.OPENAI_API_KEY)} />
  )
}

export default Page
