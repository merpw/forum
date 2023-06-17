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
  const [showPassword, setShowPassword] = useState(false)
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
              <div className={"font-light"}>
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
                    <button
                      type={"button"}
                      onClick={() => setShowPassword(!showPassword)}
                      className={"label-text text-info clickable"}
                    >
                      {showPassword ? (
                        <svg
                          xmlns={"http://www.w3.org/2000/svg"}
                          fill={"none"}
                          viewBox={"0 0 24 24"}
                          strokeWidth={1}
                          stroke={"currentColor"}
                          className={"w-5 h-5"}
                        >
                          <path
                            strokeLinecap={"round"}
                            strokeLinejoin={"round"}
                            d={
                              "M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88"
                            }
                          />
                        </svg>
                      ) : (
                        <svg
                          xmlns={"http://www.w3.org/2000/svg"}
                          fill={"none"}
                          viewBox={"0 0 24 24"}
                          strokeWidth={1}
                          stroke={"currentColor"}
                          className={"w-5 h-5"}
                        >
                          <path
                            strokeLinecap={"round"}
                            strokeLinejoin={"round"}
                            d={
                              "M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z"
                            }
                          />
                          <path
                            strokeLinecap={"round"}
                            strokeLinejoin={"round"}
                            d={"M15 12a3 3 0 11-6 0 3 3 0 016 0z"}
                          />
                        </svg>
                      )}
                    </button>
                  </label>
                  <input
                    onChange={(e) => setPassword(e.currentTarget.value)}
                    type={showPassword ? "text" : "password"}
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
