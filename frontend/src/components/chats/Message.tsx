import { FC, memo, useEffect, useState } from "react"
import dayjs from "dayjs"

import { useMe } from "@/api/auth/hooks"
import { useMessage } from "@/api/chats/messages"
import RenderMarkdown from "@/components/markdown/render"
import { Message as MessageType } from "@/ws"

const Message: FC<{ id: number }> = ({ id }) => {
  const { message } = useMessage(id)

  if (!message) {
    return null
  }

  return (
    <>
      {message.authorId === -1 ? (
        <span
          className={"mx-auto italic"}
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
        "py-1 px-3 rounded w-fit flex flex-wrap max-w-[85%]" +
        " " +
        (user?.id !== message.authorId
          ? "rounded-bl-none dark:bg-gray-800 bg-blue-200 mr-auto justify-end"
          : "rounded-br-none dark:bg-gray-900 bg-gray-200 ml-auto justify-end")
      }
    >
      <div
        className={"prose dark:prose-invert max-w-full prose-img:max-h-[50vh]"}
        dangerouslySetInnerHTML={{ __html: html }}
      />
      <span
        className={"ml-2 mt-auto text-sm opacity-75"}
        title={dayjs(message.timestamp).format("YYYY-MM-DD HH:mm:ss")}
      >
        {dayjs(message.timestamp).format("HH:mm")}
      </span>
    </div>
  )
}

export default memo(Message)
