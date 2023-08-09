"use client"

import { FC, useRef, useState } from "react"

import { useMe } from "@/api/auth/hooks"
import { CreateComment, useComments } from "@/api/posts/comment"
import { FormError } from "@/components/error"
import MarkdownEditor from "@/components/markdown/editor/Editor"

const CommentForm: FC<{ postId: number }> = ({ postId }) => {
  const { isLoggedIn, isLoading } = useMe()
  const { mutate: mutateComments } = useComments(postId)

  const [formError, setFormError] = useState<string | null>(null)

  const [isSame, setIsSame] = useState(false)

  const formRef = useRef<HTMLFormElement>(null)

  if (!isLoggedIn && !isLoading) return null

  return (
    <div className={"form-control"}>
      <form
        onChange={() => setIsSame(false)}
        ref={formRef}
        onSubmit={(e) => {
          e.preventDefault()

          if (isSame) return

          setIsSame(true)
          if (formError != null) setFormError(null)

          const commentTextarea = formRef.current?.querySelector(
            "[name=comment-text]"
          ) as HTMLTextAreaElement
          const input = commentTextarea.value

          CreateComment(postId, input)
            .then(() => {
              commentTextarea.value = ""
              mutateComments()
            })
            .catch((err) => {
              if (err.code == "ERR_BAD_REQUEST") {
                setFormError(err.response?.data as string)
              } else {
                // TODO: unexpected error
              }
            })
        }}
      >
        <div className={"bg-base-200 rounded flex flex-col p-2 relative"}>
          <label htmlFor={"comment-text"} className={"m-2 font-light"}>
            Leave a comment:
          </label>

          <MarkdownEditor
            name={"comment-text"}
            withSubmit
            id={"comment-text"}
            className={"input input-bordered w-full"}
            required
            maxRows={10}
          />
        </div>

        <FormError error={formError} />
      </form>
    </div>
  )
}

export default CommentForm
