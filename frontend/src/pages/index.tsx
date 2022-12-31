import { GetStaticProps, NextPage } from "next"

import Head from "next/head"
import { Post } from "../custom"
import { PostList } from "../components/posts/list"
import { getCategoriesLocal, getPostsLocal } from "../api/posts/fetch"
import Link from "next/link"

const Home: NextPage<{ posts: Post[]; categories: string[] }> = ({ posts, categories }) => {
  return (
    <>
      <Head>
        <title>Recent Posts - Forum</title>
        <meta name={"description"} content={"The friendliest forum"} />
        <meta property={"og:title"} content={"FORUM"} key={"title"} />
        <meta name={"og:description"} content={"The friendliest forum"} />
      </Head>
      <div className={"flex gap-2 flex-wrap justify-center mb-5"}>
        <span className={"text-3xl rounded-lg px-5 py-2 w-full text-center"}>Categories:</span>
        {categories.map((category, key) => (
          <span key={key} className={"text-3xl border rounded-lg px-5 py-2"}>
            <Link href={`/category/${category}`} className={"hover:opacity-50"}>
              {category.charAt(0).toUpperCase() + category.slice(1)}
            </Link>
          </span>
        ))}
      </div>
      <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
    </>
  )
}

export const getStaticProps: GetStaticProps<{ posts: Post[]; categories: string[] }> = async () => {
  const posts: Post[] = await getPostsLocal()
  const categories = await getCategoriesLocal()

  return {
    props: { posts, categories },
    revalidate: 10,
  }
}

export default Home
