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
        <label className={"mb-4 block"}>
          <p className={"inputbox-title"}>Username</p>
          <input
            type={"text"}
            className={"inputbox-singlerow"}
            name={"name"}
            minLength={3}
            maxLength={15}
            placeholder={"Username"}
            value={formFields.name}
            onChange={() => void 0 /* handled by Form */}
            required
          />
        </label>

        <div className={"mb-4 flex flex-wrap gap-3"}>
          <label className={"grow basis-1/3"}>
            <p className={"inputbox-title"}>First Name </p>
            <input
              type={"text"}
              className={"inputbox-singlerow"}
              name={"first_name"}
              placeholder={"First Name"}
              value={formFields.first_name}
              onChange={() => void 0 /* handled by Form */}
              maxLength={15}
              required
            />
          </label>
          <label className={"grow basis-1/3"}>
            <p className={"inputbox-title"}>Last Name</p>
            <input
              type={"text"}
              className={"inputbox-singlerow"}
              name={"last_name"}
              placeholder={"Last Name"}
              value={formFields.last_name}
              onChange={() => void 0 /* handled by Form */}
              maxLength={15}
              required
            />
          </label>
        </div>
        <div className={"mb-4 flex flex-wrap gap-3"}>
          <label className={"grow basis-1/3"}>
            <p className={"inputbox-title"}>Date of Birth </p>
            <input
              type={"date"}
              min={"1900-01-01"}
              max={new Date().toISOString().split("T")[0]}
              className={"inputbox-singlerow"}
              name={"dob"}
              placeholder={"Date of Birth"}
              required
            />
          </label>
          <label className={"grow basis-1/3 block"}>
            <span className={"inputbox-title m-0"}>Gender </span>
            <div>
              <select
                className={"inputbox-singlerow my-1 py-1"}
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
          </label>
        </div>

        <label className={"mb-4 block"}>
          <p className={"inputbox-title"}>Your email</p>
          <input
            type={"email"}
            className={"inputbox-singlerow"}
            name={"email"}
            placeholder={"Email"}
            value={formFields.email}
            onChange={() => void 0 /* handled by Form */}
            required
          />
        </label>

        <label className={"mb-4 block"}>
          <p className={"inputbox-title"}>Create password </p>
          <input
            type={"password"}
            className={"inputbox-singlerow"}
            name={"password"}
            required
            minLength={8}
          />
        </label>

        <label className={"mb-6 block"}>
          <p className={"inputbox-title"}>Repeat password </p>
          <input
            type={"password"}
            id={"repeat-password"}
            className={"inputbox-singlerow"}
            name={"passwordConfirm"}
            required
          />
        </label>

        <FormError error={formError} />
        <button type={"submit"} className={"button"}>
          Submit
        </button>
      </form>
    </>
  )
}

export default SignupPage
