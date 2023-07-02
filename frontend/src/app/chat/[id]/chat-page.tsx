"use client"

import { FC, useEffect } from "react"
import { useParams, useRouter } from "next/navigation"
import { useDispatch } from "react-redux"

import { chatActions } from "@/store/chats"
import { useChat } from "@/api/chats/chats"
import ChatInfo from "@/components/chats/ChatInfo"
import ChatMessages from "@/components/chats/Messages"
import WriteMessageForm from "@/components/chats/WriteMessageForm"

const ChatPage = () => {
  const chatId = Number(useParams().id)

  const router = useRouter()
  const dispatch = useDispatch()

  useEffect(() => {
    if (isNaN(chatId) || chatId < 0) {
      router.push("/chat")
      return
    }
    dispatch(chatActions.setActiveChatId(chatId))

    return () => {
      dispatch(chatActions.setActiveChatId(null))
    }
  }, [chatId, dispatch, router])

  if (isNaN(chatId) || chatId < 0) {
    return null
  }

  return <Chat id={chatId} />
}

const Chat: FC<{ id: number }> = ({ id }) => {
  const { chat } = useChat(id)
  const router = useRouter()

  useEffect(() => {
    if (chat === null) {
      router.push("/chat")
    }
  }, [chat, router])

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
