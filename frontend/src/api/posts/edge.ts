import { Post } from "@/custom"

export const getPostsLocal = (): Promise<Post[]> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts`).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })

export const getCategoriesLocal = (): Promise<string[]> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/categories`).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })
