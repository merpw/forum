import { FC, memo, useMemo } from "react"
import dayjs from "dayjs"

import { useMessages } from "@/api/chats/messages"
import Message from "@/components/chats/Message"

const MessagesDateGroups: FC<{ messageIds: number[]; showStickyDate: boolean }> = ({
  messageIds,
  showStickyDate,
}) => {
  // load the most recent messages first
  messageIds.reverse()

  const { messages } = useMessages(messageIds)

  // reverse back to original order
  messages.reverse()

  const { groupedMessages, processingMessages } = useMemo(() => {
    const groupedMessages = messages.reduce((acc, message) => {
      if (!message) {
        return acc
      }
      const date = dayjs(message.timestamp).format("YYYY-MM-DD")
      if (!acc[date]) {
        acc[date] = []
      }
      acc[date].push(message.id)
      return acc
    }, {} as Record<string, number[]>)

    const processingMessages = messages.filter((message) => !message)

    return { groupedMessages, processingMessages }
  }, [messages])

  const groups = Object.entries(groupedMessages)

  return (
    <>
      {processingMessages.map((_, key) => {
        // TODO: add placeholders
        return <div className={"h-20"} key={key} />
      })}
      {groups.map(([date, messages], key) => (
        <div key={date} className={"relative flex flex-col gap-1.5 items-center"}>
          <div
            className={
              "pt-1 text-base-content bg-base bg-opacity-80 backdrop-blur px-2 rounded-xl" +
              " " +
              (key === groups.length - 1 && !showStickyDate ? "" : "sticky top-0")
            }
          >
            {formatDate(date)}
          </div>
          {key === 0 &&
            processingMessages.map((_, key) => {
              // TODO: add placeholders
              return <div className={"h-20"} key={key} />
            })}
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
