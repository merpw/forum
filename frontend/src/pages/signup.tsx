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

  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [passwordConfirm, setPasswordConfirm] = useState("")
  const [first_name, setFirstName] = useState("")
  const [last_name, setLastName] = useState("")
  const [dob, setDoB] = useState("")
  const [gender, setGender] = useState("")

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
  }, [name, email, password, passwordConfirm, first_name, last_name, dob, gender])

  return (
    <>
      <Head>
        <title>Sign Up - Forum</title>
        <meta name={"title"} content={"Sign Up - Forum"} />
        <meta name={"og:title"} content={"Sign Up - Forum"} />
      </Head>
      <form
        onSubmit={async (e) => {
          e.preventDefault()

          if (isSame) return

          if (formError != null) setFormError(null)
          setIsSame(true)

          if (password != passwordConfirm) {
            setFormError("Passwords don't match")
            return
          }

          SignUp(name, email, password, first_name, last_name, dob, gender)
            .then((response) => {
              if (response.status == 200) {
                logIn(email, password).then(() => mutate())
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
        <label htmlFor={"username"} className={"field label"}>
          <p className={"field-title"}>Username</p>
          <input
            type={"text"}
            className={"input"}
            onInput={(e) => setName(e.currentTarget.value)}
            placeholder={"Username"}
            required
          />
        </label>

        <div className={"field flex flex-raw gap-3"}>
          <div className={"w-full"}>
            <label htmlFor={"first_name"} className={"label"}>
              <p className={"field-title"}>First Name </p>
              <input
                type={"text"}
                className={"input"}
                onInput={(e) => setFirstName(e.currentTarget.value)}
                placeholder={"First Name"}
                required
              />
            </label>
          </div>
          <div className={"w-full"}>
            <label htmlFor={"last_name"} className={"label"}>
              <p className={"field-title"}>Last Name </p>
              <input
                type={"text"}
                className={"input"}
                onInput={(e) => setLastName(e.currentTarget.value)}
                placeholder={"Last Name"}
                required
              />
            </label>
          </div>
        </div>
        <div className={"field flex flex-raw gap-3"}>
          <div className={"w-full flex-1 overflow-x-hidden"}>
            <label htmlFor={"dob"} className={"label"}>
              <p className={"field-title"}>Date of Birth </p>
              <input
                type={"date"}
                max={new Date().toISOString().split("T")[0]}
                className={"input"}
                onInput={(e) => setDoB(e.currentTarget.value)}
                placeholder={"Date of Birth"}
                required
              />
            </label>
          </div>
          <div className={"w-full flex-1"}>
            <label htmlFor={"gender"} className={"label"}>
              <p className={"field-title"}>Gender </p>
              <select
                className={"input"}
                onInput={(e) => setGender(e.currentTarget.value)}
                placeholder={"Gender"}
                required
              >
                <option value={""}>Select</option>
                <option value={"male"}>Male</option>
                <option value={"female"}>Female</option>
                <option value={"other"}>Other</option>
              </select>
            </label>
          </div>
        </div>

        <label htmlFor={"email"} className={"field label"}>
          <p className={"field-title"}>Your email</p>
          <input
            type={"email"}
            className={"input"}
            onInput={(e) => setEmail(e.currentTarget.value)}
            placeholder={"Email"}
            required
          />
        </label>

        <label htmlFor={"password"} className={"field label"}>
          <p className={"field-title"}>Create password </p>
          <input
            onInput={(e) => setPassword(e.currentTarget.value)}
            type={"password"}
            className={"input"}
            required
          />
        </label>

        <label htmlFor={"repeat-password"} className={"field label"}>
          <p className={"field-title"}>Repeat password </p>
          <input
            onInput={(e) => setPasswordConfirm(e.currentTarget.value)}
            type={"password"}
            id={"repeat-password"}
            className={"input"}
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
