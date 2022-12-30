import axios from "axios"

export const CreateComment = (post_id: number, text: string) =>
  axios
    .post<number>(`/api/posts/${post_id}/comment`, { content: text }, { withCredentials: true })
    .then((res) => res.data)
