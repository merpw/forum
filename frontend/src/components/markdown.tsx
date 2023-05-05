import { FC, useEffect, useState } from "react"
import rehypeHighlight from "rehype-highlight"
import remarkRehype from "remark-rehype"
import { unified } from "unified"
import remarkParse from "remark-parse"
import rehypeStringify from "rehype-stringify"
import remarkGfm from "remark-gfm"
import { remark } from "remark"
import stripMarkdown, { Handler } from "strip-markdown"

import "highlight.js/styles/github-dark.css"

export const RenderMarkdown = async (content: string) => {
  const html = await unified()
    .use(remarkParse)
    .use(remarkGfm)
    .use(remarkRehype)
    .use(rehypeHighlight, { ignoreMissing: true })
    .use(rehypeStringify)
    .process(content)
  return html.toString()
}

// Just wanted to try out some typescript stuff. If async option is true, it returns a promise, otherwise it returns a string.
export function MarkdownToPlain<
  T extends {
    limit?: number
    removeNewLines?: boolean
    async?: boolean
  }
>(content: string, options: T): T["async"] extends false ? string : Promise<string> {
  let plain = ""

  const applyOptions = () => {
    if (options.limit !== undefined && plain.length > options.limit) {
      plain = plain.slice(0, options.limit) + "..."
    }

    if (options.removeNewLines) {
      plain = plain.replace(/\n+/g, " ").trim()
    }
  }

  const codeHandler: Handler = (node) => {
    return { type: "code", value: `[code in ${node.lang}]` }
  }

  const processor = remark().use(stripMarkdown, {
    remove: [["code", codeHandler]],
  })

  if (options.async !== false) {
    return new Promise<string>(async (resolve) => {
      plain = (await processor.process(content)).toString()
      applyOptions()
      resolve(plain)
    }) as never
  } else {
    plain = processor.processSync(content).toString()
    applyOptions()
    return plain as never
  }
}

const Markdown: FC<{ content: string; className?: string; fallback?: string }> = ({
  content,
  className = "",
  fallback = "",
}) => {
  const [html, setHtml] = useState<string>()

  useEffect(() => {
    RenderMarkdown(content).then((html) => {
      setHtml(html)
    })
  }, [content, html])

  if (html === undefined) {
    return <div className={className}>{fallback}</div>
  }

  return (
    <div
      className={"prose dark:prose-invert max-w-full" + " " + className}
      dangerouslySetInnerHTML={{ __html: html }}
    />
  )
}

export default Markdown
