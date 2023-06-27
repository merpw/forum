import { FC } from "react"

import { Post } from "@/custom"
import { PostList } from "@/components/posts/list"

const Home: FC<{ posts: Post[] }> = ({ posts }) => {
  return (
    <>
      <div className={"flex flex-col"}>
        <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
      </div>
    </>
  )
}
export default Home
