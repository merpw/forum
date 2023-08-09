import useSWR from "swr"

import { Post } from "@/custom"
import fetcher from "@/api/fetcher"

export const useGroupPosts = (groupId: number) => {
  const { data, error, mutate } = useSWR<Post[] | undefined>(`/api/posts`, fetcher)

  return {
    posts: data,
    isLoading: !error && !data,
    error,
    mutate,
  }
}
