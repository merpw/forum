"use client"

import { FC, useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { useDispatch } from "react-redux"

import { useChat, useCreateChat, useUserChat } from "@/api/chats/chats"
import ChatInfo from "@/components/chats/ChatInfo"
import ChatMessages from "@/components/chats/ChatMessages"
import WriteMessageForm from "@/components/chats/WriteMessageForm"
import { chatActions } from "@/store/chats"
import { useSendMessage } from "@/api/chats/messages"

export const ChatPageFirstMessage: FC<{ userId: number }> = ({ userId }) => {
  const { chatId } = useUserChat(userId)

  const [firstMessage, setFirstMessage] = useState<string | null>(null)

  const createChat = useCreateChat()
  const sendMessage = useSendMessage()

  useEffect(() => {
    if (!firstMessage || chatId === undefined) {
      return
    }
    if (chatId === null) {
      createChat(userId)
    } else {
      sendMessage(chatId, firstMessage)
    }
  }, [chatId, createChat, firstMessage, sendMessage, userId])

  if (chatId === undefined) {
    return <div className={"text-info text-center mt-5 mb-7"}>loading...</div>
  }

  if (chatId === null) {
    return (
      <div className={"flex flex-col h-full"}>
        <ChatInfo userId={userId} />
        <div
          className={
            "flex overflow-y-auto border border-base-100 rounded-lg grow p-1 pb-3 my-3 gradient-light dark:gradient-dark"
          }
        ></div>

        <WriteMessageForm
          sendMessage={(message: string) => {
            setFirstMessage(message)
          }}
        />
      </div>
    )
  }

  return <ChatPage id={chatId} />
}

const ChatPage: FC<{ id: number }> = ({ id }) => {
  const { chat } = useChat(id)
  const dispatch = useDispatch()
  const router = useRouter()

  const sendMessage = useSendMessage()

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
    return <div className={"text-info text-center mt-5 mb-7"}>loading...</div>
  }

  return (
    <div className={"flex flex-col h-full"}>
      <ChatInfo userId={chat.companionId} />

      <ChatMessages chatId={chat.id} />

      <WriteMessageForm sendMessage={(content) => sendMessage(chat.id, content)} />
    </div>
  )
}

export default ChatPage
