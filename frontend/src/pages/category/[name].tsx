import { GetStaticPaths, GetStaticProps, NextPage } from "next"
import { NextSeo } from "next-seo"
import { AxiosError } from "axios"

import { getCategoriesLocal, getCategoryPostsLocal } from "@/api/posts/fetch"
import { PostList } from "@/components/posts/list"
import { Post } from "@/custom"
import { Capitalize } from "@/helpers/text"

const CategoryPage: NextPage<{ category_name: string; posts: Post[] }> = ({
  category_name,
  posts,
}) => {
  return (
    <>
      <NextSeo title={category_name} description={`Posts in ${category_name} category`} />
      <h1 className={"text-3xl font-light mb-5"}>
        <span className={"font-normal"}>{category_name}</span> category
      </h1>
      <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
    </>
  )
}

export const getStaticPaths: GetStaticPaths<{ name: string }> = async () => {
  if (!process.env.FORUM_BACKEND_PRIVATE_URL) {
    return { paths: [], fallback: "blocking" }
  }
  const categories = await getCategoriesLocal()
  return {
    paths: categories.map((category) => {
      return { params: { name: category } }
    }),
    // TODO: maybe remove
    fallback: "blocking", // fallback tries to regenerate ArtistPage if Artist did not exist during building
  }
}

export const getStaticProps: GetStaticProps<{ posts: Post[] }, { name: string }> = async ({
  params,
}) => {
  if (!params?.name) {
    return { notFound: true, revalidate: 60 }
  }

  try {
    const category_name = params.name.toLowerCase()
    const posts = await getCategoryPostsLocal(category_name)

    return {
      props: {
        posts: posts,
        category_name: Capitalize(category_name),
      },
      revalidate: 60,
    }
  } catch (e) {
    if ((e as AxiosError).response?.status !== 404) {
      throw e
    }
    return { notFound: true, revalidate: 60 }
  }
}

export default CategoryPage
