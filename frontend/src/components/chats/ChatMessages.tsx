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

    if (scrollRef.current.scrollTop < 100) {
      if (scrollRef.current.scrollTop < 10 && messagesCount < (chatMessages?.length ?? 0)) {
        scrollRef.current.scrollTo(0, 10)
      }
      setMessagesCount((prev) => prev + 10)
    }
  }, 100)

  useEffect(() => {
    if (lastMessage?.authorId === user?.id || isOnBottom.current) {
      setTimeout(() => {
        scrollRef.current?.scrollTo(0, scrollRef.current.scrollHeight)
      }, 200)
      return
    }

    if (lastMessage?.authorId !== user?.id) {
      setHasUnread(true)
    }
  }, [lastMessage, user])

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

  if (!chatMessages) {
    return <div>loading...</div>
  }

  const messages = chatMessages.slice(-messagesCount)

  return (
    <div ref={scrollRef} onScroll={onScroll} className={"overflow-y-auto"}>
      <div />
      <MessageGroups messageIds={messages} />
    </div>
  )
}

export default ChatMessages
