import { FC, lazy, RefObject, Suspense, useEffect, useRef, useState } from "react"
import ReactTextAreaAutosize, { TextareaAutosizeProps } from "react-textarea-autosize"

import ToolBar from "@/components/markdown/editor/ToolBar"
import uploadFile from "@/api/attachments/upload"

const ImageTypeRegex = /^image\/(png|jpe?g|gif|webp)$/

const MaxImageSize = 1024 * 1024 * 20 // 20 MB

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

  const [uploadError, setUploadError] = useState<string | null>(null)

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
        onChange={(e) => {
          props.onChange?.(e)
          setUploadError(null)
        }}
        className={isPreview ? "hidden" : props.className}
        onDrop={(e) => {
          console.log(e)
        }}
        onPaste={(e) => {
          if (e.clipboardData.files.length > 0) {
            e.preventDefault()
            const file = e.clipboardData.files[0]

            if (file.size > MaxImageSize) {
              setUploadError("File too large")
              return
            }

            if (!ImageTypeRegex.test(file.type)) {
              setUploadError("Invalid file type, only images are allowed")
              return
            }

            uploadFile(file)
              .then((url) => {
                document.execCommand("insertText", false, `![${file.name}](${url})`)
              })
              .catch((e) => {
                console.error(e)
                setUploadError("Failed to upload image")
              })
          }
        }}
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
      {uploadError && (
        <div className={"bg-error text-error-content rounded p-1"}>{uploadError}</div>
      )}
    </div>
  )
}

export default MarkdownEditor
