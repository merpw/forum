import { AxiosError } from "axios"
import { NextPage } from "next"
import Head from "next/head"
import Link from "next/link"
import { useRouter } from "next/router"
import { useEffect, useState } from "react"
import { logIn, useMe } from "../api/auth"
import { FormError } from "../components/error"

const LoginPage: NextPage = () => {
  const router = useRouter()

  const { isLoading, isLoggedIn, mutate } = useMe()

  const [login, setLogin] = useState("")
  const [password, setPassword] = useState("")
  const [formError, setFormError] = useState<string | null>(null)

  const [isSame, setIsSame] = useState(false)

  const [isRedirecting, setIsRedirecting] = useState(false) // Prevents duplicated redirects
  useEffect(() => {
    if (!isLoading && isLoggedIn && !isRedirecting) {
      setIsRedirecting(true)
      router.replace("/me")
    }
  }, [router, isLoggedIn, isRedirecting, isLoading])

  useEffect(() => {
    setIsSame(false)
  }, [login, password])

  return (
    <>
      <Head>
        <title>Login - Forum</title>
        <meta name={"title"} content={"Login - Forum"} />
        <meta name={"og:title"} content={"Login - Forum"} />
      </Head>
      <form
        onSubmit={async (e) => {
          e.preventDefault()

          if (isSame) return

          if (formError != null) setFormError(null)
          setIsSame(true)

          logIn(login, password)
            .then(() => mutate())
            .catch((err: AxiosError) => {
              setFormError(err.response?.data as string)
            })
        }}
      >
        <div className={"mb-6"}>
          <label
            htmlFor={"login"}
            className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
          >
            Your email or username
            <input
              type={login.match("@") ? "email" : "text"}
              className={
                "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
              }
              onInput={(e) => setLogin(e.currentTarget.value)}
              placeholder={"Email or username"}
              required
            />
          </label>
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

        <FormError error={formError} />
        <span className={"flex flex-wrap gap-2"}>
          <span>
            <button
              type={"submit"}
              className={
                "mb-2 text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
              }
            >
              Submit
            </button>
          </span>

          <span className={"ml-auto text-right"}>
            <div className={"font-light"}>{"Don't have an account yet?"}</div>
            <Link className={"text-2xl hover:opacity-50 "} href={"/signup"}>
              Sign Up!
            </Link>
          </span>
        </span>
      </form>
    </>
  )
}

export default LoginPage
