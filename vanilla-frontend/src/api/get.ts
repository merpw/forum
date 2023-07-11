/* IMPORTS */

/* Root */
import { backendUrl } from "../main.js"
import { ChatUser } from "../types"

/* Authorization */
import { Auth } from "../components/authorization/auth.js"

// Gets all posts of the selected endpoint
const getPosts = async (endpoint: string): Promise<object[]> => {
  return fetch(backendUrl + endpoint)
    .then((r) => r.json())
    .then((data) => data)
}

// Gets post values of selected post
const getPostValues = async (postId: string): Promise<object> => {
  return fetch(`${backendUrl}/api/posts/${postId}`)
    .then((r) => r.json())
    .then((data) => data)
}

// Gets comments of selected post
const getComments = async (postId: string): Promise<object[]> => {
  return fetch(`${backendUrl}/api/posts/${postId}/comments`)
    .then((r) => r.json())
    .then((data) => data)
}

// Gets Me information. Used in userInfo object.
const getMe = async (): Promise<object | void> => {
  return fetch(`${backendUrl}/api/me`)
    .then((r) => r.json())
    .then((data) => {
      return {
        Id: data.id,
        Name: data.name,
      }
    })
    .catch(() => {
      return Auth(false)
    })
}

// Gets user info by id
// TODO: Add online = true if ID is in online users array?
const getUserById = async (id: string): Promise<ChatUser> => {
  return await fetch(`${backendUrl}/api/users/${id}`)
    .then((r) => r.json())
    .then((data) => {
      return <ChatUser>{
        Id: data.id,
        Name: data.name,
        UnreadMsg: false,
        Online: false,
      }
    })
}

// Gets all user IDs in alphabetical order
const getUserIds = async (): Promise<number[]> => {
  const ids: number[] = []
  return await fetch(`${backendUrl}/api/users`)
    .then((r) => r.json())
    .then((data) => {
      for (const id of data) {
        ids.push(id)
      }
      return ids
    })
}

export { getPosts, getPostValues, getComments, getMe, getUserIds, getUserById }
