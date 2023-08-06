import { FC, lazy, RefObject, Suspense, useRef, useState } from "react"
import ReactTextAreaAutosize, { TextareaAutosizeProps } from "react-textarea-autosize"

const Markdown = lazy(() => import("@/components/markdown/markdown"))

const MarkdownEditor: FC<
  TextareaAutosizeProps & {
    ref?: RefObject<HTMLTextAreaElement>
  }
> = (props) => {
  const [isPreview, setIsPreview] = useState(false)

  const fallbackRef = useRef(null)

  const ref = props.ref || fallbackRef

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
      </div>
      <ReactTextAreaAutosize
        {...props}
        ref={ref}
        className={isPreview ? "hidden" : props.className}
      />
      {isPreview && (
        <div className={"min-h-16"}>
          <Suspense fallback={<div className={"prose bg-base-100 p-5 rounded-b"}>Loading...</div>}>
            <Markdown
              content={props.value || "Nothing to preview"}
              className={"bg-base-100 p-5 rounded-b"}
            />
          </Suspense>
        </div>
      )}
    </div>
  )
}

export default MarkdownEditor
