"use client"

import { useRouter } from "next/navigation"
import { FC, useEffect, useId, useState } from "react"
import ReactTextAreaAutosize from "react-textarea-autosize"
import Select from "react-select"
import { remark } from "remark"
import stripMarkdown from "strip-markdown"
import dynamic from "next/dynamic"

import { useMe } from "@/api/auth"
import { CreatePost, generateDescription } from "@/api/posts/create"
import { FormError } from "@/components/error"
import { Capitalize } from "@/helpers/text"

const Markdown = dynamic(() => import("@/components/markdown/markdown"))

const CreatePostPage: FC<{ categories: string[]; isAIEnabled: boolean }> = ({
  categories,
  isAIEnabled,
}) => {
  const { isLoading, isLoggedIn } = useMe()
  const router = useRouter()

  const [isRedirecting, setIsRedirecting] = useState(false) // Prevents duplicated redirects
  useEffect(() => {
    if (!isLoading && !isLoggedIn && !isRedirecting) {
      setIsRedirecting(true)
      router.push("/login")
    }
  }, [router, isLoggedIn, isRedirecting, isLoading])

  return <CreatePostForm categories={categories} isAIEnabled={isAIEnabled} />
}

const CreatePostForm: FC<{ categories: string[]; isAIEnabled: boolean }> = ({
  categories,
  isAIEnabled,
}) => {
  const [formFields, setFormFields] = useState<{
    title: string
    content: string
    description: string
    isDescriptionLoading: boolean
    categories: string[]
  }>({ title: "", content: "", description: "", isDescriptionLoading: false, categories: [] })

  const [formError, setFormError] = useState<string | null>(null)

  const router = useRouter()

  const [isSame, setIsSame] = useState(false)

  useEffect(() => {
    setIsSame(false)
  }, [formFields])

  return (
    <form
      onChange={(e) => {
        const { name, value } = e.target as HTMLInputElement
        setFormFields({ ...formFields, [name]: value })
      }}
      onBlur={() => {
        formFields.title = formFields.title.trim()
        formFields.content = formFields.content.trim()
        formFields.description = formFields.description.trim()
        setFormFields({ ...formFields })
      }}
      onSubmit={async (e) => {
        e.preventDefault()

        if (isSame) return

        if (formError != null) setFormError(null)
        setIsSame(true)

        if (formFields.description == "") {
          // "stupid" description generation, converts content to plain text and cuts it to 200 characters
          const description = await remark().use(stripMarkdown).process(formFields.content)
          formFields.description = description.toString().replace(/\n+/g, " ").trim()
          if (formFields.description.length > 200) {
            formFields.description = formFields.description.slice(0, 200) + "..."
          }
          setFormFields({ ...formFields })
        }

        if (document.querySelector("#preview h1")) {
          return setFormError(
            "The biggest headings are not allowed, please use title field instead"
          )
        }

        CreatePost(formFields)
          .then((id) => router.push(`/post/${id}`))
          .catch((err) => setFormError(err.message))
      }}
    >
      <div className={"max-w-3xl mx-auto flex-col"}>
        <div className={"card flex-shrink-0 bg-neutral shadow-md"}>
          <div className={"card-body"}>
            <div className={"my-2 font-Yesteryear text-3xl text-primary opacity-50 text-center"}>
              <span className={"text-neutral-content"}> {"What's on your mind?"}</span>
            </div>
            <div className={"form-control"}>
              <input
                type={"text"}
                name={"title"}
                className={"input input-bordered bg-base-100 text-sm"}
                placeholder={"Title"}
                value={formFields.title}
                onChange={() => void 0 /* handled by Form */}
                required
                maxLength={25}
              />
            </div>
            <div className={"form-control"}>
              <ReactTextAreaAutosize
                title={
                  "Content input field. Press ESC and then TAB to move the focus to the next field"
                }
                onFocus={(e) => e.currentTarget.setAttribute("data-escaped", "false")}
                onKeyDown={(e) => {
                  if (e.key === "Escape") {
                    return e.currentTarget.setAttribute("data-escaped", "true")
                  }

                  if (
                    e.currentTarget.getAttribute("data-escaped") === "false" &&
                    !e.shiftKey &&
                    e.key === "Tab"
                  ) {
                    e.preventDefault()
                    const textarea = e.target as HTMLTextAreaElement
                    const start = textarea.selectionStart
                    const end = textarea.selectionEnd
                    const value = textarea.value
                    textarea.value = value.substring(0, start) + "\t" + value.substring(end)
                    textarea.selectionStart = textarea.selectionEnd = start + 1
                    setFormFields({ ...formFields, content: textarea.value })
                  }
                }}
                name={"content"}
                className={"textarea textarea-bordered"}
                rows={5}
                minRows={5}
                placeholder={"Content"}
                value={formFields.content}
                required
                maxLength={10000}
              />
            </div>

            <div className={"mb-3 border border-base-100 rounded-xl mt-3 px-4"} id={"preview"}>
              <Markdown content={formFields.content} />
            </div>
            <ReactTextAreaAutosize
              name={"description"}
              className={"textarea textarea-bordered"}
              rows={2}
              minRows={2}
              placeholder={"Description"}
              value={formFields.description}
              maxLength={205}
            />
            {isAIEnabled && (
              <button
                onClick={async () => {
                  if (!formFields.title) {
                    return setFormError("Title is required")
                  }
                  if (!formFields.content) {
                    return setFormError("Content is required")
                  }
                  setFormFields({ ...formFields, isDescriptionLoading: true })
                  generateDescription(formFields)
                    .then((description) => {
                      setFormError(null)
                      setFormFields((prev) => ({ ...prev, description }))
                    })
                    .catch((err) => setFormError(err.message))
                    .finally(() =>
                      setFormFields((prev) => ({ ...prev, isDescriptionLoading: false }))
                    )
                }}
                type={"button"}
                className={
                  "btn btn-sm transition-none hover:opacity-100 hover:gradient-text hover:border-primary btn-outline font-normal mb-3 self-center font-xs"
                }
              >
                {formFields.isDescriptionLoading ? (
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
            )}

            <div className={"mb-3"}>
              <Select
                placeholder={"Categories"}
                instanceId={useId()}
                isMulti={true}
                name={"categories"}
                className={"react-select-container"}
                classNamePrefix={"react-select"}
                onChange={(newValue) =>
                  setFormFields({ ...formFields, categories: newValue.map((v) => v.value) })
                }
                options={categories.map((name) => ({ label: Capitalize(name), value: name }))}
                required
              />
            </div>

            <FormError error={formError} />
            <div className={"form-control"}>
              <button type={"submit"} className={"button self-center"}>
                SUBMIT
              </button>
            </div>
          </div>
        </div>
      </div>
    </form>
  )
}

export default CreatePostPage
