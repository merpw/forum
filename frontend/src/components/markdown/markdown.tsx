"use client"

import { FC, useEffect, useState } from "react"

import "highlight.js/styles/github-dark.css"
import RenderMarkdown from "@/components/markdown/render"

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
