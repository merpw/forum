import useSWR from "swr"
import { User } from "../custom"
import { user as dummy_user } from "./dummy"

export const useUser = () => {
  const { data, mutate, error } = useSWR<{ user: User | undefined }>("/api/me", getUser)

  return {
    isError: error != undefined,
    isLoading: !error && !data,
    isLoggedIn: !error && data?.user != undefined,
    mutate: mutate,
    user: data?.user,
  }
}

const getUser = async (): Promise<{ user: User | undefined }> => {
  if (document.cookie.includes("forum-test-token=admin")) {
    return Promise.resolve({ user: dummy_user })
  }
  return Promise.resolve({ user: undefined })
}

export const logIn = async (login: string, password: string): Promise<void> => {
  if (login == "admin" && password == "admin") {
    document.cookie = "forum-test-token=admin;path=/;"
    return Promise.resolve()
  }
  return Promise.reject("Invalid login or password")
}

export const logOut = async (): Promise<void> => {
  document.cookie = "forum-test-token=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/;"
  return Promise.resolve()
}

export const SignUp = async (name: string, email: string, password: string): Promise<void> => {
  return Promise.reject("Not connected to backend yet")
}
