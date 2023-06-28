import { createAction } from "@reduxjs/toolkit"

export const sendGet = createAction("ws/send", (url: string) => {
  return {
    payload: {
      type: "get",
      item: {
        url,
      },
    },
  }
})

// TODO: improve types, maybe even infer data type by url
export const sendPost = createAction("ws/send", <T>(url: string, data: T) => {
  return {
    payload: {
      type: "post",
      item: {
        url,
        data,
      },
    },
  }
})

const close = createAction("ws/close")

const connect = createAction("ws/connect")

const wsActions = {
  sendGet,
  sendPost,
  close,
  connect,
}

export default wsActions
