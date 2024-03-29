import { FC, useEffect, useMemo, useRef, useState } from "react"

import { useChatMessages, useChatTyping } from "@/api/chats/messages"
import { useMe } from "@/api/auth/hooks"
import { useAppSelector } from "@/store/hooks"
import throttle from "@/helpers/throttle"
import MessageGroups from "@/components/chats/MessagesDateGroups"
import UserName from "@/components/chats/UserName"

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

  const showStickyDateTimeout = useRef<NodeJS.Timeout | null>(null)
  const [showStickyDate, setShowStickyDate] = useState(false)

  const onScroll = useMemo(
    () =>
      throttle(() => {
        if (!scrollRef.current) {
          return
        }

        if (!showStickyDate) {
          setShowStickyDate(true)
        }

        if (showStickyDateTimeout.current) {
          clearTimeout(showStickyDateTimeout.current)
        }
        showStickyDateTimeout.current = setTimeout(() => {
          if (showStickyDate) {
            setShowStickyDate(false)
          }
        }, 1000)

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
          if (scrollRef.current.scrollTop < 10) {
            scrollRef.current.scrollTop = 10
          }
          setMessagesCount((prev) => prev + 10)
        }
      }, 500),
    [chatMessages, hasUnread, messagesCount, showStickyDate]
  )

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
    <div
      ref={scrollRef}
      onScroll={onScroll}
      className={
        "flex overflow-y-auto border border-base-100 rounded-lg grow p-1 pb-3 my-3 gradient-light dark:gradient-dark"
      }
    >
      <div ref={observerRef} className={"mt-auto w-full"}>
        <MessageGroups messageIds={messages} showStickyDate={showStickyDate} />
        <ChatTyping chatId={chatId} />
      </div>
    </div>
  )
}

const ChatTyping: FC<{ chatId: number }> = ({ chatId }) => {
  const typingData = useChatTyping(chatId)

  if (!typingData) {
    return null
  }

  return (
    <div className={"flex items-center gap-1 ml-3 mt-1 text-primary"}>
      <span>
        <UserName userId={typingData.userId} /> is typing
      </span>
      <span className={"loading loading-dots loading-xs mt-auto"}></span>
    </div>
  )
}

export default ChatMessages
