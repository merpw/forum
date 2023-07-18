import { FC, useCallback, useEffect, useMemo, useRef, useState } from "react"
import ReactTextAreaAutosize from "react-textarea-autosize"
import { useDispatch } from "react-redux"

import throttle from "@/helpers/throttle"
import wsActions from "@/store/ws/actions"
import { useAppSelector } from "@/store/hooks"

const WriteMessageForm: FC<{
  sendMessage: (content: string) => void
}> = ({ sendMessage }) => {
  const [input, setInput] = useState("")

  const chatId = useAppSelector((state) => state.chats.activeChatId)

  const formRef = useRef<HTMLFormElement>(null)
  const inputRef = useRef<HTMLTextAreaElement>(null)

  const dispatch = useDispatch()

  const sendTypingStatus = useCallback(
    (isTyping: boolean) => {
      dispatch(
        wsActions.sendPost(`/chat/${chatId}/typing`, {
          isTyping,
        })
      )
    },
    [chatId, dispatch]
  )

  const sendTyping = useMemo(
    () =>
      throttle(() => sendTypingStatus(true), 2000, {
        leading: true,
        trailing: false,
      }),
    [sendTypingStatus]
  )

  useEffect(() => {
    inputRef.current?.focus()
  }, [chatId])

  useEffect(() => {
    const listener = () => {
      sendTypingStatus(false)
    }
    window.addEventListener("blur", listener)
    return () => window.removeEventListener("blur", listener)
  }, [sendTypingStatus])

  return (
    <form
      ref={formRef}
      className={"my-1 relative"}
      onSubmit={async (e) => {
        e.preventDefault()
        sendMessage(input)
        setInput("")
        sendTypingStatus(false)
      }}
      onBlur={() => {
        sendTypingStatus(false)
      }}
    >
      <ReactTextAreaAutosize
        ref={inputRef}
        className={"input input-bordered border-primary bg-base-100 p-2 pr-10 w-full"}
        onChange={(e) => {
          setInput(e.target.value)
          sendTyping()
        }}
        onBlur={() => setInput(input.trim())}
        value={input}
        maxRows={5}
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
      <button
        className={"absolute z-10 clickable disabled:opacity-50 right-2.5 bottom-3.5"}
        type={"submit"}
      >
        <svg
          xmlns={"http://www.w3.org/2000/svg"}
          fill={"none"}
          viewBox={"0 0 24 24"}
          strokeWidth={2}
          stroke={"currentColor"}
          className={"w-7 h-7 text-primary"}
        >
          <path
            strokeLinecap={"round"}
            strokeLinejoin={"round"}
            d={
              "M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"
            }
          />
        </svg>
      </button>
    </form>
  )
}

export default WriteMessageForm
