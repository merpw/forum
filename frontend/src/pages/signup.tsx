import { AxiosError } from "axios"
import { NextPage } from "next"
import Head from "next/head"
import { useRouter } from "next/router"
import { useEffect, useState } from "react"
import { logIn, SignUp, useMe } from "../api/auth"
import { FormError } from "../components/error"

const SignupPage: NextPage = () => {
  const router = useRouter()

  const { isLoading, isLoggedIn, mutate } = useMe()

  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [passwordConfirm, setPasswordConfirm] = useState("")

  const [isSame, setIsSame] = useState(false)

  const [formError, setFormError] = useState<string | null>(null)

  const [isRedirecting, setIsRedirecting] = useState(false) // Prevents duplicated redirects
  useEffect(() => {
    if (!isLoading && isLoggedIn && !isRedirecting) {
      setIsRedirecting(true)
      router.replace("/me")
    }
  }, [router, isLoggedIn, isRedirecting, isLoading])

  useEffect(() => {
    setIsSame(false)
  }, [name, email, password, passwordConfirm])

  return (
    <>
      <Head>
        <title>Sign Up - Forum</title>
        <meta name={"title"} content={"Sign Up - Forum"} />
        <meta name={"og:title"} content={"Sign Up - Forum"} />
      </Head>
      <form
        onSubmit={async (e) => {
          e.preventDefault()

          if (isSame) return

          if (formError != null) setFormError(null)
          setIsSame(true)

          if (password != passwordConfirm) {
            setFormError("Passwords don't match")
            return
          }

          SignUp(name, email, password)
            .then((response) => {
              if (response.status == 200) {
                logIn(email, password).then(() => mutate())
              }
            })
            .catch((err: AxiosError) => {
              if (err.code == "ERR_BAD_REQUEST") {
                setFormError(err.response?.data as string)
              } else {
                // TODO: unexpected error
              }
            })
        }}
      >
        <div className={"mb-6"}>
          <label
            htmlFor={"username"}
            className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
          >
            Username
          </label>
          <input
            type={"text"}
            id={"username"}
            className={
              "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            }
            onInput={(e) => setName(e.currentTarget.value)}
            placeholder={"Username"}
            required
          />
        </div>
        <div className={"mb-6"}>
          <label
            htmlFor={"email"}
            className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
          >
            Your email
          </label>
          <input
            type={"email"}
            id={"email"}
            className={
              "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            }
            onInput={(e) => setEmail(e.currentTarget.value)}
            placeholder={"Email"}
            required
          />
        </div>
        <div className={"mb-3"}>
          <label
            htmlFor={"password"}
            className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
          >
            Create password
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
        <div className={"mb-6"}>
          <label
            htmlFor={"repeat-password"}
            className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
          >
            Repeat password
          </label>
          <input
            onInput={(e) => setPasswordConfirm(e.currentTarget.value)}
            type={"password"}
            id={"repeat-password"}
            className={
              "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            }
            required
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
    </>
  )
}

export default SignupPage
