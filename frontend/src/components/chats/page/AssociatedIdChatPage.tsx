"use client"

import { FC, useEffect, useState } from "react"

import { useAssociatedChat, useCreateChat } from "@/api/chats/chats"
import { useSendMessage } from "@/api/chats/messages"
import ChatInfo from "@/components/chats/ChatInfo"
import WriteMessageForm from "@/components/chats/WriteMessageForm"
import ChatIdPage from "@/components/chats/page/ChatIdPage"

export const AssociatedIdChatPage: FC<{ userId: number } | { groupId: number }> = (props) => {
  const { chatId } = useAssociatedChat(props)

  const [firstMessage, setFirstMessage] = useState<string | null>(null)

  const createChat = useCreateChat()
  const sendMessage = useSendMessage()

  useEffect(() => {
    if (!firstMessage || chatId === undefined) {
      return
    }
    if (chatId === null) {
      createChat(props)
    } else {
      sendMessage(chatId, firstMessage)
    }
  }, [chatId, createChat, firstMessage, props, sendMessage])

  if (chatId === undefined) {
    return <div className={"text-info text-center mt-5 mb-7"}>loading...</div>
  }

  if (chatId === null) {
    return (
      <div className={"flex flex-col h-full"}>
        <ChatInfo {...props} />
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

  return <ChatIdPage id={chatId} />
}
