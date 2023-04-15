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
      <Head>
        <title>Sign Up - Forum</title>
        <meta name={"title"} content={"Sign Up - Forum"} />
        <meta name={"og:title"} content={"Sign Up - Forum"} />
      </Head>
      <form
        onChange={(e) => {
          const target = e.target as HTMLInputElement
          setFormFields({ ...formFields, [target.name]: target.value })
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

          SignUp(
            formFields.name,
            formFields.email,
            formFields.password,
            formFields.first_name,
            formFields.last_name,
            formFields.dob,
            formFields.gender
          )
            .then((response) => {
              if (response.status == 200) {
                logIn(formFields.email, formFields.password).then(() => mutate())
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
        <label htmlFor={"username"} className={"mb-6 label"}>
          <p className={"inputbox-title"}>Username</p>
          <input
            type={"text"}
            className={"inputbox"}
            name={"name"}
            placeholder={"Username"}
            required
          />
        </label>

        <div className={"mb-6 flex flex-raw gap-3"}>
          <div className={"w-full"}>
            <label htmlFor={"first_name"} className={"label"}>
              <p className={"inputbox-title"}>First Name </p>
              <input
                type={"text"}
                className={"inputbox"}
                name={"first_name"}
                placeholder={"First Name"}
                required
              />
            </label>
          </div>
          <div className={"w-full"}>
            <label htmlFor={"last_name"} className={"label"}>
              <p className={"inputbox-title"}>Last Name </p>
              <input
                type={"text"}
                className={"inputbox"}
                name={"last_name"}
                placeholder={"Last Name"}
                required
              />
            </label>
          </div>
        </div>
        <div className={"mb-6 flex flex-raw gap-3"}>
          <div className={"w-full flex-1 overflow-x-hidden"}>
            <label htmlFor={"dob"} className={"label"}>
              <p className={"inputbox-title"}>Date of Birth </p>
              <input
                type={"date"}
                max={new Date().toISOString().split("T")[0]}
                className={"inputbox"}
                name={"dob"}
                placeholder={"Date of Birth"}
                required
              />
            </label>
          </div>
          <div className={"w-full flex-1"}>
            <label htmlFor={"gender"} className={"label"}>
              <p className={"inputbox-title"}>Gender </p>
              <select className={"inputbox"} name={"gender"} placeholder={"Gender"} required>
                <option value={""}>Select</option>
                <option value={"male"}>Male</option>
                <option value={"female"}>Female</option>
                <option value={"other"}>Other</option>
              </select>
            </label>
          </div>
        </div>

        <label htmlFor={"email"} className={"mb-6 label"}>
          <p className={"inputbox-title"}>Your email</p>
          <input
            type={"email"}
            className={"inputbox"}
            name={"email"}
            placeholder={"Email"}
            required
          />
        </label>

        <label htmlFor={"password"} className={"mb-6 label"}>
          <p className={"inputbox-title"}>Create password </p>
          <input type={"password"} className={"inputbox"} name={"password"} required />
        </label>

        <label htmlFor={"repeat-password"} className={"mb-6 label"}>
          <p className={"inputbox-title"}>Repeat password </p>
          <input
            type={"password"}
            id={"repeat-password"}
            className={"inputbox"}
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
