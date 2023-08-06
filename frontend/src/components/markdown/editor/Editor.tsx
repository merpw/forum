import { FC, lazy, RefObject, Suspense, useEffect, useRef, useState } from "react"
import ReactTextAreaAutosize, { TextareaAutosizeProps } from "react-textarea-autosize"

import ToolBar from "@/components/markdown/editor/ToolBar"

const Markdown = lazy(() => import("@/components/markdown/markdown"))

const MarkdownEditor: FC<
  TextareaAutosizeProps & {
    ref?: RefObject<HTMLTextAreaElement>
  }
> = (props) => {
  const [isPreview, setIsPreview] = useState(false)

  const fallbackRef = useRef<HTMLTextAreaElement>(null)

  const ref = props.ref || fallbackRef

  const selectionBackup = useRef<{ selectionStart: number; selectionEnd: number } | null>(null)

  useEffect(() => {
    const textarea = ref.current as HTMLTextAreaElement
    if (isPreview) {
      selectionBackup.current = {
        selectionStart: textarea.selectionStart,
        selectionEnd: textarea.selectionEnd,
      }
    } else {
      if (selectionBackup.current != null) {
        textarea.selectionStart = selectionBackup.current.selectionStart
        textarea.selectionEnd = selectionBackup.current.selectionEnd
        textarea.focus()
      }
    }
  }, [isPreview, ref])

  return (
    <div className={"bg-base-200 rounded"}>
      <div className={"tabs"}>
        <span
          className={"tab tab-lifted" + " " + (isPreview ? "" : "tab-active")}
          onClick={() => setIsPreview(false)}
        >
          Edit
        </span>
        <span
          className={"tab tab-lifted" + " " + (isPreview ? "tab-active" : "")}
          onClick={() => setIsPreview(true)}
        >
          Preview
        </span>
        {!isPreview && <ToolBar textareaRef={ref} />}
      </div>
      <ReactTextAreaAutosize
        contentEditable={true}
        {...props}
        ref={ref}
        className={isPreview ? "hidden" : props.className}
      />
      {isPreview && (
        <div className={"min-h-16"}>
          <Suspense fallback={<div className={"prose bg-base-100 p-5 rounded-b"}>Loading...</div>}>
            <Markdown
              content={ref.current?.value || "Nothing to preview"}
              className={"bg-base-100 p-5 rounded-b"}
            />
          </Suspense>
        </div>
      )}
    </div>
  )
}

export default MarkdownEditor
