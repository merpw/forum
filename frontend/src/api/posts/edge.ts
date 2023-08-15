import { Comment, Post } from "@/custom"

export const getPostsLocal = (): Promise<Post[]> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts`, {
    headers: {
      "Internal-Auth": process.env.FORUM_BACKEND_SECRET || "secret",
    },
  }).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })

export const getPostLocal = (id: number): Promise<Post> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/${id}`, {
    headers: {
      "Internal-Auth": process.env.FORUM_BACKEND_SECRET || "secret",
    },
  }).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })

export const getPostCommentsLocal = (id: number): Promise<Comment[]> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/${id}/comments`, {
    headers: {
      "Internal-Auth": process.env.FORUM_BACKEND_SECRET || "secret",
    },
  }).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })

export const getCategoriesLocal = (): Promise<string[]> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/categories`, {
    headers: {
      "Internal-Auth": process.env.FORUM_BACKEND_SECRET || "secret",
    },
  }).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })
