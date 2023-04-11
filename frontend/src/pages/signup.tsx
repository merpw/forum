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
        <div className={"mb-6"}>
          <label
            htmlFor={"username"}
            className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
          >
            Username
            <input
              type={"text"}
              className={
                "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
              }
              onInput={(e) => setName(e.currentTarget.value)}
              placeholder={"Username"}
              required
            />
          </label>
        </div>
        <div className={"flex flex-raw mb-3 gap-3"}>
          <div className={"w-full"}>
            <label
              htmlFor={"first_name"}
              className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
            >
              First Name
              <input
                type={"text"}
                className={
                  "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                }
                onInput={(e) => setFirstName(e.currentTarget.value)}
                placeholder={"First Name"}
                required
              />
            </label>
          </div>
          <div className={"w-full"}>
            <label
              htmlFor={"last_name"}
              className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
            >
              Last Name
              <input
                type={"text"}
                className={
                  "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                }
                onInput={(e) => setLastName(e.currentTarget.value)}
                placeholder={"Last Name"}
                required
              />
            </label>
          </div>
        </div>
        <div className={"flex flex-raw mb-6 gap-3"}>
          <div className={"w-full flex-1 overflow-x-hidden"}>
            <label
              htmlFor={"dob"}
              className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
            >
              Date of Birth
              <input
                type={"date"}
                max={new Date().toISOString().split("T")[0]}
                className={
                  "h-10 bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                }
                onInput={(e) => setDoB(e.currentTarget.value)}
                placeholder={"Date of Birth"}
                required
              />
            </label>
          </div>
          <div className={"w-full flex-1"}>
            <label
              htmlFor={"gender"}
              className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
            >
              Gender
              <select
                className={
                  "h-10 bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                }
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
        <div className={"mb-6"}>
          <label
            htmlFor={"email"}
            className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
          >
            Your email
            <input
              type={"email"}
              className={
                "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
              }
              onInput={(e) => setEmail(e.currentTarget.value)}
              placeholder={"Email"}
              required
            />
          </label>
        </div>
        <div className={"mb-3"}>
          <label
            htmlFor={"password"}
            className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
          >
            Create password
            <input
              onInput={(e) => setPassword(e.currentTarget.value)}
              type={"password"}
              className={
                "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
              }
              required
            />
          </label>
        </div>
        <div className={"mb-6"}>
          <label
            htmlFor={"repeat-password"}
            className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
          >
            Repeat password
          </label>
          <input
            onInput={(e) => setPasswordConfirm(e.currentTarget.value)}
            type={"password"}
            id={"repeat-password"}
            className={
              "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            }
            required
          />
        </div>
        <FormError error={formError} />
        <button
          type={"submit"}
          className={
            "text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
          }
        >
          Submit
        </button>
      </form>
    </>
  )
}

export default SignupPage
