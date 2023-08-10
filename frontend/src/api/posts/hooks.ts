import useSWR from "swr"

import fetcher from "@/api/fetcher"
import { Post } from "@/custom"

export const usePost = (id: string) => {
  const { data, error } = useSWR<Post>(`/api/posts/${id}`, fetcher)

  return {
    post: data,
    isLoading: !error && !data,
    error,
  }
}

export const usePostList = (postIds: number[]) => {
  const { data, error } = useSWR<Post[]>(
    postIds.map((id) => `/api/posts/${id}`),
    (...urls: string[]) => Promise.all(urls.map((url) => fetcher(url)))
  )

  return {
    posts: data,
    isLoading: !error && !data,
    error,
  }
}
