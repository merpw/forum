import axios from "axios"

export const CreatePost = async (postData: {
  title: string
  content: string
  description: string
  categories: string[]
}) => {
  if (postData.title.length === 0) {
    throw new Error("Title is too short")
  }
  if (postData.title.length > 25) {
    throw new Error("Title is too long")
  }
  if (postData.content.length === 0) {
    throw new Error("Content is too short")
  }
  if (postData.content.length > 10000) {
    throw new Error(`Content is too long, ${postData.content.length}/10000`)
  }
  if (postData.description.length === 0) {
    throw new Error("Description is too short")
  }
  if (postData.description.length > 205) {
    throw new Error(`Description is too long, ${postData.description.length}/200`)
  }
  if (postData.categories.length === 0) {
    throw new Error("Categories are not selected")
  }

  return axios
    .post<number>("/api/posts/create", postData, { withCredentials: true })
    .then((res) => res.data)
    .catch((err) => {
      throw new Error(err.response?.data?.length < 200 ? err.response.data : "Unexpected error")
    })
}

export const generateDescription = async ({
  title,
  content,
}: {
  title: string
  content: string
}) => {
  if (title.length === 0) {
    throw new Error("Title is too short")
  }
  if (title.length > 50) {
    throw new Error("Title is too long")
  }
  if (content.length === 0) {
    throw new Error("Content is too short")
  }
  return axios
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
}

// TODO: implement all of the requests with the same pattern (move the logic from components and pages to api). Add unexpected error handling.
