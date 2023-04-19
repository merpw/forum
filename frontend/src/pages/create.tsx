import { GetStaticProps, NextPage } from "next"
import { useRouter } from "next/router"
import { FC, useEffect, useId, useState } from "react"
import ReactTextAreaAutosize from "react-textarea-autosize"
import { NextSeo } from "next-seo"
import Select from "react-select"

import { useMe } from "@/api/auth"
import { CreatePost } from "@/api/posts/create"
import { FormError } from "@/components/error"
import { getCategoriesLocal } from "@/api/posts/fetch"
import { Capitalize } from "@/helpers/text"

const CreatePostPage: NextPage<{ categories: string[] }> = ({ categories }) => {
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

      <CreatePostForm categories={categories} />
    </>
  )
}

const CreatePostForm: FC<{ categories: string[] }> = ({ categories }) => {
  const [formFields, setFormFields] = useState<{
    title: string
    content: string
    description: string
    categories: string[]
  }>({ title: "", content: "", description: "", categories: [] })

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
      onSubmit={(e) => {
        e.preventDefault()

        if (isSame) return

        if (formError != null) setFormError(null)
        setIsSame(true)
        if (formFields.categories.length == 0) {
          setFormError("Category is not selected")
          return
        }

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

export const getStaticProps: GetStaticProps = async () => {
  if (!process.env.FORUM_BACKEND_PRIVATE_URL) {
    return { notFound: true, revalidate: 60 }
  }
  const categories = await getCategoriesLocal()
  return {
    props: {
      categories,
    },
    revalidate: 60,
  }
}

export default CreatePostPage
