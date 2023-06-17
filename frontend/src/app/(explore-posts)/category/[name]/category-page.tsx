import { NextPage } from "next"

import { PostList } from "@/components/posts/list"
import { Post } from "@/custom"

const CategoryPage: NextPage<{ posts: Post[] }> = ({ posts }) => {
  return (
    <>
      <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
    </>
  )
}

export default CategoryPage
