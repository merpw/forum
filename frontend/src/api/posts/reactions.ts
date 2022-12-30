import axios from "axios"

export const getPostReaction = (postID: number) =>
  document.cookie.includes("forum-token")
    ? axios
        .get<number>(`/api/posts/${postID}/reaction`, { withCredentials: true })
        .then((res) => res.data)
        .catch(() => 0)
    : Promise.resolve(0)

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
