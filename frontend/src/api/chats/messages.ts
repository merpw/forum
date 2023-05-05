import { useEffect } from "react"

import { useAppDispatch, useAppSelector } from "@/store/hooks"
import { sendWSGet, sendWSPost } from "@/store/wsMiddleware"

export const useChatMessages = (chatId: number) => {
  const chatMessages = useAppSelector((state) => state.chats.chatMessages?.[chatId])
  const dispatch = useAppDispatch()

  useEffect(() => {
    if (!chatMessages) {
      dispatch(sendWSGet(`/chat/${chatId}/messages`))
    }
  }, [chatId, chatMessages, dispatch])

  return { chatMessages }
}

export const useMessage = (messageId: number) => {
  const message = useAppSelector((state) => state.chats.messages?.[messageId])
  const dispatch = useAppDispatch()

  useEffect(() => {
    if (!message) {
      dispatch(sendWSGet(`/message/${messageId}`))
    }
  }, [dispatch, message, messageId])

  return { message }
}

export const useSendMessage = () => {
  const dispatch = useAppDispatch()

  return (chatId: number, content: string) => {
    dispatch(sendWSPost(`/chat/${chatId}/message`, { content }))
  }
}
