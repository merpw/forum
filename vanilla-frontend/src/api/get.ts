import { userInfo } from "../components/auth.js"
import { backendUrl } from "../main.js"
import { ChatUser } from "../types"

const getPosts = async (endpoint: string): Promise<object[]> => {
  return fetch(backendUrl + endpoint)
    .then((response) => response.json())
    .then((data) => data) 
}

const getPostValues = async (postId: string): Promise<object> => {
return fetch(`${backendUrl}/api/posts/${postId}`)
    .then((r) => r.json())
    .then((data) => data)
}

const getComments = async (postId: string): Promise<object[]> => {
return fetch(`${backendUrl}/api/posts/${postId}/comments`)
  .then((r) => r.json())
  .then((data) => data)
}


const getMe = async (): Promise<object> => {
  return fetch(`${backendUrl}/api/me`)
    .then((r) => r.json())
    .then((data) => {
      return {
        Id: data.id,
        Name: data.name
      }
    })
}

const getUserById = async (id: string): Promise<ChatUser> => {
  return await fetch(`${backendUrl}/api/user/${id}`)
  .then((r) => r.json())
  .then((data) => {
        return <ChatUser>{
          Id: data.id,
          Name: data.name,
          UnreadMsg: false,
          Online: false
        }
    })
}

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
