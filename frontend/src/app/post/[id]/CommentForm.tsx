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
        className={"relative mb-7"}
      >
        <div className={"mb-3"}>
          <label htmlFor={"comment-text"} className={"block mb-2 font-light"}>
            Leave a comment:
          </label>

          <ReactTextareaAutosize
            id={"comment-text"}
            className={"absolute input input-bordered w-full py-3 px-3 pr-10"}
            value={text}
            onInput={(e) => setText(e.currentTarget.value)}
            required
            maxRows={2}
          />
        </div>

        <FormError error={formError} />

        <button
          type={"submit"}
          className={"absolute z-10 clickable disabled:opacity-50 right-1.5 m-1.5 mt-1.5"}
        >
          <svg
            xmlns={"http://www.w3.org/2000/svg"}
            fill={"none"}
            viewBox={"0 0 24 24"}
            strokeWidth={2}
            stroke={"currentColor"}
            className={"w-7 h-7 text-primary"}
          >
            <path
              strokeLinecap={"round"}
              strokeLinejoin={"round"}
              d={
                "M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"
              }
            />
          </svg>
        </button>
      </form>
    </div>
  )
}

export default CommentForm
