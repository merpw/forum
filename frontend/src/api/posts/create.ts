import axios from "axios"

export const CreatePost = (title: string, content: string, category: string[]) =>
  axios
    .post<number>("/api/posts/create", { title, content, category }, { withCredentials: true })
    .then((res) => res.data)
