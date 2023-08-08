import { FC, lazy, RefObject, Suspense, useCallback, useEffect, useRef, useState } from "react"
import ReactTextAreaAutosize, { TextareaAutosizeProps } from "react-textarea-autosize"
import { useDropzone } from "react-dropzone"

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

  const uploadImage = useCallback(
    async (file: File) => {
      setUploadError(null)
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
          ref.current?.focus()
          document.execCommand("insertText", false, `![${file.name}](${url})`)
        })
        .catch((e) => {
          console.error(e)
          setUploadError("Failed to upload image")
        })
    },
    [ref]
  )

  const {
    getRootProps,
    getInputProps,
    isDragActive,
    open: openFileUpload,
  } = useDropzone({
    accept: {
      "image/png": [".png"],
      "image/jpeg": [".jpg", ".jpeg"],
      "image/gif": [".gif"],
      "image/webp": [".webp"],
    },
    onDrop: (acceptedFiles) => {
      if (acceptedFiles.length > 0) {
        uploadImage(acceptedFiles[0])
      } else {
        setUploadError("Invalid file type, only images are allowed")
      }
    },
    noClick: true,
  })

  return (
    <div
      className={
        "bg-base-200 rounded" + " " + (isDragActive ? "ring ring-opacity-60 ring-primary" : "")
      }
      {...getRootProps()}
      tabIndex={-1}
    >
      <input {...getInputProps()} />
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
        className={isPreview ? "hidden" : props.className + " " + ""}
        onPaste={(e) => {
          if (e.clipboardData.files.length > 0) {
            e.preventDefault()
            uploadImage(e.clipboardData.files[0])
          }
        }}
      />
      {isPreview ? (
        <div className={"min-h-16"}>
          <Suspense fallback={<div className={"prose bg-base-100 p-5 rounded-b"}>Loading...</div>}>
            <Markdown
              content={ref.current?.value || "Nothing to preview"}
              className={"bg-base-100 p-5 rounded-b"}
            />
          </Suspense>
        </div>
      ) : (
        <>
          <button
            type={"button"}
            className={
              "text-sm  px-2 pb-1 rounded-md" + " " + (uploadError ? "text-error" : "text-info")
            }
            onClick={() => {
              setUploadError(null)
              openFileUpload()
            }}
          >
            {isDragActive
              ? "Drop the files here ..."
              : uploadError
              ? uploadError
              : "Attach images by dragging & dropping, selecting or pasting them."}
          </button>
        </>
      )}
    </div>
  )
}

export default MarkdownEditor
