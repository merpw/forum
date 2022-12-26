import { GetStaticProps, NextPage } from "next"

import Head from "next/head"
import { Post } from "../custom"
import { PostList } from "../components/posts"
import { getPostsLocal } from "../api/posts/fetch"

const Home: NextPage<{ posts: Post[] }> = ({ posts }) => {
  return (
    <>
      <Head>
        <title>Recent Posts - Forum</title>
        <meta name={"description"} content={"The friendliest forum"} />
        <meta property={"og:title"} content={"FORUM"} key={"title"} />
        <meta name={"og:description"} content={"The friendliest forum"} />
      </Head>
      <PostList posts={posts} />
    </>
  )
}

export const getStaticProps: GetStaticProps<{ posts: Post[] }> = async () => {
  const posts: Post[] = await getPostsLocal()

  return {
    props: { posts: posts },
  }
}

export default Home
