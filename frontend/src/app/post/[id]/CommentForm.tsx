"use client"

import { FC, useEffect, useState } from "react"
import ReactTextareaAutosize from "react-textarea-autosize"

import { useMe } from "@/api/auth"
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
        <label
          htmlFor={"comment-text"}
          className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
        >
          Write a comment:
        </label>

        <ReactTextareaAutosize
          id={"comment-text"}
          className={"inputbox"}
          value={text}
          onInput={(e) => setText(e.currentTarget.value)}
          required
          rows={1}
        />
      </div>

      <FormError error={formError} />

      <button type={"submit"} className={"button"}>
        Submit
      </button>
    </form>
  )
}

export default CommentForm
