import Link from "next/link"
import { FC } from "react"

import { Category, CommentsCount, ReactionsButtons } from "./reactions"

import { Post } from "@/custom"
import AutoDate from "@/components/AutoDate"

export const PostList: FC<{ posts: Post[] }> = ({ posts }) => {
  if (posts.length == 0) {
    return <div className={"text-info text-center mt-5 mb-7"}>There are no posts yet...</div>
  }
  return (
    <div>
      {posts.map((post, key) => (
        <PostCard post={post} key={key} />
      ))}
    </div>
  )
}

const PostCard: FC<{ post: Post }> = ({ post }) => {
  return (
    <div
      className={
        "max-w-3xl rounded-lg border border-base-100 shadow-base-100 shadow-sm mb-4 mx-auto"
      }
    >
      <div className={"mx-5"}>
        <div className={"clickable flex flex-wrap text-xl my-3"}>
          <Link href={`/post/${post.id}`}>• {post.title} •</Link>
          <span className={"ml-auto self-center"}>
            <CommentsCount post={post} />
          </span>
        </div>
        <p className={"mb-5 mx-3 font-light"}>{post.description}</p>
      </div>

      <div className={"bg-base-100 m-1 rounded-lg flex flex-wrap py-2 px-5 gap-2"}>
        <span className={"self-center mr-3"}>
          <ReactionsButtons post={post} />
        </span>
        <span className={"flex flex-wrap mx-auto sm:m-0 gap-1"}>
          <Category post={post} />
        </span>
        <span className={"ml-auto font-normal"}>
          <span className={"text-info text-sm align-middle"}>
            <AutoDate date={post.date} />
            {" by "}
          </span>
          <span className={"clickable font-Alatsi align-middle"}>
            <Link href={`/user/${post.author.id}`}>{post.author.name}</Link>
          </span>
        </span>
      </div>
    </div>
  )
}
