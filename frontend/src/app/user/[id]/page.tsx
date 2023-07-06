import { Metadata } from "next"
import { notFound, redirect } from "next/navigation"
import { cookies } from "next/headers"

import UserPage from "@/app/user/[id]/user-page"
import { getUserLocal } from "@/api/users/edge"
import { getUserPostsLocal } from "@/api/posts/edge"
import checkSession from "@/api/auth/edge"

type Props = { params: { id: string } }

export const revalidate = 60

export const generateMetadata = async ({ params }: Props): Promise<Metadata> => {
  const id = +params.id
  if (isNaN(id)) {
    return notFound()
  }
  const user = await getUserLocal(id).catch(notFound)

  return {
    title: user.username,
    description: `Posts created by ${user.username}`,
  }
}

const Page = async ({ params }: Props) => {
  const token = cookies().get("forum-token")?.value

  if (token) {
    const requestUserId = await checkSession(token)
    if (requestUserId === +params.id) {
      return redirect("/me")
    }
  }

  const user = await getUserLocal(+params.id).catch(notFound)
  const posts = await getUserPostsLocal(user.id)

  return <UserPage user={user} posts={posts} />
}

export default Page
