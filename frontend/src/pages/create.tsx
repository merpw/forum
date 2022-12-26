import { motion } from "framer-motion"
import { NextPage } from "next"
import Head from "next/head"
import { useRouter } from "next/router"
import { useEffect, useState } from "react"
import { useUser } from "../api/auth"
import { CreatePost } from "../api/posts/create"

const CreatePostPage: NextPage = () => {
  const { isLoading, isLoggedIn } = useUser()
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
      <Head>
        <title>Create new post - Forum</title>
      </Head>
      <CreatePostForm />
    </>
  )
}

const CreatePostForm = () => {
  const [title, setTitle] = useState("")
  const [content, setContent] = useState("")
  // const [categories, setCategories] = useState([])
  const [formError, setFormError] = useState<string | null>(null)
  const router = useRouter()

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault()
        setFormError(null)
        CreatePost(title, content)
          .then((id) => {
            router.push(`/post/${id}`)
          })
          .catch((err) => {
            setFormError(err)
          })
      }}
    >
      <div className={"mb-3"}>
        <input
          type={"text"}
          title={"title"}
          className={
            "bg-gray-50 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-700 dark:border-gray-600 "
          }
          onInput={(e) => setTitle(e.currentTarget.value)}
          placeholder={"Title"}
          required
        />
      </div>
      <div className={"mb-6"}>
        <textarea
          title={"content"}
          className={
            "bg-gray-50 border border-gray-300 text-gray-900 rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
          }
          onInput={(e) => setContent(e.currentTarget.value)}
          placeholder={"Content"}
          required
        />
      </div>
      {formError != null && (
        <motion.div
          className={
            "transition ease-in-out -translate-y-1 p-4 mb-4 text-sm text-red-700 bg-red-100 rounded-lg dark:bg-inherit dark:border-2 dark:border-red-900 dark:text-white"
          }
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ ease: "easeOut", duration: 0.25 }}
          exit={{ opacity: 0 }}
          role={"alert"}
        >
          <span className={"font-medium"}>{formError}</span>
        </motion.div>
      )}
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

export default CreatePostPage
