import { FC, useEffect, useState } from "react"
import rehypeHighlight from "rehype-highlight"
import remarkRehype from "remark-rehype"
import { unified } from "unified"
import remarkParse from "remark-parse"
import rehypeStringify from "rehype-stringify"

import "highlight.js/styles/github-dark.css"

export const RenderMarkdown = async (content: string) => {
  const html = await unified()
    .use(remarkParse)
    .use(remarkRehype)
    .use(rehypeHighlight, { ignoreMissing: true })
    .use(rehypeStringify)
    .process(content)
  return html.toString()
}

const Markdown: FC<{ content: string }> = ({ content }) => {
  const [html, setHtml] = useState("")
  useEffect(() => {
    RenderMarkdown(content).then((html) => {
      setHtml(html)
    })
  }, [content, html])

  return <div className={"prose dark:prose-invert"} dangerouslySetInnerHTML={{ __html: html }} />
}

export default Markdown
