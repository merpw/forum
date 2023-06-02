import { useEffect } from "react"

import { useAppDispatch, useAppSelector } from "@/store/hooks"
import { wsActions } from "@/store/wsMiddleware"

export const useChatIds = () => {
  const chatIds = useAppSelector((state) => state.chats.chatIds)
  const dispatch = useAppDispatch()

  useEffect(() => {
    if (!chatIds) {
      dispatch(wsActions.sendWSGet("/chat/all"))
    }
  }, [chatIds, dispatch])

  return { chatIds }
}

export const useChat = (id: number) => {
  const chat = useAppSelector((state) => state.chats.chats?.[id])
  const dispatch = useAppDispatch()

  useEffect(() => {
    if (!chat) {
      dispatch(wsActions.sendWSGet(`/chat/${id}`))
    }
  }, [chat, dispatch, id])

  return { chat }
}

export const useCreateChat = () => {
  const dispatch = useAppDispatch()

  return (userId: number) => {
    dispatch(wsActions.sendWSPost("/chat/create", { userId }))
  }
}

export const useUserChat = (userId: number) => {
  const chatId = useAppSelector((state) => state.chats.userChats?.[userId])
  const dispatch = useAppDispatch()

  useEffect(() => {
    if (!chatId) {
      dispatch(wsActions.sendWSGet(`/users/${userId}/chat`))
    }
  }, [chatId, dispatch, userId])

  return { chatId }
}
