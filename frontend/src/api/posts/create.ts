import axios from "axios"

export const CreatePost = async (postData: {
  title: string
  content: string
  description: string
  categories: string[]
  group_id?: number
}) =>
  axios
    .post<number>("/api/posts/create", postData)
    .then((res) => res.data)
    .catch((err) => {
      throw new Error(err.response?.data?.length < 200 ? err.response.data : "Unexpected error")
    })

export const generateDescription = async ({ title, content }: { title: string; content: string }) =>
  axios
    .post<string>(
      "/api/next-public/ai",
      {
        action: "GENERATE_DESCRIPTION",
        body: {
          title,
          content,
        },
      },
      { withCredentials: true }
    )
    .then((res) => res.data)
    .catch((err) => {
      throw new Error(err.response?.data?.length < 200 ? err.response.data : "Unexpected error")
    })
