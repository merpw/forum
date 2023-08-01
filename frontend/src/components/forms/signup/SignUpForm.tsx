"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"

import { logIn, SignUp } from "@/api/auth/hooks"
import { FormError } from "@/components/error"
import ChooseAvatar from "@/components/forms/signup/ChooseAvatar"
import Passwords from "@/components/forms/signup/Passwords"
import FullName from "@/components/forms/signup/FullName"
import Username from "@/components/forms/signup/Username"
import Email from "@/components/forms/signup/Email"
import DateOfBirth from "@/components/forms/signup/DateOfBirth"

const SignUpForm = () => {
  const router = useRouter()
  const [formFields, setFormFields] = useState<{
    username: string
    email: string
    password: string
    passwordConfirm: string
    first_name: string
    last_name: string
    dob: string
    gender: string

    avatar?: string
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
  const [isSame, setIsSame] = useState(false)

  const [formError, setFormError] = useState<string | null>(null)

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
            .then(() =>
              logIn(formFields.email, formFields.password).then(() => {
                router.refresh()
              })
            )
            .catch((err) => setFormError(err.message))
        }}
        className={"flex justify-center min-h-full bg-base-200"}
      >
        <div className={"hero-content flex-col"}>
          <div className={"text-center"}>
            <h1 className={"text-6xl gradient-text font-Yesteryear"}>Welcome</h1>
            <div className={"opacity-50 font-light"}>
              to the place <span className={"font-Alatsi text-primary text-xl"}>FOR</span>{" "}
              <span className={"font-Alatsi text-primary text-xl"}>U</span>nleashing{" "}
              <span className={"font-Alatsi text-primary text-xl"}> M</span>inds!
            </div>
          </div>

          <div className={"card flex-shrink-0 w-full max-w-2xl shadow-xl bg-base-100"}>
            <div className={"card-body"}>
              <ChooseAvatar
                avatar={formFields.avatar}
                setAvatar={(newAvatar) => setFormFields({ ...formFields, avatar: newAvatar })}
              />

              <Username {...formFields} />

              <FullName first_name={formFields.first_name} last_name={formFields.last_name} />

              <DateOfBirth dob={formFields.dob} />

              <Email email={formFields.email} />

              <Passwords />

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
      </form>
    </>
  )
}

export default SignUpForm
