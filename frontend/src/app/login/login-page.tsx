"use client"

import Link from "next/link"
import { useRouter } from "next/navigation"
import { useEffect, useState } from "react"

import { logIn, useMe } from "@/api/auth"
import { FormError } from "@/components/error"

const LoginPage = () => {
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
      <form
        onSubmit={async (e) => {
          e.preventDefault()

          if (isSame) return

          if (formError != null) setFormError(null)
          setIsSame(true)

          logIn(login, password)
            .then(() => {
              router.refresh()
              mutate()
            })
            .catch((err) => setFormError(err.message))
        }}
      >
        <div className={"pt-12 pb-56 md:py-52 bg-base-200"}>
          <div className={"hero-content flex-col md:flex-row-reverse min-h-full"}>
            <div className={"text-center md:text-left"}>
              <h1 className={"text-6xl mb-5 font-Yesteryear gradient-text"}>Welcome!</h1>
              <div className={"font-thin"}>
                <p>{"Don't have an account yet?"}</p>
                <Link className={"clickable text-xl font-normal"} href={"/signup"}>
                  Sign Up!
                </Link>
              </div>
            </div>
            <div className={"card flex-shrink-0 w-full max-w-sm shadow-xl bg-base-100"}>
              <div className={"card-body"}>
                <div className={"form-control"}>
                  <label className={"label"}>
                    <span className={"label-text"}>Email or Username</span>
                  </label>
                  <input
                    type={login.match("@") ? "email" : "text"}
                    placeholder={"email or username"}
                    className={"input input-bordered"}
                    value={login}
                    onChange={(e) => setLogin(e.currentTarget.value.trim())}
                    required
                  />
                </div>
                <div className={"form-control"}>
                  <label className={"label"}>
                    <span className={"label-text"}>Password</span>
                  </label>
                  <input
                    onChange={(e) => setPassword(e.currentTarget.value)}
                    type={"password"}
                    placeholder={"password"}
                    className={"input input-bordered"}
                    required
                  />
                </div>
                <FormError error={formError} />
                <div className={"form-control mt-6"}>
                  <button type={"submit"} className={"btn button m-auto"}>
                    Login
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </form>
    </>
  )
}

export default LoginPage
