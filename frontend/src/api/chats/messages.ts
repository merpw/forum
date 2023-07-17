import { useEffect, useMemo, useRef } from "react"

import { useAppDispatch, useAppSelector } from "@/store/hooks"
import wsActions from "@/store/ws/actions"
import { TypingData } from "@/ws"
import { chatActions } from "@/store/chats"

export const useChatMessages = (chatId: number) => {
  const chatMessages = useAppSelector((state) => state.chats.chatMessages?.[chatId])
  const dispatch = useAppDispatch()

  useEffect(() => {
    if (!chatMessages) {
      dispatch(wsActions.sendGet(`/chat/${chatId}/messages`))
    }
  }, [chatId, chatMessages, dispatch])

  return { chatMessages }
}

export const useMessage = (messageId: number) => {
  const message = useAppSelector((state) => state.chats.messages?.[messageId])
  const dispatch = useAppDispatch()

  useEffect(() => {
    if (!message) {
      dispatch(wsActions.sendGet(`/message/${messageId}`))
    }
  }, [dispatch, message, messageId])

  return { message }
}

export const useMessages = (messageIds: number[]) => {
  const allMessages = useAppSelector((state) => state.chats.messages)

  const messages = useMemo(() => {
    return messageIds.map((id) => allMessages?.[id])
  }, [allMessages, messageIds])

  const dispatch = useAppDispatch()

  useEffect(() => {
    messageIds.forEach((id) => {
      if (!allMessages?.[id]) {
        dispatch(wsActions.sendGet(`/message/${id}`))
      }
    })
  }, [dispatch, messageIds, allMessages])

  return { messages }
}

export const useSendMessage = () => {
  const dispatch = useAppDispatch()

  return (chatId: number, content: string) => {
    dispatch(wsActions.sendPost(`/chat/${chatId}/message`, { content }))
  }
}

export const useChatTyping = (chatId: number): TypingData | null => {
  const typingData = useAppSelector((state) => state.chats.chatsTyping?.[chatId])
  const timeout = useRef<NodeJS.Timeout>()

  const dispatch = useAppDispatch()

  useEffect(() => {
    if (!typingData) {
      return
    }

    if (timeout.current) {
      clearTimeout(timeout.current)
    }

    timeout.current = setTimeout(() => {
      dispatch(chatActions.resetChatTyping(chatId))
    }, 3000)
  }, [chatId, dispatch, typingData])

  return typingData ?? null
}
