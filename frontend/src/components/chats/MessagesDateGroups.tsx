import { FC, memo, useMemo } from "react"
import dayjs from "dayjs"

import { useMessages } from "@/api/chats/messages"
import Message from "@/components/chats/Message"
import { formatDate } from "@/helpers/dates"

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
      {groups.map(([date, messages], key) => (
        <span key={date} className={"relative flex flex-col gap-1.5 items-center"}>
          <span
            className={
              "z-50 text-xs text-info bg-primary bg-opacity-10 backdrop-brightness-200 dark:backdrop-brightness-50 backdrop-blur-sm mt-1 px-3 rounded-xl" +
              " " +
              (key === groups.length - 1 && !showStickyDate ? "" : "sticky top-0")
            }
          >
            {formatDate(date)}
          </span>
          {key === 0 &&
            processingMessages.map((_, key) => {
              // TODO: add placeholders
              return <div className={"h-30"} key={key} />
            })}
          {messages.map((messageId) => (
            <Message key={messageId} id={messageId} />
          ))}
        </span>
      ))}
    </>
  )
}

export default memo(MessagesDateGroups)
