"use client"

import { FC, useContext, useEffect, useRef } from "react"
import Link from "next/link"
import { useParams } from "next/navigation"
import dayjs from "dayjs"

import { useChat } from "@/api/chats/chats"
import { useUser } from "@/api/users/hooks"
import { useMessage } from "@/api/chats/messages"
import { useMe } from "@/api/auth"
import { MarkdownToPlain } from "@/components/markdown/render"
import { ChatSectionCollapsedContext } from "@/components/layout"
import { Message } from "@/ws"

const ChatList: FC<{ chatIds: number[] }> = ({ chatIds }) => {
  const { id } = useParams()

  return (
    <div>
      <h1 className={"text-xl mb-3"}>Total: {chatIds.length} chats</h1>
      <ul className={"flex flex-col gap-2"}>
        {chatIds.map((chatId) => (
          <li key={chatId} className={"w-full"}>
            <ChatCard id={chatId} isActive={Number(id) === chatId} />
          </li>
        ))}
      </ul>
    </div>
  )
}

const ChatCard: FC<{ id: number; isActive?: boolean }> = ({ id, isActive = false }) => {
  const { chat } = useChat(id)

  const { user: lastMessageUser } = useUser(chat?.companionId)

  const { isCollapsed, setIsCollapsed } = useContext(ChatSectionCollapsedContext)

  if (chat === undefined) {
    return <div>loading...</div>
  }
  if (chat === null) {
    return null
  }
  return (
    <Link
      href={`/chat/${id}`}
      className={"w-full"}
      onClick={() => {
        if (!isCollapsed && window.innerWidth < 640) {
          setIsCollapsed(true)
        }
      }}
    >
      <div className={"border p-4 rounded" + " " + (isActive && "border-blue-900")}>
        <div className={"text-xl"}>{lastMessageUser?.name}</div>
        <MessageInfo id={chat.lastMessageId} />
      </div>
    </Link>
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

  return <>{author?.name}: </>
}

export default ChatList
