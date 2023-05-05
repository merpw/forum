import axios from "axios"
import useSWR from "swr"

import { User } from "@/custom"

export const useMe = () => {
  const { data, mutate, error } = useSWR<{ user: User | null }>("/api/me", getMe, {
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
          return { user: null }
        })
    : { user: null }

export const logIn = async (login: string, password: string) =>
  axios.post("/api/login", { login, password }, { withCredentials: true }).catch((err) => {
    throw new Error(err.response?.data?.length < 200 ? err.response.data : "Unexpected error")
  })

export const logOut = async (): Promise<void> =>
  axios.post("/api/logout", {}, { withCredentials: true })

export const SignUp = async (data: {
  name: string
  email: string
  password: string
  first_name: string
  last_name: string
  dob: string
  gender: string
}) =>
  axios.post("/api/signup", data, { withCredentials: true }).catch((err) => {
    throw new Error(err.response?.data?.length < 200 ? err.response.data : "Unexpected error")
  })
