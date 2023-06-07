import { backendUrl } from "../main.js"

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


export { getPosts, getPostValues, getComments, getMe }
