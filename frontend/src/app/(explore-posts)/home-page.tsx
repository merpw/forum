"use client"

import { FC } from "react"

import { PostList } from "@/components/posts/list"
import { usePosts } from "@/api/posts/hooks"

const Home: FC = () => {
  const { posts } = usePosts()

  if (!posts) {
    return null
  }

  return (
    <>
      <div className={"flex flex-col"}>
        <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
      </div>
    </>
  )
}
export default Home
