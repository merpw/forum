"use client"

import { useRouter } from "next/navigation"
import { useEffect, useState } from "react"

import { logIn, SignUp, useMe } from "@/api/auth/hooks"
import { FormError } from "@/components/error"

const SignupPage = () => {
  const router = useRouter()

  const { isLoading, isLoggedIn, mutate } = useMe()

  const [formFields, setFormFields] = useState<{
    username: string
    email: string
    password: string
    passwordConfirm: string
    first_name: string
    last_name: string
    dob: string
    gender: string
  }>({
    username: "",
    email: "",
    password: "",
    passwordConfirm: "",
    first_name: "",
    last_name: "",
    dob: "",
    gender: "",
  })
  const [showPassword, setShowPassword] = useState(false)

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
  }, [formFields])

  return (
    <>
      <form
        onChange={(e) => {
          const target = e.target as HTMLInputElement
          setFormFields({ ...formFields, [target.name]: target.value })
        }}
        onBlur={() => {
          formFields.username = formFields.username.trim()
          formFields.first_name = formFields.first_name.trim()
          formFields.last_name = formFields.last_name.trim()
          formFields.email = formFields.email.trim()

          setFormFields({ ...formFields })
        }}
        onSubmit={async (e) => {
          e.preventDefault()

          if (isSame) return

          if (formError != null) setFormError(null)
          setIsSame(true)

          if (formFields.password != formFields.passwordConfirm) {
            setFormError("Passwords don't match")
            return
          }

          SignUp(formFields)
            .then(() => logIn(formFields.email, formFields.password).then(() => mutate()))
            .catch((err) => setFormError(err.message))
        }}
      >
        <div className={"py-1 bg-base-200"}>
          <div className={"hero-content flex-col py-7"}>
            <div className={"text-center"}>
              <h1 className={"text-6xl gradient-text font-Yesteryear"}>Welcome</h1>
              <div className={"opacity-50 font-light"}>
                to the place <span className={"font-Alatsi text-primary text-xl"}>FOR</span>{" "}
                <span className={"font-Alatsi text-primary text-xl"}>U</span>nleashing{" "}
                <span className={"font-Alatsi text-primary text-xl"}> M</span>inds!
              </div>
            </div>

            <div className={"card flex-shrink-0 w-full max-w-2xl shadow-xl bg-base-100 mb-10"}>
              <div className={"card-body"}>
                <div className={"form-control"}>
                  <label className={"label pt-0"}>
                    <span className={"label-text"}>Username</span>
                  </label>
                  <input
                    type={"text"}
                    className={"input input-bordered"}
                    name={"username"}
                    minLength={3}
                    maxLength={15}
                    placeholder={"username"}
                    value={formFields.username}
                    onChange={() => void 0 /* handled by Form */}
                    required
                  />
                </div>

                <div className={"flex flex-wrap flex-row gap-3"}>
                  <div className={"grow basis-1/3"}>
                    <label className={"label"}>
                      <span className={"label-text"}>First Name </span>
                    </label>
                    <input
                      type={"text"}
                      className={"input input-bordered w-full"}
                      name={"first_name"}
                      placeholder={"first name"}
                      value={formFields.first_name}
                      onChange={() => void 0 /* handled by Form */}
                      maxLength={15}
                      required
                    />
                  </div>
                  <div className={"grow basis-1/3"}>
                    <label className={"label"}>
                      <span className={"label-text"}>Last Name </span>
                    </label>
                    <input
                      type={"text"}
                      className={"input input-bordered w-full"}
                      name={"last_name"}
                      placeholder={"last name"}
                      value={formFields.last_name}
                      onChange={() => void 0 /* handled by Form */}
                      maxLength={15}
                      required
                    />
                  </div>
                </div>

                <div className={"flex flex-wrap flex-row gap-3"}>
                  <div className={"grow basis-1/3"}>
                    <label className={"label"}>
                      <span className={"label-text"}>Date of Birth</span>
                    </label>
                    <input
                      type={"date"}
                      min={"1900-01-01"}
                      max={new Date().toISOString().split("T")[0]}
                      className={"input input-bordered w-full text-sm"}
                      name={"dob"}
                      placeholder={"Date of Birth"}
                      required
                    />
                  </div>
                  <div className={"grow basis-1/3"}>
                    <label className={"label"}>
                      <span className={"label-text"}>Gender</span>
                    </label>
                    <div>
                      <select
                        title={"Your gender"}
                        className={"input input-bordered w-full text-sm"}
                        name={"gender"}
                        placeholder={"Gender"}
                        required
                      >
                        <option value={""}>Select</option>
                        <option value={"male"}>Male</option>
                        <option value={"female"}>Female</option>
                        <option value={"other"}>Other</option>
                      </select>
                    </div>
                  </div>
                </div>

                <div className={"form-control"}>
                  <label className={"label"}>
                    <span className={"label-text"}>Email</span>
                  </label>
                  <input
                    type={"email"}
                    className={"input input-bordered"}
                    name={"email"}
                    placeholder={"email"}
                    value={formFields.email}
                    onChange={() => void 0 /* handled by Form */}
                    required
                  />
                </div>

                <div className={"form-control"}>
                  <label className={"label"}>
                    <span className={"label-text"}>Create password </span>
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
                    type={showPassword ? "text" : "password"}
                    className={"input input-bordered"}
                    name={"password"}
                    placeholder={"password"}
                    required
                    minLength={8}
                  />
                </div>

                <div className={"form-control"}>
                  <label className={"label"}>
                    <span className={"label-text"}>Repeat password </span>
                  </label>
                  <input
                    title={"Repeat password"}
                    type={showPassword ? "text" : "password"}
                    id={"repeat-password"}
                    className={"input input-bordered"}
                    name={"passwordConfirm"}
                    placeholder={"repeat password"}
                    required
                  />
                </div>
                <FormError error={formError} />
                <div className={"form-control inline mt-5 mx-auto"}>
                  <span className={"font-black mr-3 text-neutral"}>• •</span>
                  <button type={"submit"} className={"button"}>
                    Signup
                  </button>
                  <span className={"font-black ml-3 text-neutral"}>• •</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </form>
    </>
  )
}

export default SignupPage
