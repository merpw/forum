import { Comment, Post } from "@/custom"

export const getPostsLocal = (): Promise<Post[]> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts`).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })

export const getPostLocal = (id: number): Promise<Post> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/${id}`).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })

export const getPostCommentsLocal = (id: number): Promise<Comment[]> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/${id}/comments`).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })

export const getCategoriesLocal = (): Promise<string[]> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/categories`).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })

export const getUserPostsLocal = (user_id: number): Promise<Post[]> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/user/${user_id}/posts`).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })
