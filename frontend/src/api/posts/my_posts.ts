import axios from "axios"
import useSWR from "swr"

import { Post } from "@/custom"

const getMyPosts = async () =>
  document.cookie.includes("forum-token")
    ? axios
        .get<Post[]>(`/api/me/posts`, { withCredentials: true })
        .then((res) => {
          return { posts: res.data }
        })
        .catch(() => {
          return { posts: undefined }
        })
    : { posts: undefined }

export const useMyPosts = () => {
  const { data, error } = useSWR("/api/me/posts", getMyPosts)

  return {
    isError: error != undefined,
    isLoading: !error && !data,
    posts: data?.posts,
  }
}
