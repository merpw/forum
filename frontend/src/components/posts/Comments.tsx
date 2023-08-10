"use client"

import { FC } from "react"
import { SWRConfig } from "swr"

import { Comment } from "@/custom"
import { useComments } from "@/api/posts/comment"
import { ReactionsCommentButtons } from "@/components/posts/reactions"
import AutoDate from "@/components/AutoDate"
import UserLink from "@/components/UserLink"
import Markdown from "@/components/markdown/markdown"

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
      <div className={"mx-5 mb-3"}>
        <div className={"clickable flex flex-wrap text-xl my-1"}>
          <UserLink userId={comment.author.id}>
            <h3 className={"clickable font-Alatsi text-base"}>
              {comment.author.username}{" "}
              <span className={"font-sans font-bold text-xl gradient-text"}>:</span>
            </h3>
          </UserLink>
        </div>
        <Markdown className={"ml-5"} content={comment.content} />
      </div>
      <div className={"bg-base-100 m-1 rounded-lg flex flex-wrap py-1.5 px-3 gap-x-1"}>
        <span className={"self-center mr-3"}>
          <ReactionsCommentButtons postId={postId} comment={comment} />
        </span>

        <AutoDate
          date={comment.date}
          className={"self-center ml-auto font-light text-sm text-info"}
        />
      </div>
    </div>
  )
}

export default Comments
