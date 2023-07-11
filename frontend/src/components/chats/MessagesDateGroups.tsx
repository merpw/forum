import { FC, memo, useMemo } from "react"
import dayjs from "dayjs"

import { useMessages } from "@/api/chats/messages"
import Message from "@/components/chats/Message"

const MessagesDateGroups: FC<{ messageIds: number[]; showStickyDate: boolean }> = ({
  messageIds,
  showStickyDate,
}) => {
  const { messages: reversedMessages } = useMessages(messageIds.slice().reverse())

  const { groupedMessages, processingMessages } = useMemo(() => {
    const messages = reversedMessages.slice().reverse()
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
  }, [reversedMessages])

  const groups = Object.entries(groupedMessages)

  return (
    <>
      {processingMessages.map((_, key) => {
        // TODO: add placeholders
        return <div className={"h-20"} key={key} />
      })}
      {groups.map(([date, messages], key) => (
        <span key={date} className={"relative flex flex-col gap-1.5 items-center"}>
          <span
            className={
              "z-50 text-xs text-info bg-primary bg-opacity-10 backdrop-brightness-200 dark:backdrop-brightness-50 backdrop-blur-sm px-3 rounded-xl" +
              " " +
              (key === groups.length - 1 && !showStickyDate ? "" : "sticky top-0")
            }
          >
            {formatDate(date)}
          </span>
          {key === 0 &&
            processingMessages.map((_, key) => {
              // TODO: add placeholders
              return <div className={"h-20"} key={key} />
            })}
          {messages.map((messageId) => (
            <Message key={messageId} id={messageId} />
          ))}
        </span>
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
