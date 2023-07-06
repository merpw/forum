"use client"

import { FC, useEffect, useState } from "react"
import ReactTextareaAutosize from "react-textarea-autosize"

import { useMe } from "@/api/auth/hooks"
import { CreateComment, useComments } from "@/api/posts/comment"
import { FormError } from "@/components/error"

const CommentForm: FC<{ postId: number }> = ({ postId }) => {
  const { isLoggedIn, isLoading } = useMe()
  const { mutate: mutateComments } = useComments(postId)

  const [text, setText] = useState("")
  const [formError, setFormError] = useState<string | null>(null)

  const [isSame, setIsSame] = useState(false)

  useEffect(() => setIsSame(false), [text])

  if (!isLoggedIn && !isLoading) return null

  return (
    <div className={"form-control"}>
      <form
        onSubmit={(e) => {
          e.preventDefault()

          if (isSame) return

          setIsSame(true)
          if (formError != null) setFormError(null)

          CreateComment(postId, text)
            .then(() => {
              setText("")
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
        className={"mb-5"}
      >
        <div className={"mb-3"}>
          <label htmlFor={"comment-text"} className={"block mb-2 font-light"}>
            Write a comment:
          </label>

          <ReactTextareaAutosize
            id={"comment-text"}
            className={"textarea textarea-bordered w-full"}
            value={text}
            onInput={(e) => setText(e.currentTarget.value)}
            required
            rows={1}
          />
        </div>

        <FormError error={formError} />

        <div className={"text-center"}>
          <button type={"submit"} className={"btn button"}>
            Submit
          </button>
        </div>
      </form>
    </div>
  )
}

export default CommentForm
