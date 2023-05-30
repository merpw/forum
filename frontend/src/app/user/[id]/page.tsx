import { Metadata } from "next"
import { notFound } from "next/navigation"

import UserPage from "@/app/user/[id]/user-page"
import { getUserLocal } from "@/api/users/edge"
import { getUserPostsLocal } from "@/api/posts/edge"

type Props = { params: { id: string } }

export const revalidate = 60

export const generateMetadata = async ({ params }: Props): Promise<Metadata> => {
  const id = +params.id
  if (isNaN(id)) {
    return notFound()
  }
  const user = await getUserLocal(id).catch(notFound)

  return {
    title: user.name,
    description: `Posts created by ${user.name}`,
  }
}

const Page = async ({ params }: Props) => {
  const user = await getUserLocal(+params.id).catch(notFound)
  const posts = await getUserPostsLocal(user.id)

  return <UserPage user={user} posts={posts} />
}

export default Page
