import { NextPage } from "next"
import Link from "next/link"
import { useRouter } from "next/router"
import { useEffect, useState } from "react"
import { NextSeo } from "next-seo"

import { logIn, useMe } from "@/api/auth"
import { FormError } from "@/components/error"

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
      <NextSeo title={"Login"} />
      <form
        onSubmit={async (e) => {
          e.preventDefault()

          if (isSame) return

          if (formError != null) setFormError(null)
          setIsSame(true)

          logIn(login, password)
            .then(() => mutate())
            .catch((err) => setFormError(err.message))
        }}
      >
        <div className={"mb-6"}>
          <label className={"label"}>
            <p className={"inputbox-title"}>Your email or username</p>
            <input
              type={login.match("@") ? "email" : "text"}
              className={"inputbox-singlerow"}
              value={login}
              onChange={(e) => setLogin(e.currentTarget.value.trim())}
              placeholder={"Email or username"}
              required
            />
          </label>
        </div>
        <div className={"mb-6"}>
          <label className={"label"}>
            <p className={"inputbox-title"}>Your password</p>
            <input
              onChange={(e) => setPassword(e.currentTarget.value)}
              type={"password"}
              id={"password"}
              className={"inputbox-singlerow"}
              required
            />
          </label>
        </div>

        <FormError error={formError} />
        <span className={"flex flex-wrap gap-2"}>
          <span>
            <button type={"submit"} className={"button mb-2"}>
              Submit
            </button>
          </span>

          <span className={"ml-auto text-right"}>
            <div className={"font-light"}>{"Don't have an account yet?"}</div>
            <Link className={"clickable text-2xl"} href={"/signup"}>
              Sign Up!
            </Link>
          </span>
        </span>
      </form>
    </>
  )
}

export default LoginPage
