"use client"

import { FC, useContext, useEffect, useRef } from "react"
import Link from "next/link"
import dayjs from "dayjs"

import { useChat } from "@/api/chats/chats"
import { useIsUserOnline, useUser } from "@/api/users/hooks"
import { useMessage } from "@/api/chats/messages"
import { useMe } from "@/api/auth/hooks"
import { MarkdownToPlain } from "@/components/markdown/render"
import { ChatSectionCollapsedContext } from "@/components/layout"
import { Message } from "@/ws"
import { useAppSelector } from "@/store/hooks"

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
      <div className={"border p-4 rounded" + " " + (chatId === activeChatId && "border-blue-900")}>
        <CompanionData userId={chat.companionId} />
        <MessageInfo id={chat.lastMessageId} />
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
      <div className={"text-xl"}>
        {user.username} {isOnline ? "🟢" : "🔴"}
      </div>
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
