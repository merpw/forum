import { FC, useEffect, useRef } from "react"

import { useChatMessages } from "@/api/chats/messages"
import Message from "@/components/chats/Message"

const ChatMessages: FC<{ chatId: number }> = ({ chatId }) => {
  const { chatMessages } = useChatMessages(chatId)

  const scrollRef = useRef<HTMLDivElement>(null)

  const isOnBottom = useRef(true)

  // TODO: add lazy loading

  // TODO: add dates between days

  useEffect(() => {
    // TODO: improve this
    setTimeout(() => {
      scrollRef.current?.scrollTo(0, 0)
    }, 100)
  }, [chatMessages])

  if (!chatMessages) {
    return <div>loading...</div>
  }

  return (
    <div
      ref={scrollRef}
      onScroll={(e) => {
        // scroll is reversed, so scrollTop is always negative
        isOnBottom.current = -50 < e.currentTarget.scrollTop
      }}
      className={"overflow-y-auto flex flex-col-reverse gap-1.5 items-center"}
    >
      {chatMessages
        .slice()
        .reverse()
        .map((messageId) => (
          <Message key={messageId} id={messageId} />
        ))}
    </div>
  )
}

export default ChatMessages
