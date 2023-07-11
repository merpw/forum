import { FC, useEffect, useRef, useState } from "react"

import { useChatMessages } from "@/api/chats/messages"
import { useMe } from "@/api/auth/hooks"
import { useAppSelector } from "@/store/hooks"
import throttle from "@/helpers/throttle"
import MessageGroups from "@/components/chats/MessagesDateGroups"

const ChatMessages: FC<{ chatId: number }> = ({ chatId }) => {
  const { chatMessages } = useChatMessages(chatId)

  const { user } = useMe()

  const lastMessage = useAppSelector((state) => {
    const lastMessageId = state.chats.chats[chatId]?.lastMessageId
    if (!lastMessageId) {
      return
    }
    return state.chats.messages[lastMessageId]
  })

  useEffect(() => {
    if (lastMessage?.authorId === user?.id) {
      isOnBottom.current = true
      setTimeout(() => {
        // fallback, sometimes the observer misses the scroll (when the message is too small)
        scrollRef.current?.scrollTo(0, scrollRef.current?.scrollHeight)
      }, 200)
    }
  }, [chatMessages, lastMessage?.authorId, user?.id])

  const scrollRef = useRef<HTMLDivElement>(null)

  const isOnBottom = useRef(true)

  const [messagesCount, setMessagesCount] = useState(10)

  const [hasUnread, setHasUnread] = useState(false)

  const onScroll = throttle(() => {
    if (!scrollRef.current) {
      return
    }

    isOnBottom.current =
      scrollRef.current.scrollHeight -
        scrollRef.current.scrollTop -
        scrollRef.current.clientHeight <
      50

    if (isOnBottom.current && hasUnread) {
      setHasUnread(false)
    }

    if (
      scrollRef.current.scrollTop < scrollRef.current.clientHeight / 2 &&
      messagesCount < (chatMessages?.length ?? 0)
    ) {
      if (scrollRef.current.scrollTop === 0) {
        scrollRef.current.scrollTo(0, 1)
      }
      setMessagesCount((prev) => prev + 10)
    }
  }, 500)

  useEffect(() => {
    setTimeout(() => {
      if (!scrollRef.current || !observerRef.current) {
        return
      }

      const observer = new ResizeObserver(() => {
        if (isOnBottom.current) {
          scrollRef.current?.scrollTo(0, scrollRef.current?.scrollHeight)
          setTimeout(() => {
            scrollRef.current?.scrollTo(0, scrollRef.current?.scrollHeight)
          }, 100)
        }
      })

      observer.observe(observerRef.current)
      observer.observe(scrollRef.current)
    }, 100)
  }, [])

  useEffect(() => {
    setTimeout(() => {
      if (
        messagesCount < (chatMessages?.length ?? 0) &&
        scrollRef.current?.scrollHeight === scrollRef.current?.clientHeight
      ) {
        setMessagesCount((prev) => prev + 10)
      }
    }, 100)
  }, [chatMessages, messagesCount])

  const observerRef = useRef<HTMLDivElement>(null)

  if (!chatMessages) {
    return null
  }

  const messages = chatMessages.slice(-messagesCount)

  return (
    <div ref={scrollRef} onScroll={onScroll} className={"overflow-y-auto"}>
      <div ref={observerRef}>
        <MessageGroups messageIds={messages} />
      </div>
    </div>
  )
}

export default ChatMessages
