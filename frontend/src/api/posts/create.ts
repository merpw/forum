import axios from "axios"

export const CreatePost = (title: string, content: string, categories: string[]) =>
  axios
    .post<number>("/api/posts/create", { title, content, categories }, { withCredentials: true })
    .then((res) => res.data)
