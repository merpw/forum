import { Metadata } from "next"
import { cookies } from "next/headers"
import { notFound, redirect } from "next/navigation"

import { getCategoriesLocal } from "@/api/posts/edge"
import checkSession from "@/api/auth/edge"
import CreatePostForm from "@/components/forms/create-post/CreatePostForm"
import { getGroupLocal } from "@/api/groups/edge"

type Props = { params: { id: string } }

export const generateMetadata = async ({ params }: Props): Promise<Metadata> => {
  const id = +params.id
  if (isNaN(id)) {
    return notFound()
  }
  const group = await getGroupLocal(id).catch(notFound)

  return {
    title: `Create a post in group "${group.title}"`,
  }
}

const Page = async ({ params }: Props) => {
  if (!params.id || isNaN(+params.id)) {
    return redirect("/groups")
  }

  const group = await getGroupLocal(+params.id).catch(notFound)
  if (!group) {
    return redirect("/groups")
  }

  // TODO: add check if user is a member of the group

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
    <CreatePostForm
      categories={categories}
      isAIEnabled={Boolean(process.env.OPENAI_API_KEY)}
      group={group}
    />
  )
}

export default Page
