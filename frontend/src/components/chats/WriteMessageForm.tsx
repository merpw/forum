import { FC, useEffect, useRef, useState } from "react"
import ReactTextAreaAutosize from "react-textarea-autosize"

const WriteMessageForm: FC<{
  sendMessage: (content: string) => void
}> = ({ sendMessage }) => {
  const [input, setInput] = useState("")

  const formRef = useRef<HTMLFormElement>(null)
  const inputRef = useRef<HTMLTextAreaElement>(null)

  useEffect(() => {
    inputRef.current?.focus()
  }, [sendMessage])

  return (
    <form
      ref={formRef}
      className={"mb-6"}
      onSubmit={async (e) => {
        e.preventDefault()
        sendMessage(input)
        setInput("")
      }}
    >
      <ReactTextAreaAutosize
        ref={inputRef}
        className={"textarea w-full my-3"}
        onChange={(e) => setInput(e.currentTarget.value)}
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
      <button className={"button disabled:opacity-50"} type={"submit"}>
        Send
      </button>
    </form>
  )
}

export default WriteMessageForm
