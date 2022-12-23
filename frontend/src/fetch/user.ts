import useSWR from "swr"
import { User } from "../custom"
import { user } from "./dummy"

export const useUser = () => {
  const { data, error } = useSWR("/api/me", getUser)

  return {
    isError: error != undefined,
    isLoading: !error && !data,
    isLoggedIn: !error && data?.user != undefined,
    user: data?.user,
  }
}

const getUser = async (): Promise<{ user: User | undefined }> => {
  return Promise.resolve({ user: user })
}
