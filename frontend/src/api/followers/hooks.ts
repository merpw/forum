import useSWR from "swr"

import fetcher from "@/api/fetcher"

export const useFollowers = () => {
  const { data, error, mutate } = useSWR<number[]>("/api/me/followers", fetcher)

  const loading = !data && !error

  return {
    followers: data,
    loading,
    error,
    mutate,
  }
}

export const useFollowed = () => {
  const { data, error, mutate } = useSWR<number[]>("/api/me/followed", fetcher)

  const loading = !data && !error

  return {
    followed: data,
    loading,
    error,
    mutate,
  }
}
