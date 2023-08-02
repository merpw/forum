"use client"

import { FC, useEffect, useRef } from "react"
import Link from "next/link"
import dayjs from "dayjs"

import { useChat } from "@/api/chats/chats"
import { useUser } from "@/api/users/hooks"
import { useChatTyping, useMessage } from "@/api/chats/messages"
import { useMe } from "@/api/auth/hooks"
import { MarkdownToPlain } from "@/components/markdown/render"
import { Message } from "@/ws"
import { useAppSelector } from "@/store/hooks"
import Avatar from "@/components/Avatar"
import UserName from "@/components/chats/UserName"
import { useCollapseIfMobile } from "@/components/chats/section/ChatSection"

const ChatList: FC<{ chatIds: number[] }> = ({ chatIds }) => {
  return (
    <div>
      <h1 className={"text-info text-sm my-1"}>Total: {chatIds.length} chats</h1>
      <ul className={"flex flex-col gap-2"}>
        {chatIds.map((chatId) => (
          <li key={chatId} className={"w-full"}>
            <ChatCard chatId={chatId} />
          </li>
        ))}
      </ul>
    </div>
  )
}

const ChatCard: FC<{ chatId: number }> = ({ chatId }) => {
  const { chat } = useChat(chatId)

  const collapseIfMobile = useCollapseIfMobile()

  const activeChatId = useAppSelector((state) => state.chats.activeChatId)

  const typingData = useChatTyping(chatId)

  const unreadMessagesCount = useAppSelector(
    (state) => state.chats.unreadMessagesChatIds.filter((id) => id === chatId).length
  )

  if (chat === undefined) {
    return <div className={"text-info"}>loading...</div>
  }
  if (chat === null) {
    return null
  }

  return (
    <Link href={`/chat/${chatId}`} className={"relative w-full"} onClick={collapseIfMobile}>
      <div
        className={
          "bg-base-200 p-3 pt-2 pb-1 rounded-lg hover:bg-neutral break-all flex gap-1" +
          " " +
          (chatId === activeChatId ? "gradient-light dark:gradient-dark shadow-sm" : "") +
          " " +
          (unreadMessagesCount > 0 ? "ring-2 ring-secondary" : "")
        }
      >
        <CompanionAvatar userId={chat.companionId} />
        {unreadMessagesCount > 0 && (
          <div className={"absolute badge badge-secondary top-3 right-3 font-bold py-3"}>
            {unreadMessagesCount}
          </div>
        )}
        <div className={"w-full"}>
          <CompanionData userId={chat.companionId} />
          {typingData ? (
            <div className={"flex items-center gap-1 mb-5 text-sm"}>
              <UserName userId={typingData.userId} /> is typing
              <span
                className={"loading loading-dots loading-xs mt-2 opacity-60 dark:opacity-80"}
              ></span>
            </div>
          ) : (
            <MessageInfo id={chat.lastMessageId} />
          )}
        </div>
      </div>
    </Link>
  )
}

const CompanionAvatar: FC<{ userId: number }> = ({ userId }) => {
  const { user } = useUser(userId)

  if (!user) {
    return null
  }

  return <Avatar user={user} size={50} className={"mr-2 mt-2 w-12 mb-auto"} />
}

const CompanionData: FC<{ userId: number }> = ({ userId }) => {
  const { user } = useUser(userId)

  if (user === undefined) {
    return <div className={"text-info"}>loading...</div>
  }
  if (user === null) {
    return <div className={"text-info text-center"}>User not found</div>
  }

  return (
    <div>
      <div className={"font-Alatsi text-lg text-primary"}>{user.username}</div>
    </div>
  )
}

const MessageInfo: FC<{ id: number }> = ({ id }) => {
  const fallbackMessage = useRef<null | Message>()
  const { message: newMessage } = useMessage(id)
  const { user } = useMe()

  const message = newMessage ?? fallbackMessage.current

  useEffect(() => {
    if (newMessage !== undefined) {
      fallbackMessage.current = newMessage
    }
  }, [newMessage])

  if (message === undefined) {
    return <div className={"text-info"}>loading...</div>
  }

  if (message === null) {
    return <div className={"text-info"}>Message not found</div>
  }

  return (
    <>
      <div
        className={"text-sm " + (message.authorId === -1 ? " text-info italic font-light " : "")}
      >
        {message.authorId !== -1 &&
          (message.authorId === user?.id ? (
            <span className={"end-dot text-info font-light"}>You</span>
          ) : (
            ""
          ))}
        {MarkdownToPlain(message.content, { async: false, removeNewLines: true, limit: 50 })}
      </div>
      <div className={"flex justify-end text-info text-xs mt-1"}>
        <FormattedDate timestamp={message.timestamp} />
      </div>
    </>
  )
}

const FormattedDate: FC<{ timestamp: string }> = ({ timestamp }) => {
  const date = dayjs(timestamp)
  return (
    <div title={date.format("YYYY-MM-DD HH:mm:ss")}>
      {date.diff(dayjs(), "day") === 0 ? date.format("HH:mm") : date.format("YYYY-MM-DD HH:mm")}
    </div>
  )
}

export default ChatList
