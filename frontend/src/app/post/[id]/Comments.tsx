"use client"

import { FC } from "react"
import Link from "next/link"
import { SWRConfig } from "swr"

import { Comment } from "@/custom"
import { useComments } from "@/api/posts/comment"
import { ReactionsCommentButtons } from "@/components/posts/reactions"
import AutoDate from "@/components/AutoDate"

const Comments: FC<{ postId: number; fallback: Comment[] }> = ({ postId, fallback }) => {
  return (
    <SWRConfig
      value={{
        fallback: {
          [`/api/posts/${postId}/comments`]: fallback,
        },
      }}
    >
      <CommentList postId={postId} />
    </SWRConfig>
  )
}

const CommentList: FC<{ postId: number }> = ({ postId }) => {
  const { comments } = useComments(postId)
  if (!comments) return null

  if (comments.length == 0) {
    return <div>There are no comments yet, write one first!</div>
  }

  return (
    <div className={"flex flex-col gap-3"}>
      {comments
        .sort((a, b) => b.date.localeCompare(a.date))
        .map((comment, key) => (
          <CommentCard comment={comment} postId={postId} key={key} />
        ))}
    </div>
  )
}

const CommentCard: FC<{ comment: Comment; postId: number }> = ({ comment, postId }) => {
  return (
    <div className={"border rounded p-5"}>
      <Link href={`/user/${comment.author.id}`}>
        <h3 className={"clickable text-lg"}>{comment.author.name}</h3>
      </Link>
      <p className={"whitespace-pre-line"}>{comment.content}</p>
      <hr className={"mt-4 mb-2"}></hr>
      <span className={"flex"}>
        <ReactionsCommentButtons postId={postId} comment={comment} />

        <AutoDate date={comment.date} className={"ml-auto"} />
      </span>
    </div>
  )
}

export default Comments
