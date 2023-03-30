import { Comment, Post } from "../../custom"
import axios from "axios"

export const getPostsLocal = (): Promise<Post[]> =>
  process.env.FORUM_BACKEND_PRIVATE_URL
    ? axios(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts`).then((res) => res.data)
    : Promise.resolve([])

export const getPostLocal = (id: number) =>
  process.env.FORUM_BACKEND_PRIVATE_URL
    ? axios<Post | undefined>(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/${id}`).then(
        (res) => res.data
      )
    : Promise.resolve(undefined)

export const getPostCommentsLocal = (id: number) =>
  process.env.FORUM_BACKEND_PRIVATE_URL
    ? axios<Comment[]>(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/${id}/comments`).then(
        (res) => res.data
      )
    : Promise.resolve([])

export const getUserPostsLocal = (user_id: number) =>
  process.env.FORUM_BACKEND_PRIVATE_URL
    ? axios<Post[]>(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/user/${user_id}/posts`).then(
        (res) => res.data
      )
    : Promise.resolve([])

export const getCategoriesLocal = () =>
  process.env.FORUM_BACKEND_PRIVATE_URL
    ? axios
        .get<string[]>(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/categories`)
        .then((res) => res.data)
    : Promise.resolve([])

export const getCategoryPostsLocal = (category: string) =>
  process.env.FORUM_BACKEND_PRIVATE_URL
    ? axios<Post[]>(
        `${process.env.FORUM_BACKEND_PRIVATE_URL}/api/posts/categories/${category}`
      ).then((res) => res.data)
    : Promise.resolve([])
