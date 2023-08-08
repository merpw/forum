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

export const useFollowing = () => {
  const { data, error, mutate } = useSWR<number[]>("/api/me/following", fetcher)

  const loading = !data && !error

  return {
    following: data,
    loading,
    error,
    mutate,
  }
}
