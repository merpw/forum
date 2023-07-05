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
      className={"my-1 relative"}
      onSubmit={async (e) => {
        e.preventDefault()
        sendMessage(input)
        setInput("")
      }}
    >
      <ReactTextAreaAutosize
        ref={inputRef}
        className={"absolute input input-bordered border-primary bg-base-100 p-2 w-full"}
        onChange={(e) => setInput(e.currentTarget.value)}
        onBlur={() => setInput(input.trim())}
        value={input}
        maxRows={4}
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
        className={"absolute z-10 clickable disabled:opacity-50 right-1.5 m-1.5"}
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
