"use client"

import { useRouter } from "next/navigation"
import { Dispatch, FC, SetStateAction, useId, useRef, useState } from "react"
import ReactTextAreaAutosize from "react-textarea-autosize"
import Select from "react-select"

import { CreatePost, generateDescription } from "@/api/posts/create"
import { FormError } from "@/components/error"
import { Capitalize } from "@/helpers/text"
import { MarkdownToPlain } from "@/components/markdown/render"
import MarkdownEditor from "@/components/markdown/editor"
import { trimInput } from "@/helpers/input"

const CreatePostPage: FC<{ categories: string[]; isAIEnabled: boolean }> = ({
  categories,
  isAIEnabled,
}) => {
  const [formError, setFormError] = useState<string | null>(null)

  const router = useRouter()

  const [isSame, setIsSame] = useState(false)

  const formRef = useRef<HTMLFormElement>(null)

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

        console.log(formFields)

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
            <div className={"mb-3"}>
              <Select
                placeholder={"Categories"}
                instanceId={useId()}
                isMulti={true}
                name={"categories"}
                className={"react-select-container"}
                classNamePrefix={"react-select"}
                options={categories.map((name) => ({ label: Capitalize(name), value: name }))}
                required
              />
            </div>

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

const GenerateDescriptionButton: FC<{
  formRef: React.RefObject<HTMLFormElement>
  setFormError: Dispatch<SetStateAction<string | null>>
}> = ({ formRef, setFormError }) => {
  const [isLoading, setIsLoading] = useState(false)

  return (
    <button
      onClick={async () => {
        const formData = new FormData(formRef.current as HTMLFormElement)

        const formFields = {
          title: formData.get("title") as string,
          content: formData.get("content") as string,
        }

        if (!formFields.title) {
          return setFormError("Title is required")
        }
        if (!formFields.content) {
          return setFormError("Content is required")
        }
        setIsLoading(true)
        generateDescription(formFields)
          .then((description) => {
            setFormError(null)
            const descriptionTextarea = formRef.current?.querySelector(
              "[name=description]"
            ) as HTMLInputElement
            descriptionTextarea.value = description
          })
          .catch((err) => setFormError(err.message))
          .finally(() => setIsLoading(false))
      }}
      type={"button"}
      className={
        "btn btn-sm transition-none hover:opacity-100 hover:gradient-text hover:border-primary btn-outline font-normal mb-3 self-center font-xs"
      }
    >
      {isLoading ? (
        <span className={"text-primary loading loading-ring"} />
      ) : (
        <svg
          xmlns={"http://www.w3.org/2000/svg"}
          fill={"none"}
          viewBox={"0 0 24 24"}
          strokeWidth={1}
          stroke={"currentColor"}
          className={"w-5 h-5 mr-1 fill-primary"}
        >
          <path
            strokeLinecap={"round"}
            strokeLinejoin={"round"}
            d={
              "M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 00-2.456 2.456zM16.894 20.567L16.5 21.75l-.394-1.183a2.25 2.25 0 00-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 001.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 001.423 1.423l1.183.394-1.183.394a2.25 2.25 0 00-1.423 1.423z"
            }
          />
        </svg>
      )}
      Generate description
    </button>
  )
}

export default CreatePostPage
