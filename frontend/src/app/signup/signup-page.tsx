"use client"

import { useRouter } from "next/navigation"
import { useEffect, useState } from "react"

import { logIn, SignUp, useMe } from "@/api/auth"
import { FormError } from "@/components/error"

const SignupPage = () => {
  const router = useRouter()

  const { isLoading, isLoggedIn, mutate } = useMe()

  const [formFields, setFormFields] = useState<{
    name: string
    email: string
    password: string
    passwordConfirm: string
    first_name: string
    last_name: string
    dob: string
    gender: string
  }>({
    name: "",
    email: "",
    password: "",
    passwordConfirm: "",
    first_name: "",
    last_name: "",
    dob: "",
    gender: "",
  })

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
          formFields.name = formFields.name.trim()
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
                    name={"name"}
                    minLength={3}
                    maxLength={15}
                    placeholder={"username"}
                    value={formFields.name}
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
                    <label className={"label "}>
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
                  </label>
                  <input
                    type={"password"}
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
                    type={"password"}
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
