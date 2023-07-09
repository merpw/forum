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
        "px-3 rounded-2xl w-fit flex flex-wrap max-w-[73%]" +
        " " +
        (user?.id !== message.authorId
          ? "rounded-bl-none bg-secondary-content text-info ml-2 mr-auto justify-end"
          : "rounded-br-none bg-base-300 text-info brightness-110 mr-2 ml-auto justify-end")
      }
    >
      <div
        className={"prose dark:prose-invert max-w-full prose-img:max-h-[50vh] font-medium p-1"}
        dangerouslySetInnerHTML={{ __html: html }}
      />
      <span
        className={"mb-1 ml-2 mt-auto text-xs opacity-75"}
        title={dayjs(message.timestamp).format("YYYY-MM-DD HH:mm:ss")}
      >
        {dayjs(message.timestamp).format("HH:mm")}
      </span>
    </div>
  )
}

export default memo(Message)
