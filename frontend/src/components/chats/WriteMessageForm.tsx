import { FC, useCallback, useEffect, useMemo, useRef } from "react"
import { useDispatch } from "react-redux"

import throttle from "@/helpers/throttle"
import wsActions from "@/store/ws/actions"
import { useAppSelector } from "@/store/hooks"
import MarkdownEditor from "@/components/markdown/editor/Editor"
import { trimInput } from "@/helpers/input"

const WriteMessageForm: FC<{
  sendMessage: (content: string) => void
}> = ({ sendMessage }) => {
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
      className={"my-1 bg-base-200 p-2 rounded"}
      onSubmit={async (e) => {
        e.preventDefault()

        console.log(inputRef.current)

        const inputMessage = inputRef.current as HTMLTextAreaElement
        const trimmedInput = inputMessage.value.trim()
        if (!trimmedInput) return

        inputMessage.value = ""

        sendMessage(trimmedInput)
        sendTypingStatus(false)
      }}
      onBlur={() => {
        sendTypingStatus(false)
      }}
    >
      <MarkdownEditor
        name={"message"}
        withSubmit
        ref={inputRef}
        className={"input input-bordered w-full"}
        onChange={sendTyping}
        onBlur={trimInput}
        maxRows={5}
        required
      />
    </form>
  )
}

export default WriteMessageForm
