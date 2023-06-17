import { Metadata } from "next"
import { notFound } from "next/navigation"

import { getCategoriesLocal, getCategoryPostsLocal } from "@/api/posts/edge"
import CategoryPage from "@/app/(explore-posts)/category/[name]/category-page"
import { Capitalize } from "@/helpers/text"

type Props = { params: { name: string } }

export const revalidate = 60

export const generateMetadata = async ({ params }: Props): Promise<Metadata> => {
  const posts = await getCategoryPostsLocal(params.name).catch(notFound)

  return {
    title: `${Capitalize(params.name)} category`,
    description: `More than ${posts.length} posts in ${params.name} category`,
  }
}

export const generateStaticParams = async () => {
  if (!process.env.FORUM_BACKEND_PRIVATE_URL) {
    return []
  }
  const categories = await getCategoriesLocal()

  return categories.map((category) => ({
    name: category,
  }))
}

const Page = async ({ params }: Props) => {
  const posts = await getCategoryPostsLocal(params.name).catch(notFound)

  return <CategoryPage posts={posts} />
}

export default Page
