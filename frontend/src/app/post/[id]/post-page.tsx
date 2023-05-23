import Link from "next/link"
import { FC } from "react"

import { Comment, Post } from "@/custom"
import { Category, ReactionsButtons } from "@/components/posts/reactions"
import CommentForm from "@/app/post/[id]/CommentForm"
import Comments from "@/app/post/[id]/Comments"
import AutoDate from "@/components/AutoDate"

import "highlight.js/styles/github-dark.css"

const PostPage: FC<{ post: Post; comments: Comment[] }> = ({ post, comments }) => {
  return (
    <div className={"m-5"}>
      <div className={"mb-3"}>
        <h1 className={"text-3xl mb-2 "}>{post.title}</h1>
        <hr />
      </div>
      <div
        className={"prose dark:prose-invert"}
        dangerouslySetInnerHTML={{ __html: post.content }}
      />
      <hr className={"mt-4"} />
      <div className={"border-t py-2 flex flex-wrap"}>
        <ReactionsButtons post={post} />
        <Category post={post} />

        <span className={"ml-auto"}>
          <AutoDate date={post.date} />
          {" by "}
          <span className={"clickable text-xl"}>
            <Link href={`/user/${post.author.id}`}>{post.author.name}</Link>
          </span>
        </span>
      </div>
      <div>
        <h2 className={"text-2xl my-4"}>Comments:</h2>
        <CommentForm postId={post.id} />
        <Comments postId={post.id} fallback={comments} />
      </div>
    </div>
  )
}

export default PostPage
