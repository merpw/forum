import Link from "next/link"
import { FC } from "react"

import { PostList } from "@/components/posts/list"
import { Post } from "@/custom"

const Home: FC<{ posts: Post[]; categories: string[] }> = ({ posts, categories }) => {
  return (
    <>
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
export default Home
