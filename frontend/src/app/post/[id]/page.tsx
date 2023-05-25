import { notFound } from "next/navigation"
import { Metadata } from "next"

import { getPostCommentsLocal, getPostLocal, getPostsLocal } from "@/api/posts/edge"
import PostPage from "@/app/post/[id]/post-page"
import RenderMarkdown from "@/components/markdown/render"

type Props = { params: { id: string } }

export const generateMetadata = async ({ params }: Props): Promise<Metadata> => {
  const id = +params.id
  if (isNaN(id)) {
    return notFound()
  }
  const post = await getPostLocal(id).catch(notFound)

  return {
    title: post.title,
    description: post.description,
  }
}

export const generateStaticParams = async () => {
  if (!process.env.FORUM_BACKEND_PRIVATE_URL) {
    return []
  }
  const posts = await getPostsLocal()

  return posts.map((post) => ({
    id: post.id.toString(),
  }))
}

const Page = async ({ params }: { params: { id: string } }) => {
  const post = await getPostLocal(+params.id).catch(notFound)
  post.content = await RenderMarkdown(post.content)

  const comments = await getPostCommentsLocal(+params.id).catch(() => [])

  return <PostPage post={post} comments={comments} />
}

export default Page
