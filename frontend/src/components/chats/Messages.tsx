import { FC, useEffect, useRef, useState } from "react"

import { useChatMessages } from "@/api/chats/messages"
import Message from "@/components/chats/Message"
import throttle from "@/helpers/throttle"

const ChatMessages: FC<{ chatId: number }> = ({ chatId }) => {
  const { chatMessages } = useChatMessages(chatId)

  const scrollRef = useRef<HTMLDivElement>(null)

  const isOnBottom = useRef(true)

  const [messagesCount, setMessagesCount] = useState(10)

  // TODO: add dates between days

  const onScroll = throttle(() => {
    if (!scrollRef.current) {
      return
    }

    isOnBottom.current =
      scrollRef.current.scrollHeight -
        scrollRef.current.scrollTop -
        scrollRef.current.clientHeight <
      50

    console.log(scrollRef.current.scrollTop)

    if (scrollRef.current.scrollTop < 100) {
      if (scrollRef.current.scrollTop < 10 && messagesCount < (chatMessages?.length ?? 0)) {
        scrollRef.current.scrollTo(0, 10)
      }
      setMessagesCount((prev) => prev + 10)
    }
  }, 100)

  useEffect(() => {
    setTimeout(() => {
      if (isOnBottom.current) {
        scrollRef.current?.scrollTo(0, scrollRef.current.scrollHeight)
      }

      if (
        messagesCount < (chatMessages?.length ?? 0) &&
        scrollRef.current?.scrollHeight === scrollRef.current?.clientHeight
      ) {
        setMessagesCount((prev) => prev + 10)
      }
    }, 100)
  }, [chatMessages, messagesCount])

  if (!chatMessages) {
    return <div>loading...</div>
  }

  return (
    <div
      ref={scrollRef}
      onScroll={onScroll}
      className={"overflow-y-auto flex flex-col gap-1.5 items-center"}
    >
      <div />
      {chatMessages.slice(-messagesCount).map((messageId) => (
        <Message key={messageId} id={messageId} />
      ))}
    </div>
  )
}

export default ChatMessages
