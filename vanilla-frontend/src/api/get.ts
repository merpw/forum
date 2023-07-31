import { User } from "../types"

import { Auth } from "../components/authorization/auth.js"

const getPosts = async (endpoint: string): Promise<object[]> => {
  return fetch(endpoint)
    .then((r) => r.json())
    .then((data) => data)
}

const getPostValues = async (postId: number): Promise<object> => {
  return fetch(`/api/posts/${postId}`)
    .then((r) => r.json())
    .then((data) => data)
}

const getComments = async (postId: number): Promise<object[]> => {
  return fetch(`/api/posts/${postId}/comments`)
    .then((r) => r.json())
    .then((data) => data)
}

const getMe = async (): Promise<User | void> => {
  return fetch(`/api/me`)
    .then((r) => r.json())
    .then((data) => {
      return {
        id: data.id,
        name: data.username,
      }
    })
    .catch(() => {
      return Auth(false)
    })
}

const getUserById = async (id: number): Promise<User> => {
  return await fetch(`/api/users/${id}`)
    .then((r) => r.json())
    .then((data) => {
      return {
        id: data.id,
        name: data.username,
      }
    })
}

const getUserIds = async (): Promise<number[]> => {
  return await fetch(`/api/users`)
    .then((r) => r.json())
    .then(([...ids]) => ids as number[])
}

export { getPosts, getPostValues, getComments, getMe, getUserIds, getUserById }
