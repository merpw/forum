"use client"

import { FC, useEffect, useRef, useState } from "react"
import ReactTextAreaAutosize from "react-textarea-autosize"
import { useParams, useRouter } from "next/navigation"
import Link from "next/link"
import dayjs from "dayjs"
import { useDispatch } from "react-redux"

import Markdown from "@/components/markdown/markdown"
import { useChatMessages, useMessage, useSendMessage } from "@/api/chats/messages"
import { useIsUserOnline, useUser } from "@/api/users/hooks"
import { useChat } from "@/api/chats/chats"
import { chatActions } from "@/store/chats"
import { useMe } from "@/api/auth"

const ChatPage = () => {
  const params = useParams()
  const chatId = Number(params?.id || NaN)

  const router = useRouter()
  const dispatch = useDispatch()

  useEffect(() => {
    if (chatId && isNaN(chatId)) {
      router.push("/chat")
      return
    }
    dispatch(chatActions.setActiveChatId(chatId))

    return () => {
      dispatch(chatActions.setActiveChatId(null))
    }
  }, [chatId, dispatch, router])

  if (chatId === null || isNaN(chatId)) {
    return null
  }

  return <Chat id={chatId} />
}

const Chat: FC<{ id: number }> = ({ id }) => {
  const { chat } = useChat(id)
  const router = useRouter()

  useEffect(() => {
    if (chat === null) {
      router.push("/chat")
    }
  }, [chat, router])

  if (!chat) {
    return <div>loading...</div>
  }

  return (
    <div className={"flex flex-col h-full"}>
      <ChatInfo userId={chat.companionId} />
      <ChatMessages chatId={chat.id} />
      <WriteMessageForm chatId={chat.id} />
    </div>
  )
}

const WriteMessageForm: FC<{ chatId: number }> = ({ chatId }) => {
  const [input, setInput] = useState("")

  const formRef = useRef<HTMLFormElement>(null)
  const inputRef = useRef<HTMLTextAreaElement>(null)

  const sendMessage = useSendMessage()

  useEffect(() => {
    inputRef.current?.focus()
  }, [chatId])

  return (
    <form
      ref={formRef}
      className={"mb-6"}
      onSubmit={async (e) => {
        e.preventDefault()
        sendMessage(chatId, input)
        // await mutateMessages()
        setInput("")
      }}
    >
      <ReactTextAreaAutosize
        ref={inputRef}
        className={"textarea w-full my-3"}
        onChange={(e) => setInput(e.currentTarget.value)}
        onBlur={() => setInput(input.trim())}
        value={input}
        onKeyDown={(e) => {
          if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault()
            const trimmedInput = input.trim()
            setInput(trimmedInput)
            if (trimmedInput === "") {
              return
            }
            formRef.current?.dispatchEvent(new Event("submit", { cancelable: true, bubbles: true }))
          }
        }}
        required
      />
      <button className={"button disabled:opacity-50"} type={"submit"}>
        Send
      </button>
    </form>
  )
}

const ChatInfo: FC<{ userId: number }> = ({ userId }) => {
  const { user } = useUser(userId)
  const isOnline = useIsUserOnline(userId)

  if (!user) {
    return <div>loading...</div>
  }
  return (
    <h1 className={"text-2xl mb-auto pb-2 border-b"}>
      Chat with{" "}
      <Link href={`/user/${userId}`}>
        <span className={"font-bold clickable"}>{user?.name}</span> {isOnline ? "🟢" : "🔴"}
      </Link>
    </h1>
  )
}

const ChatMessages: FC<{ chatId: number }> = ({ chatId }) => {
  const { chatMessages } = useChatMessages(chatId)

  const scrollRef = useRef<HTMLDivElement>(null)

  const isOnBottom = useRef(true)

  // TODO: add lazy loading

  // TODO: add dates between days

  useEffect(() => {
    // TODO: maybe remove this?
    if (isOnBottom.current) {
      setTimeout(() => {
        scrollRef.current?.scrollTo(0, 0)
      }, 100)
    }
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
          <MessageComponent key={messageId} id={messageId} />
        ))}
    </div>
  )
}

const MessageComponent: FC<{ id: number }> = ({ id }) => {
  const { user } = useMe()
  const { message } = useMessage(id)

  const ref = useRef<HTMLDivElement>(null)

  if (!message) {
    return null
  }

  if (message.authorId === -1) {
    return (
      <span
        className={"mx-auto italic"}
        title={dayjs(message.timestamp).format("YYYY-MM-DD HH:mm:ss")}
      >
        {message.content} at {dayjs(message.timestamp).format("HH:mm")}
      </span>
    )
  }

  return (
    <div
      ref={ref}
      className={
        "py-1 px-3 rounded w-fit flex flex-wrap max-w-[85%]" +
        " " +
        (message.authorId === -1
          ? "mx-auto"
          : user?.id !== message.authorId
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
  )
}

export default ChatPage
