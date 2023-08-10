import {
  forwardRef,
  ForwardRefRenderFunction,
  lazy,
  RefObject,
  Suspense,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react"
import ReactTextAreaAutosize, { TextareaAutosizeProps } from "react-textarea-autosize"
import { useDropzone } from "react-dropzone"

import ToolBar from "@/components/markdown/editor/ToolBar"
import uploadFile from "@/api/attachments/upload"

const ImageTypeRegex = /^image\/(png|jpe?g|gif|webp)$/

const MaxImageSize = 1024 * 1024 * 20 // 20 MB

const Markdown = lazy(() => import("@/components/markdown/markdown"))

const MarkdownEditor: ForwardRefRenderFunction<
  HTMLTextAreaElement,
  TextareaAutosizeProps & {
    withSubmit?: boolean
  }
> = (props, forwardedRef) => {
  const [isPreview, setIsPreview] = useState(false)

  const fallbackRef = useRef<HTMLTextAreaElement>(null)

  const ref = (forwardedRef || fallbackRef) as RefObject<HTMLTextAreaElement>

  const selectionBackup = useRef<{ selectionStart: number; selectionEnd: number } | null>(null)

  const [uploadError, setUploadError] = useState<string | null>(null)

  useEffect(() => {
    if (ref.current?.value === "" && isPreview) {
      setIsPreview(false)
    }
  }, [isPreview, ref, ref.current?.value])

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

  const innerProps = useMemo(() => {
    const propsCopy = { ...props }
    delete propsCopy.withSubmit
    return propsCopy
  }, [props])

  return (
    <div
      className={
        "relative bg-base-200 rounded" +
        " " +
        (isDragActive ? "ring ring-opacity-60 ring-primary" : "")
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
          onClick={() => {
            if (ref.current?.value !== "") {
              setIsPreview(true)
            }
          }}
        >
          Preview
        </span>
        {!isPreview && <ToolBar textareaRef={ref} />}
      </div>
      <ReactTextAreaAutosize
        {...innerProps}
        ref={ref}
        value={props.value}
        onChange={(e) => {
          props.onChange?.(e)
          setUploadError(null)
        }}
        className={
          isPreview ? "hidden" : props.className + " " + (props.withSubmit ? "py-3 px-3 pr-10" : "")
        }
        onPaste={(e) => {
          if (e.clipboardData.files.length > 0) {
            e.preventDefault()
            uploadImage(e.clipboardData.files[0])
          }
        }}
        onKeyDown={
          props.withSubmit
            ? (e) => {
                if (e.key === "Enter" && !e.shiftKey) {
                  e.preventDefault()

                  // trim input
                  ref.current?.blur()
                  ref.current?.focus()

                  ref.current?.form?.dispatchEvent(
                    new Event("submit", { cancelable: true, bubbles: true })
                  )
                }
              }
            : undefined
        }
      />
      {isPreview ? (
        <div className={"min-h-[5.15rem]"}>
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
      {props.withSubmit && (
        <button
          className={
            "absolute z-10 clickable disabled:opacity-50 right-3 bottom-11" +
            " " +
            (ref.current?.value === "" ? "btn-disabled opacity-50" : "")
          }
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
      )}
    </div>
  )
}

export default forwardRef(MarkdownEditor)
