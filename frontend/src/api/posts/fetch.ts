import { Comment, Post } from "../../custom"
import axios from "axios"

export const getPostsLocal = (): Promise<Post[]> =>
  axios(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts`).then((res) => res.data)

export const getPostLocal = (id: number) =>
  axios<Post | undefined>(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/${id}`).then(
    (res) => res.data
  )

export const getPostCommentsLocal = (id: number) =>
  axios<Comment[]>(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/${id}/comments`).then(
    (res) => res.data
  )

export const getUserPostsLocal = (user_id: number) =>
  axios<Post[]>(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/user/${user_id}/posts`).then(
    (res) => res.data
  )

export const getCategoriesLocal = () =>
  axios
    .get<string[]>(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/categories`)
    .then((res) => res.data)

export const getCategoryPostsLocal = (category: string) =>
  axios<Post[]>(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/categories/${category}`).then(
    (res) => res.data
  )
