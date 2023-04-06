import { GetStaticProps, NextPage } from "next"
import Link from "next/link"
import { NextSeo } from "next-seo"
import { AxiosError } from "axios"

import { Post } from "@/custom"
import { PostList } from "@/components/posts/list"
import { getCategoriesLocal, getPostsLocal } from "@/api/posts/fetch"

const Home: NextPage<{ posts: Post[]; categories: string[] }> = ({ posts, categories }) => {
  return (
    <>
      <NextSeo title={"Recent Posts"} />

      <div className={"flex gap-2 flex-wrap justify-center mb-5"}>
        <span className={"text-3xl rounded-lg px-5 py-2 w-full text-center"}>Categories:</span>
        {categories.map((category, key) => (
          <span key={key} className={"text-3xl border rounded-lg px-5 py-2 capitalize"}>
            <Link href={`/category/${category}`} className={"hover:opacity-50"}>
              {category}
            </Link>
          </span>
        ))}
      </div>
      <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
    </>
  )
}

export const getStaticProps: GetStaticProps<{ posts: Post[]; categories: string[] }> = async () => {
  if (!process.env.FORUM_BACKEND_PRIVATE_URL) {
    return { notFound: true }
  }

  try {
    const posts: Post[] = await getPostsLocal()
    const categories = await getCategoriesLocal()
    return {
      props: { posts, categories },
      revalidate: 1,
    }
  } catch (e) {
    if ((e as AxiosError).response?.status !== 404) {
      throw e
    }
    return { notFound: true, revalidate: 1 }
  }
}

export default Home
