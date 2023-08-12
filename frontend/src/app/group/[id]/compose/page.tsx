import { Metadata } from "next"
import { cookies } from "next/headers"
import { notFound, redirect } from "next/navigation"

import { getCategoriesLocal } from "@/api/posts/edge"
import checkSession from "@/api/auth/edge"
import CreatePostForm from "@/components/forms/create-post/CreatePostForm"
import { getGroupLocal, getGroupMembersLocal } from "@/api/groups/edge"

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

  const categories = await getCategoriesLocal()
  const token = cookies().get("forum-token")?.value
  if (!token) {
    return redirect("/login")
  }

  const userId = await checkSession(token).catch(notFound)

  const groupMembers = await getGroupMembersLocal(group.id)
  if (!groupMembers.includes(userId)) {
    return redirect(`/group/${group.id}`)
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
