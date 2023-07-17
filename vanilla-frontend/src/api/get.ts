import { ChatUser } from "../types"

import { Auth } from "../components/authorization/auth.js"

const getPosts = async (endpoint: string): Promise<object[]> => {
  return fetch(endpoint)
    .then(r => r.json())
    .then(data => data)
}

const getPostValues = async (postId: string): Promise<object> => {
  return fetch(`/api/posts/${postId}`)
    .then(r => r.json())
    .then(data => data)
}

const getComments = async (postId: string): Promise<object[]> => {
  return fetch(`/api/posts/${postId}/comments`)
    .then(r => r.json())
    .then(data => data)
}

const getMe = async (): Promise<object | void> => {
  return fetch(`/api/me`)
    .then(r => r.json())
    .then((data) => {
      return {
        Id: data.id,
        Name: data.username,
      }
    })
    .catch(() => {
      return Auth(false)
    })
}

const getUserById = async (id: string): Promise<ChatUser> => {
  return await fetch(`/api/users/${id}`)
    .then((r) => r.json())
    .then((data) => {
      return <ChatUser>{
        Id: data.id,
        Name: data.username,
        UnreadMsg: false,
        Online: false,
      }
    })
}

const getUserIds = async (): Promise<number[]> => {
  return await fetch(`/api/users`)
    .then((r) => r.json())
    .then(([...ids]) => ids as number[])
}

export { getPosts, getPostValues, getComments, getMe, getUserIds, getUserById }
