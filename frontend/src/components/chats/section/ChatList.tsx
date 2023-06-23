"use client"

import { FC, useContext, useEffect, useRef } from "react"
import Link from "next/link"
import dayjs from "dayjs"

import { useChat } from "@/api/chats/chats"
import { useUser } from "@/api/users/hooks"
import { useMessage } from "@/api/chats/messages"
import { useMe } from "@/api/auth/hooks"
import { MarkdownToPlain } from "@/components/markdown/render"
import { ChatSectionCollapsedContext } from "@/components/layout"
import { Message } from "@/ws"
import { useAppSelector } from "@/store/hooks"
import Avatar from "@/components/Avatar"

const ChatList: FC<{ chatIds: number[] }> = ({ chatIds }) => {
  return (
    <div>
      <h1 className={"text-xl mb-3"}>Total: {chatIds.length} chats</h1>
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

  const { isCollapsed, setIsCollapsed } = useContext(ChatSectionCollapsedContext)

  const activeChatId = useAppSelector((state) => state.chats.activeChatId)

  if (chat === undefined) {
    return <div>loading...</div>
  }
  if (chat === null) {
    return null
  }

  return (
    <Link
      href={`/chat/${chatId}`}
      className={"w-full"}
      onClick={() => {
        if (!isCollapsed && window.innerWidth < 640) {
          setIsCollapsed(true)
        }
      }}
    >
      <div
        className={
          "border border-base-300 bg-base-200 p-3 pb-1.5 rounded-lg" +
          " " +
          (chatId === activeChatId && "border-base-200 gradient-light dark:gradient-dark shadow-lg")
        }
      >
        <div className={"flex gap-1"}>
          <div className={"mr-2 mt-2 w-14"}>
            <Avatar userId={chat.companionId} />
          </div>
          <div className={"w-full space-y-0.5"}>
            <CompanionData userId={chat.companionId} />
            <MessageInfo id={chat.lastMessageId} />
          </div>
        </div>
      </div>
    </Link>
  )
}

const CompanionData: FC<{ userId: number }> = ({ userId }) => {
  const { user } = useUser(userId)
  const isOnline = useIsUserOnline(userId)

  if (user === undefined) {
    return <div>loading...</div>
  }
  if (user === null) {
    return <div>User not found</div>
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
    return <div>loading...</div>
  }

  if (message === null) {
    return <div>Message not found</div>
  }

  return (
    <>
      <div>
        {message.authorId !== -1 &&
          (message.authorId === user?.id ? "You: " : <UserName userId={message.authorId} />)}
        {MarkdownToPlain(message.content, { async: false, removeNewLines: true, limit: 50 })}
      </div>
      <div className={"flex justify-end"}>
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

const UserName: FC<{ userId: number }> = ({ userId }) => {
  const { user: author } = useUser(userId)

  return <>{author?.username}: </>
}

export default ChatList
