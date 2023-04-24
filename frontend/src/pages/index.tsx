import { GetStaticProps, NextPage } from "next"
import Link from "next/link"
import { NextSeo } from "next-seo"

import { Post } from "@/custom"
import { PostList } from "@/components/posts/list"
import { getCategoriesLocal, getPostsLocal } from "@/api/posts/fetch"

const Home: NextPage<{ posts: Post[]; categories: string[] }> = ({ posts, categories }) => {
  return (
    <>
      <NextSeo title={"Recent Posts"} />

      <div className={"flex gap-2 flex-wrap justify-center mb-5"}>
        <span className={"text-2xl my-3 w-full text-center"}>Categories:</span>
        {categories.map((category, key) => (
          <span key={key} className={"text-3xl border rounded-lg px-5 py-2 capitalize"}>
            <Link href={`/category/${category}`} className={"clickable"}>
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
    return { notFound: true, revalidate: 60 }
  }

  const posts: Post[] = await getPostsLocal()
  const categories = await getCategoriesLocal()
  return {
    props: { posts, categories },
    revalidate: 60,
  }
}

export default Home
