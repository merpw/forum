"use client"

import { FC, useEffect, useRef, useState } from "react"

import { useMe } from "@/api/auth/hooks"
import { CreateComment, useComments } from "@/api/posts/comment"
import { FormError } from "@/components/error"
import MarkdownEditor from "@/components/markdown/editor/Editor"

const CommentForm: FC<{ postId: number }> = ({ postId }) => {
  const { isLoggedIn, isLoading } = useMe()
  const { mutate: mutateComments } = useComments(postId)

  const [input, setInput] = useState("")
  const [formError, setFormError] = useState<string | null>(null)

  const [isSame, setIsSame] = useState(false)

  useEffect(() => setIsSame(false), [input])

  const formRef = useRef<HTMLFormElement>(null)

  if (!isLoggedIn && !isLoading) return null

  return (
    <div className={"form-control"}>
      <form
        ref={formRef}
        onSubmit={(e) => {
          e.preventDefault()

          if (isSame) return

          setIsSame(true)
          if (formError != null) setFormError(null)

          CreateComment(postId, input)
            .then(() => {
              setInput("")
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
            id={"comment-text"}
            className={"input input-bordered w-full py-3 px-3 pr-10"}
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter" && !e.shiftKey) {
                e.preventDefault()
                const trimmedInput = input.trim()
                setInput(trimmedInput)
                if (trimmedInput === "") {
                  return
                }
                formRef.current?.dispatchEvent(
                  new Event("submit", { cancelable: true, bubbles: true })
                )
              }
            }}
            required
            maxRows={10}
          />

          <button
            className={
              "absolute z-10 clickable disabled:opacity-50 right-3 bottom-5" +
              " " +
              (input === "" ? " btn-disabled opacity-50" : "")
            }
            type={"submit"}
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
        </div>

        <FormError error={formError} />
      </form>
    </div>
  )
}

export default CommentForm
