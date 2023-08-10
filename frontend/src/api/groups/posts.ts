import useSWR from "swr"

import fetcher from "@/api/fetcher"

export const useGroupPosts = (groupId: number) => {
  const { data, error, mutate } = useSWR<number[] | undefined>(
    `/api/groups/${groupId}/posts`,
    fetcher
  )

  return {
    posts: data,
    isLoading: !error && !data,
    error,
    mutate,
  }
}
