import { FC, memo, useEffect, useState } from "react"
import dayjs from "dayjs"
import Link from "next/link"

import { useMe } from "@/api/auth/hooks"
import { useMessage } from "@/api/chats/messages"
import RenderMarkdown from "@/components/markdown/render"
import { Message as MessageType } from "@/ws"
import { useAppSelector } from "@/store/hooks"
import { useUser } from "@/api/users/hooks"
import Avatar from "@/components/Avatar"
import UserName from "@/components/chats/UserName"

import "highlight.js/styles/github-dark.css"

const Message: FC<{ id: number }> = ({ id }) => {
  const { message } = useMessage(id)

  if (!message) {
    return null
  }

  return (
    <>
      {message.authorId === -1 ? (
        <span
          className={"mx-auto text-info text-xs italic"}
          title={dayjs(message.timestamp).format("YYYY-MM-DD HH:mm:ss")}
        >
          {message.content} at {dayjs(message.timestamp).format("HH:mm")}
        </span>
      ) : (
        <MarkdownMessage message={message} />
      )}
    </>
  )
}

const MarkdownMessage: FC<{ message: MessageType }> = ({ message }) => {
  const [html, setHtml] = useState<string>()

  const { user } = useMe()

  const isGroup = useAppSelector((state) => {
    const activeChat = state.chats.chats[state.chats.activeChatId as number] || {}
    return "groupId" in activeChat
  })

  const messagePosition = useAppSelector<"first" | "last" | "single" | null>((state) => {
    if (!isGroup) return null

    const activeChatMessages = state.chats.chatMessages[state.chats.activeChatId as number] || []

    const position = activeChatMessages.indexOf(message.id)
    const nextMessage = state.chats.messages[activeChatMessages[position + 1]]
    const prevMessage = state.chats.messages[activeChatMessages[position - 1]]

    const isNextMessageFromSameAuthor = nextMessage?.authorId === message.authorId
    const isPrevMessageFromSameAuthor = prevMessage?.authorId === message.authorId

    if (!isNextMessageFromSameAuthor && !isPrevMessageFromSameAuthor) {
      return "single"
    }
    if (isNextMessageFromSameAuthor && !isPrevMessageFromSameAuthor) {
      return "first"
    }
    if (!isNextMessageFromSameAuthor && isPrevMessageFromSameAuthor) {
      return "last"
    }
    return null
  })

  const showSenderDetails = isGroup && user && user.id !== message.authorId

  const shouldShowAvatar =
    showSenderDetails && (messagePosition === "last" || messagePosition === "single")

  const shouldShowName =
    showSenderDetails && (messagePosition === "first" || messagePosition === "single")

  useEffect(() => {
    RenderMarkdown(message.content).then((html) => {
      setHtml(html)
    })
  })

  if (!html) {
    return null
  }

  return (
    <div
      className={
        "relative flex max-w-[90%] sm:max-w-[70%] w-fit" +
        " " +
        (user?.id !== message.authorId ? "mr-auto" : "ml-auto")
      }
    >
      {shouldShowAvatar && <AvatarWrapper userId={message.authorId} />}
      <div
        className={
          "px-3 rounded-2xl w-fit" +
          " " +
          (showSenderDetails ? "ml-12" : "") +
          " " +
          (user?.id !== message.authorId
            ? "rounded-bl-none bg-secondary-content text-info justify-end"
            : "rounded-br-none bg-base-300 text-info brightness-110 justify-end")
        }
      >
        {shouldShowName && (
          <span className={"badge bg-base-300 text-white"}>
            <UserName userId={user.id} />
          </span>
        )}
        <div className={"flex flex-wrap"}>
          <div
            className={
              "prose dark:prose-invert max-w-full prose-img:max-h-[50vh] font-medium p-1 break-all"
            }
            dangerouslySetInnerHTML={{ __html: html }}
          />
          <span
            className={"mb-1 ml-2 mt-auto text-xs bg-opacity-80"}
            title={dayjs(message.timestamp).format("YYYY-MM-DD HH:mm:ss")}
          >
            {dayjs(message.timestamp).format("HH:mm")}
          </span>
        </div>
      </div>
    </div>
  )
}

const AvatarWrapper: FC<{ userId: number }> = ({ userId }) => {
  const { user } = useUser(userId)

  if (!user) {
    return null
  }

  return (
    <Link href={`/user/${user.id}`} title={user.username} className={"w-10 h-10 absolute bottom-0"}>
      <Avatar user={user} size={10} />
    </Link>
  )
}

export default memo(Message)
