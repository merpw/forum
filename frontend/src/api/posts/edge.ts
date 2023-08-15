import { Comment, Post } from "@/custom"
import edgeFetcher from "@/api/edge-fetcher"

export const getPostsLocal = () => edgeFetcher<Post[]>("/api/posts")

export const getPostLocal = (id: number) => edgeFetcher<Post>(`/api/posts/${id}`)

export const getPostCommentsLocal = (id: number) =>
  edgeFetcher<Comment[]>(`/api/posts/${id}/comments`)

export const getCategoriesLocal = () => edgeFetcher<string[]>("/api/posts/categories")
