import { FC, memo, useMemo } from "react"
import dayjs from "dayjs"

import { useMessages } from "@/api/chats/messages"
import Message from "@/components/chats/Message"

const MessagesDateGroups: FC<{ messageIds: number[] }> = ({ messageIds }) => {
  const { messages } = useMessages(messageIds)

  const groupedMessages = useMemo(
    () =>
      messages.reduce((acc, message) => {
        if (!message) {
          return acc
        }
        const date = dayjs(message.timestamp).format("YYYY-MM-DD")
        if (!acc[date]) {
          acc[date] = []
        }
        acc[date].push(message.id)
        return acc
      }, {} as Record<string, number[]>),
    [messages]
  )

  return (
    <>
      {Object.entries(groupedMessages).map(([date, messages]) => (
        <div key={date} className={"relative flex flex-col gap-1.5 items-center"}>
          <div
            className={
              "sticky top-1 text-base-content bg-base bg-opacity-80 backdrop-blur px-2 rounded-xl"
            }
          >
            {formatDate(date)}
          </div>
          {messages.map((messageId) => (
            <Message key={messageId} id={messageId} />
          ))}
        </div>
      ))}
    </>
  )
}

const formatDate = (timestamp: string) => {
  const date = dayjs(timestamp)
  const today = dayjs()

  if (date.format("YYYY-MM-DD") === today.format("YYYY-MM-DD")) {
    return "Today"
  }

  if (date.add(1, "day").format("YYYY-MM-DD") === today.format("YYYY-MM-DD")) {
    return "Yesterday"
  }

  return date.format(date.year() === dayjs().year() ? "MMMM D" : "MMMM D, YYYY")
}

export default memo(MessagesDateGroups)
