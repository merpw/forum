import axios from "axios"
import useSWR from "swr"
import { User } from "../custom"

export const useMe = () => {
  const { data, mutate, error } = useSWR<{ user: User | undefined }>("/api/me", getMe, {
    revalidateOnFocus: false,
  })
  return {
    isError: error != undefined,
    isLoading: !error && !data,
    isLoggedIn: !error && data?.user != undefined,
    mutate: mutate,
    user: data?.user,
  }
}

const getMe = async () =>
  document.cookie.includes("forum-token")
    ? axios
        .get<User>("/api/me", { withCredentials: true })
        .then((res) => {
          return { user: res.data }
        })
        .catch(() => {
          return { user: undefined }
        })
    : { user: undefined }

export const logIn = async (login: string, password: string) =>
  axios.post("/api/login", { login, password }, { withCredentials: true })

export const logOut = async (): Promise<void> =>
  axios.post("/api/logout", {}, { withCredentials: true })

export const SignUp = async (
  name: string,
  email: string,
  password: string,
  first_name: string,
  last_name: string,
  age: string,
  gender: string
) =>
  axios.post(
    "/api/signup",
    { name, email, password, first_name, last_name, age, gender },
    { withCredentials: true }
  )
