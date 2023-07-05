import { FC, useRef } from "react"
import dayjs from "dayjs"

import { useMe } from "@/api/auth/hooks"
import { useMessage } from "@/api/chats/messages"
import Markdown from "@/components/markdown/markdown"

const Message: FC<{ id: number }> = ({ id }) => {
  const { user } = useMe()
  const { message } = useMessage(id)

  const ref = useRef<HTMLDivElement>(null)

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
        <div
          ref={ref}
          className={
            "py-1 px-3 rounded w-fit flex flex-wrap max-w-[85%]" +
            " " +
            (user?.id !== message.authorId
              ? "rounded-bl-none dark:bg-gray-800 bg-blue-200 mr-auto justify-end"
              : "rounded-br-none dark:bg-gray-900 bg-gray-200 ml-auto justify-end")
          }
        >
          <Markdown
            className={"prose-img:max-h-[50vh]"}
            content={message.content}
            fallback={message.content}
          />
          <span
            className={"ml-2 mt-auto text-sm opacity-75"}
            title={dayjs(message.timestamp).format("YYYY-MM-DD HH:mm:ss")}
          >
            {dayjs(message.timestamp).format("HH:mm")}
          </span>
        </div>
      )}
    </>
  )
}

export default Message
