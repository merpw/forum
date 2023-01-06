import axios from "axios"
import useSWR from "swr"

export const useReactions = (postID: number) => {
  const { data, mutate, error } = useSWR(`/api/posts/${postID}/reaction`, getPostReaction)
  return {
    isError: error != undefined,
    isLoading: !error && !data,
    mutate: mutate,
    reaction: data?.reaction,
    likes_count: data?.likes_count,
    dislikes_count: data?.dislikes_count,
  }
}

export const useCommentReactions = (postID: number, commentID: number | undefined) => {
  const { data, mutate, error } = useSWR(
    `/api/posts/${postID}/comment/${commentID}/reaction`,
    getCommentReaction
  )
  return {
    isError: error != undefined,
    isLoading: !error && !data,
    mutate,
    reaction: data?.reaction,
    likes_count: data?.likes_count,
    dislikes_count: data?.dislikes_count,
  }
}

type ReactionResponse = {
  reaction: number
  likes_count: number
  dislikes_count: number | undefined
}

export const getPostReaction = (path: string) =>
  document.cookie.includes("forum-token")
    ? axios
        .get<ReactionResponse>(path, {
          withCredentials: true,
        })
        .then((res) => res.data)
        .catch(() => undefined)
    : undefined

export const getCommentReaction = (path: string) =>
  document.cookie.includes("forum-token")
    ? axios
        .get<ReactionResponse>(path, {
          withCredentials: true,
        })
        .then((res) => res.data)
        .catch(() => undefined)
    : undefined

export const dislikePost = (postID: number) =>
  document.cookie.includes("forum-token")
    ? axios
        .post<number>(`/api/posts/${postID}/dislike`, null, { withCredentials: true })
        .then((res) => res.data)
    : Promise.resolve(0)

export const likePost = (postID: number) =>
  document.cookie.includes("forum-token")
    ? axios
        .post<number>(`/api/posts/${postID}/like`, null, { withCredentials: true })
        .then((res) => res.data)
    : Promise.resolve(0)

export const likeComment = (postID: number, commentID: number) =>
  document.cookie.includes("forum-token")
    ? axios
        .post<number>(`/api/posts/${postID}/comment/${commentID}/like`, null, {
          withCredentials: true,
        })
        .then((res) => res.data)
    : Promise.resolve(0)

export const dislikeComment = (postID: number, commentID: number) =>
  document.cookie.includes("forum-token")
    ? axios
        .post<number>(`/api/posts/${postID}/comment/${commentID}/dislike`, null, {
          withCredentials: true,
        })
        .then((res) => res.data)
    : Promise.resolve(0)
