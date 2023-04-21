import { NextPage } from "next"
import { useRouter } from "next/router"
import { useEffect, useState } from "react"
import ReactTextAreaAutosize from "react-textarea-autosize"
import { NextSeo } from "next-seo"

import { useMe } from "@/api/auth"
import { getCategories } from "@/api/posts/categories"
import { CreatePost } from "@/api/posts/create"
import { FormError } from "@/components/error"

const CreatePostPage: NextPage = () => {
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

      <CreatePostForm />
    </>
  )
}

const CreatePostForm = () => {
  const [title, setTitle] = useState("")
  const [content, setContent] = useState("")
  const [category, setCategory] = useState<string[]>([])
  const [categories, setCategories] = useState<string[]>([])
  const [formError, setFormError] = useState<string | null>(null)

  const router = useRouter()

  const [isSame, setIsSame] = useState(false)

  useEffect(() => {
    setIsSame(false)
  }, [title, content, category])

  useEffect(() => {
    getCategories().then(setCategories)
  }, [])

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault()

        if (isSame) return

        if (formError != null) setFormError(null)
        setIsSame(true)
        if (category.length == 0) {
          setFormError("Category is not selected")
          return
        }

        CreatePost(title, content, category)
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
          title={"title"}
          className={"inputbox-singlerow"}
          onInput={(e) => setTitle(e.currentTarget.value)}
          placeholder={"Title"}
          required
        />
      </div>
      <div className={"mb-3"}>
        <ReactTextAreaAutosize
          title={"content"}
          className={"inputbox"}
          onInput={(e) => setContent(e.currentTarget.value)}
          minRows={5}
          placeholder={"Content"}
          required
        />
      </div>
      <div className={"mb-6"}>
        <label
          htmlFor={"cats"}
          className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
        ></label>
        <select
          title={"category"}
          multiple
          required
          id={"cats"}
          className={"inputbox text-xl capitalize"}
          onChange={(e) =>
            setCategory(Array.from(e.currentTarget.selectedOptions, (option) => option.value))
          }
        >
          {categories.map((cat, key) => (
            <option key={key}>{cat}</option>
          ))}
        </select>
      </div>
      <FormError error={formError} />

      <button type={"submit"} className={"button"}>
        Submit
      </button>
    </form>
  )
}

export default CreatePostPage
