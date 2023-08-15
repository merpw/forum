"use client"

import { FC, useEffect } from "react"
import { useRouter } from "next/navigation"
import { useDispatch } from "react-redux"

import { useChat } from "@/api/chats/chats"
import ChatInfo from "@/components/chats/ChatInfo"
import ChatMessages from "@/components/chats/ChatMessages"
import WriteMessageForm from "@/components/chats/WriteMessageForm"
import { chatActions } from "@/store/chats"
import { useSendMessage } from "@/api/chats/messages"

const ChatIdPage: FC<{ id: number }> = ({ id }) => {
  const { chat } = useChat(id)
  const dispatch = useDispatch()
  const router = useRouter()

  const sendMessage = useSendMessage()

  useEffect(() => {
    if (chat === null) {
      router.replace("/chat")
    }
    if (chat) {
      dispatch(chatActions.setActiveChatId(chat.id))

      return () => {
        dispatch(chatActions.setActiveChatId(null))
      }
    }
  }, [chat, dispatch, id, router])

  if (!chat) {
    return <div className={"text-info text-center mt-5 mb-7"}>loading...</div>
  }

  return (
    <div className={"flex flex-col h-full"}>
      <ChatInfo {...chat} />

      <ChatMessages chatId={chat.id} />

      <WriteMessageForm sendMessage={(content) => sendMessage(chat.id, content)} />
    </div>
  )
}

export default ChatIdPage
