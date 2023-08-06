import { MdFormatBold, MdFormatItalic, MdFormatQuote, MdLink } from "react-icons/md"
import { IconType } from "react-icons"
import { FC, RefObject } from "react"

import { wrapWith } from "@/components/markdown/editor/helpers"

type Button = {
  icon: IconType
  text: string
  onClick: (textAreaRef: HTMLTextAreaElement) => void
}

const buttons: Button[] = [
  {
    icon: MdFormatBold,
    text: "Bold",
    onClick: (textAreaRef) => {
      wrapWith(textAreaRef, "**", "**")
    },
  },
  {
    icon: MdFormatItalic,
    text: "Italic",
    onClick: (textAreaRef) => {
      wrapWith(textAreaRef, "*", "*")
    },
  },
  {
    icon: MdFormatQuote,
    text: "Quote",
    onClick: (textAreaRef) => {
      wrapWith(textAreaRef, "> ", "")
    },
  },
  {
    icon: MdLink,
    text: "Link",
    onClick: (textAreaRef) => {
      wrapWith(textAreaRef, "[", "](url)")
    },
  },
]

const ToolBar: FC<{ textareaRef: RefObject<HTMLTextAreaElement> }> = ({ textareaRef }) => {
  return (
    <ul className={"ml-auto flex"}>
      {buttons.map((button, i) => (
        <li key={i}>
          <button
            type={"button"}
            tabIndex={-1} // TODO: improve accessibility
            className={"btn btn-square btn-ghost btn-sm"}
            onClick={() => textareaRef.current && button.onClick(textareaRef.current)}
            title={button.text}
          >
            <button.icon size={20} />
          </button>
        </li>
      ))}
    </ul>
  )
}

export default ToolBar
