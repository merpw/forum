import { Metadata } from "next"

import { getCategoriesLocal } from "@/api/posts/edge"
import CategoryPage from "@/app/(explore-posts)/category/[name]/category-page"
import { Capitalize } from "@/helpers/text"

type Props = { params: { name: string } }

export const revalidate = 60

export const generateMetadata = async ({ params }: Props): Promise<Metadata> => {
  return {
    title: `${Capitalize(params.name)} category`,
    description: `Posts in ${Capitalize(params.name)} category`,
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
  return <CategoryPage category={params.name} />
}

export default Page
