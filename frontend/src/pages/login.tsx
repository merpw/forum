import { motion } from "framer-motion"
import { NextPage } from "next"
import { useRouter } from "next/router"
import { useEffect, useState } from "react"
import { logIn, useUser } from "../fetch/user"

const LoginPage: NextPage = () => {
  const router = useRouter()

  const { isLoggedIn, mutate } = useUser()

  const [login, setLogin] = useState("")
  const [password, setPassword] = useState("")
  const [formError, setFormError] = useState<string | null>(null)

  const [isRedirecting, setIsRedirecting] = useState(false) // Prevents duplicated redirects
  useEffect(() => {
    if (isLoggedIn && !isRedirecting) {
      setIsRedirecting(true)
      router.push("/me")
    }
  }, [router, isLoggedIn, isRedirecting])

  return (
    <div>
      <form
        onSubmit={async (e) => {
          e.preventDefault()
          if (formError != null) setFormError(null)

          logIn(login, password)
            .then(() => mutate())
            .catch((error) => setFormError(error))
        }}
      >
        <div className={"mb-6"}>
          <label
            htmlFor={"login"}
            className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
          >
            Your email or username
          </label>
          <input
            type={login.match("@") ? "email" : "text"}
            id={"login"}
            className={
              "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            }
            onInput={(e) => setLogin(e.currentTarget.value)}
            placeholder={"Email or username"}
            required
          />
        </div>
        <div className={"mb-6"}>
          <label
            htmlFor={"password"}
            className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
          >
            Your password
          </label>
          <input
            onInput={(e) => setPassword(e.currentTarget.value)}
            type={"password"}
            id={"password"}
            className={
              "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            }
            required
          />
        </div>
        {/* <div className={"flex items-start mb-6"}>
          <div className={"flex items-center h-5"}>
            <input
              id={"remember"}
              type={"checkbox"}
              value={""}
              className={
                "w-4 h-4 bg-gray-50 rounded border border-gray-300 focus:ring-3 focus:ring-blue-300 dark:bg-gray-700 dark:border-gray-600 dark:focus:ring-blue-600 dark:ring-offset-gray-800"
              }
              required
            />
          </div>
          <label
            htmlFor={"remember"}
            className={"ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"}
          >
            Remember me
          </label>
        </div> */}
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
    </div>
  )
}

export default LoginPage
