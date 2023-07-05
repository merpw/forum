"use client"

import { FC, useEffect } from "react"
import { useRouter } from "next/navigation"
import { useDispatch } from "react-redux"

import { useChat } from "@/api/chats/chats"
import ChatInfo from "@/components/chats/ChatInfo"
import ChatMessages from "@/components/chats/ChatMessages"
import WriteMessageForm from "@/components/chats/WriteMessageForm"
import { chatActions } from "@/store/chats"

const ChatPage: FC<{ id: number }> = ({ id }) => {
  const { chat } = useChat(id)
  const dispatch = useDispatch()
  const router = useRouter()

  useEffect(() => {
    if (chat === null) {
      router.push("/chat")
    }
    if (chat) {
      dispatch(chatActions.setActiveChatId(chat.id))

      return () => {
        dispatch(chatActions.setActiveChatId(null))
      }
    }
  }, [chat, dispatch, id, router])

  if (!chat) {
    return <div>loading...</div>
  }

  return (
    <div className={"flex flex-col h-full"}>
      <ChatInfo userId={chat.companionId} />
      <ChatMessages chatId={chat.id} />
      <WriteMessageForm chatId={chat.id} />
    </div>
  )
}

export default ChatPage
