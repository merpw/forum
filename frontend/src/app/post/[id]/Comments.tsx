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
    return (
      <div className={"text-info text-center mt-3 mb-7"}>
        There are no comments yet, be the first!
      </div>
    )
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
    <div className={"max-w-3xl rounded-lg border border-base-100 shadow-base-100 shadow-sm"}>
      <div className={"mx-5"}>
        <h1 className={"clickable flex flex-wrap text-xl my-1"}>
          <Link href={`/user/${comment.author.id}`}>
            <h3 className={"clickable font-Alatsi text-base"}>
              {comment.author.name}{" "}
              <span className={"font-sans font-bold text-xl gradient-text"}>:</span>
            </h3>
          </Link>
        </h1>
        <p className={"whitespace-pre-line mb-3 ml-5 font-light"}>{comment.content}</p>
      </div>
      <div className={"bg-base-100 m-1 rounded-lg flex flex-wrap py-1.5 px-3 gap-x-1"}>
        <span className={"self-center mr-3"}>
          <ReactionsCommentButtons postId={postId} comment={comment} />
        </span>

        <AutoDate
          date={comment.date}
          className={"self-center ml-auto text-sm font-light text-info"}
        />
      </div>
    </div>
  )
}

export default Comments
