import axios from "axios"

export const CreatePost = (postData: {
  title: string
  content: string
  description: string
  categories: string[]
}) =>
  axios
    .post<number>("/api/posts/create", postData, { withCredentials: true })
    .then((res) => res.data)
