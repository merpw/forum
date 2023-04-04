import axios from "axios"
import useSWR from "swr"
import { Comment } from "@/custom"
export const CreateComment = (post_id: number, text: string) =>
  axios
    .post<number>(`/api/posts/${post_id}/comment`, { content: text }, { withCredentials: true })
    .then((res) => res.data)

const getComments = (path: string) => axios.get<Comment[]>(path).then((res) => res.data)

export const useComments = (post_id: number) => {
  const { data, error, mutate } = useSWR<Comment[]>(`/api/posts/${post_id}/comments`, getComments)
  return {
    comments: data,
    isLoading: !error && !data,
    isError: error,
    mutate,
  }
}
