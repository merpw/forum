import { Post, Comment } from "../../custom"
import axios from "axios"

export const getPostsLocal = (): Promise<Post[]> =>
  axios(`http://${process.env.FORUM_BACKEND_LOCALHOST}/api/posts`).then((res) => res.data)

export const getPostLocal = (id: number) =>
  axios<Post | undefined>(`http://${process.env.FORUM_BACKEND_LOCALHOST}/api/posts/${id}`)
    .then((res) => res.data)
    .catch(() => undefined)

export const getPostCommentsLocal = (id: number) =>
  axios<Comment[]>(`http://${process.env.FORUM_BACKEND_LOCALHOST}/api/posts/${id}/comments`)
    .then((res) => res.data)
    .catch(() => [])

export const getUserPostsLocal = (user_id: number) =>
  axios<Post[]>(`http://${process.env.FORUM_BACKEND_LOCALHOST}/api/user/${user_id}/posts`)
    .then((res) => {
      return { posts: res.data }
    })
    .catch(() => {
      return { posts: [] as Post[] }
    })

export const getCategoriesLocal = () =>
  axios
    .get<string[]>(`http://${process.env.FORUM_BACKEND_LOCALHOST}/api/posts/categories`)
    .then((res) => res.data)

export const getCategoryPostsLocal = (category: string) =>
  axios<Post[]>(`http://${process.env.FORUM_BACKEND_LOCALHOST}/api/posts/categories/${category}`)
    .then((res) => {
      return { posts: res.data }
    })
    .catch(() => {
      return { posts: undefined }
    })
