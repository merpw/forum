import axios from "axios"
import useSWR from "swr"
import { Post } from "@/custom"

const getMyPostsLiked = async () =>
  document.cookie.includes("forum-token")
    ? axios
        .get<Post[]>(`/api/me/posts/liked`, { withCredentials: true })
        .then((res) => {
          return { posts: res.data }
        })
        .catch(() => {
          return { posts: undefined }
        })
    : { posts: undefined }

export const useMyPostsLiked = () => {
  const { data, error } = useSWR("/api/me/posts/liked", getMyPostsLiked)

  return {
    isError: error != undefined,
    isLoading: !error && !data,
    posts: data?.posts,
  }
}
