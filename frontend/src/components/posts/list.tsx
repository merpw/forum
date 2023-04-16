import Link from "next/link"
import { FC } from "react"

import { Category, CommentsCount, ReactionsButtons } from "./reactions"

import { Post } from "@/custom"
import useDates from "@/helpers/dates"

export const PostList: FC<{ posts: Post[] }> = ({ posts }) => {
  if (posts.length == 0) {
    return <div>There are no posts yet...</div>
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
  const { localDate, relativeDate } = useDates(post.date)

  return (
    <div
      className={
        "rounded-lg border opacity-90 hover:opacity-100 hover:shadow dark:shadow-white mb-4"
      }
    >
      <div className={"m-5"}>
        <div className={"mb-3"}>
          <h1 className={"clickable text-2xl max-w-fit"}>
            <Link href={`/post/${post.id}`}>{post.title}</Link>
          </h1>
          <hr className={"mt-2"} />
        </div>

        {/* TODO: change to description */}
        <p>{post.content}</p>
      </div>

      <div className={"border-t flex flex-wrap p-3 gap-y-2 "}>
        <ReactionsButtons post={post} />
        <CommentsCount post={post} />
        <Category post={post} />
        <span className={"ml-auto"}>
          <span suppressHydrationWarning title={localDate}>
            {relativeDate}
          </span>
          {" by "}
          <span className={"clickable text-xl"}>
            <Link href={`/user/${post.author.id}`}>{post.author.name}</Link>
          </span>
        </span>
      </div>
    </div>
  )
}
