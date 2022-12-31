import axios from "axios"

export const getCategories = () =>
  axios.get<string[]>("/api/posts/categories").then((res) => res.data)
