"use client"

import { FC, useRef, useState } from "react"
import ReactTextAreaAutosize from "react-textarea-autosize"
import { useRouter } from "next/navigation"

import { CreatePost } from "@/api/posts/create"
import { FormError } from "@/components/error"
import { MarkdownToPlain } from "@/components/markdown/render"
import MarkdownEditor from "@/components/markdown/editor/Editor"
import { trimInput } from "@/helpers/input"
import { GenerateDescriptionButton } from "@/components/forms/create-post/GenerateDescriptionButton"
import SelectCategories from "@/components/forms/create-post/SelectCategories"

const CreatePostForm: FC<{ categories: string[]; isAIEnabled: boolean }> = ({
  categories,
  isAIEnabled,
}) => {
  const formRef = useRef<HTMLFormElement>(null)

  const [isSame, setIsSame] = useState(false)

  const [formError, setFormError] = useState<string | null>(null)

  const router = useRouter()

  return (
    <form
      ref={formRef}
      onChange={() => setIsSame(false)}
      onSubmit={async (e) => {
        e.preventDefault()

        if (isSame) return

        const form = formRef.current as HTMLFormElement

        const formData = new FormData(form)

        const formFields = {
          title: formData.get("title") as string,
          content: formData.get("content") as string,
          description: formData.get("description") as string,
          categories: formData.getAll("categories") as string[],
        }

        if (formError != null) setFormError(null)
        setIsSame(true)

        if (formFields.description == "") {
          // "stupid" description generation, converts content to plain text and cuts it to 200 characters
          formFields.description = await MarkdownToPlain(formFields.content, {
            limit: 200,
            removeNewLines: true,
          })

          const descriptionTextarea = form.querySelector("[name=description]") as HTMLInputElement

          descriptionTextarea.value = formFields.description
        }

        CreatePost(formFields)
          .then((id) => router.push(`/post/${id}`))
          .catch((err) => setFormError(err.message))
      }}
    >
      <div className={"max-w-3xl mx-auto flex-col"}>
        <div className={"card flex-shrink-0 bg-base-200 shadow-md"}>
          <div className={"card-body"}>
            <div className={"my-2 font-Yesteryear text-3xl text-primary opacity-50 text-center"}>
              <span className={"text-neutral-content"}> {"What's on your mind?"}</span>
            </div>

            <div className={"form-control"}>
              <input
                onBlur={trimInput}
                type={"text"}
                name={"title"}
                className={"input input-bordered bg-base-100 text-sm"}
                placeholder={"Title"}
                required
                maxLength={25}
              />
            </div>

            <div className={"form-control"}>
              <MarkdownEditor
                title={
                  "Content input field. Press ESC and then TAB to move the focus to the next field"
                }
                name={"content"}
                className={"textarea textarea-bordered w-full"}
                rows={5}
                minRows={5}
                placeholder={"Content"}
                required
                maxLength={10000}
              />
            </div>

            <ReactTextAreaAutosize
              onBlur={trimInput}
              name={"description"}
              className={"textarea textarea-bordered"}
              rows={2}
              minRows={2}
              placeholder={"Description"}
              maxLength={205}
            />
            {isAIEnabled && (
              <GenerateDescriptionButton formRef={formRef} setFormError={setFormError} />
            )}

            <SelectCategories categories={categories} />

            <FormError error={formError} />

            <div className={"form-control"}>
              <button type={"submit"} className={"button self-center"}>
                <svg
                  xmlns={"http://www.w3.org/2000/svg"}
                  fill={"none"}
                  viewBox={"0 0 24 24"}
                  strokeWidth={1.5}
                  stroke={"currentColor"}
                  className={"w-4 h-4"}
                >
                  <path
                    strokeLinecap={"round"}
                    strokeLinejoin={"round"}
                    d={
                      "M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"
                    }
                  />
                </svg>
                SUBMIT
              </button>
            </div>
          </div>
        </div>
      </div>
    </form>
  )
}

export default CreatePostForm
