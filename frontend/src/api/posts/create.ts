import axios from "axios"

export const CreatePost = (title: string, content: string) =>
  axios.post<number>("/api/posts/create", { title, content }, { withCredentials: true })
