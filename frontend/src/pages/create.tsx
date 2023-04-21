import { GetServerSideProps, NextPage } from "next"
import { useRouter } from "next/router"
import { FC, useEffect, useId, useState } from "react"
import ReactTextAreaAutosize from "react-textarea-autosize"
import { NextSeo } from "next-seo"
import Select from "react-select"
import { remark } from "remark"
import stripMarkdown from "strip-markdown"
import dynamic from "next/dynamic"

import { useMe } from "@/api/auth"
import { CreatePost, generateDescription } from "@/api/posts/create"
import { FormError } from "@/components/error"
import { getCategoriesLocal } from "@/api/posts/fetch"
import { Capitalize } from "@/helpers/text"

const Markdown = dynamic(() => import("@/components/markdown"))

const CreatePostPage: NextPage<{ categories: string[]; isAIEnabled: boolean }> = ({
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

  return (
    <>
      <NextSeo title={"Create new post"} />

      <CreatePostForm categories={categories} isAIEnabled={isAIEnabled} />
    </>
  )
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
      onSubmit={async (e) => {
        e.preventDefault()

        if (isSame) return

        if (formError != null) setFormError(null)
        setIsSame(true)
        if (formFields.categories.length == 0) {
          setFormError("Category is not selected")
          return
        }
        if (formFields.description == "") {
          const description = await remark().use(stripMarkdown).process(formFields.content)
          formFields.description = description.toString().replace(/\n+/g, " ").trim()
        }

        if (document.querySelector("#preview h1")) {
          return setFormError(
            "The biggest headings are not allowed, please use title field instead"
          )
        }

        setFormFields({ ...formFields }) // to show generated description in the input field

        CreatePost(formFields)
          .then((id) => router.push(`/post/${id}`))
          .catch((err) => {
            if (err.code == "ERR_BAD_REQUEST") {
              setFormError(err.response?.data as string)
            } else {
              // TODO: unexpected error
            }
          })
      }}
    >
      <div className={"mb-3"}>
        <input
          type={"text"}
          name={"title"}
          className={
            "bg-gray-50 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-700 dark:border-gray-600 "
          }
          placeholder={"Title"}
          required
        />
      </div>
      <div className={"mb-3"}>
        <ReactTextAreaAutosize
          title={"Content input field. Press ESC and then TAB to move the focus to the next field"}
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
          className={
            "bg-gray-50 border border-gray-300 text-gray-900 rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
          }
          rows={5}
          minRows={5}
          placeholder={"Content"}
          required
        />
      </div>

      <div className={"mb-3"} id={"preview"}>
        <Markdown content={formFields.content} />
      </div>

      <ReactTextAreaAutosize
        name={"description"}
        className={
          "mb-2 bg-gray-50 border border-gray-300 text-gray-900 rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
        }
        rows={2}
        minRows={2}
        placeholder={"Description"}
        value={formFields.description}
      />
      {isAIEnabled && (
        <button
          onClick={async () => {
            setFormFields({ ...formFields, isDescriptionLoading: true })
            generateDescription(formFields)
              .then((description) => {
                setFormError(null)
                setFormFields((prev) => ({ ...prev, description }))
              })
              .catch((err) => {
                setFormError(err.message)
              })
              .finally(() => setFormFields((prev) => ({ ...prev, isDescriptionLoading: false })))
          }}
          type={"button"}
          className={
            "mb-3 flex flex-row row justify-center text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-2.5 py-2 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
          }
        >
          {formFields.isDescriptionLoading ? (
            <svg
              aria-hidden={"true"}
              className={"w-5 h-5 mr-1 text-gray-200 animate-spin fill-blue-600"}
              viewBox={"0 0 100 101"}
              fill={"none"}
              xmlns={"http://www.w3.org/2000/svg"}
            >
              <path
                d={
                  "M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
                }
                fill={"currentColor"}
              />
              <path
                d={
                  "M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
                }
                fill={"currentFill"}
              />
            </svg>
          ) : (
            <svg
              xmlns={"http://www.w3.org/2000/svg"}
              fill={"none"}
              viewBox={"0 0 24 24"}
              strokeWidth={1.5}
              stroke={"currentColor"}
              className={"w-5 h-5 mr-1"}
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
        />
      </div>
      <FormError error={formError} />

      <button
        type={"submit"}
        className={
          "text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
        }
      >
        Submit
      </button>
    </form>
  )
}

export const getServerSideProps: GetServerSideProps = async () => {
  const categories = await getCategoriesLocal()
  return {
    props: {
      categories,
      isAIEnabled: !!process.env.OPENAI_API_KEY,
    },
  }
}

export default CreatePostPage
