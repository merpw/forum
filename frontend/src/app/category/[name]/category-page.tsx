import { NextPage } from "next"

import { PostList } from "@/components/posts/list"
import { Post } from "@/custom"

const CategoryPage: NextPage<{ categoryName: string; posts: Post[] }> = ({
  categoryName,
  posts,
}) => {
  return (
    <>
      <h1 className={"text-3xl font-light mb-5"}>
        <span className={"font-normal"}>{categoryName}</span> category
      </h1>
      <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
    </>
  )
}

export default CategoryPage
