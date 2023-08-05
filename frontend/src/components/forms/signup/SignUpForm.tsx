"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"

import { logIn, SignUp } from "@/api/auth/hooks"
import { FormError } from "@/components/error"
import ChooseAvatar from "@/components/forms/signup/ChooseAvatar"
import Passwords from "@/components/forms/signup/Passwords"
import FullName from "@/components/forms/signup/FullName"
import Username from "@/components/forms/signup/Username"
import Email from "@/components/forms/signup/Email"
import DateOfBirth from "@/components/forms/signup/DateOfBirth"
import Bio from "@/components/forms/signup/Bio"

const SignUpForm = () => {
  const router = useRouter()

  const [isSame, setIsSame] = useState(false)

  const [formError, setFormError] = useState<string | null>(null)

  return (
    <>
      <form
        onChange={() => {
          if (isSame) setIsSame(false)
        }}
        onSubmit={async (e) => {
          e.preventDefault()

          if (isSame) return

          const formData = new FormData(e.target as HTMLFormElement)

          const formFields = {
            username: formData.get("username") as string,
            email: formData.get("email") as string,
            password: formData.get("password") as string,
            passwordConfirm: formData.get("passwordConfirm") as string,
            first_name: formData.get("first_name") as string,
            last_name: formData.get("last_name") as string,
            dob: formData.get("dob") as string,
            gender: formData.get("gender") as string,

            avatar: formData.get("avatar") || undefined,
            bio: formData.get("bio") || undefined,
          }

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
              <ChooseAvatar />

              <Username />

              <FullName />

              <DateOfBirth />

              <Bio />

              <Email />

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
