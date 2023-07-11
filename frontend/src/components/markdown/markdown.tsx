"use client"

import { FC, memo, useEffect, useState } from "react"

import "highlight.js/styles/github-dark.css"

import RenderMarkdown from "@/components/markdown/render"

const Markdown: FC<{ content: string; className?: string; fallback?: string }> = ({ content }) => {
  const [html, setHtml] = useState<string>()

  useEffect(() => {
    RenderMarkdown(content).then((html) => {
      setHtml(html)
    })
  }, [content, html])

  if (!html) {
    return null
  }

  return (
    <div
      className={"prose dark:prose-invert max-w-full"}
      dangerouslySetInnerHTML={{ __html: html }}
    />
  )
}

export default memo(Markdown)
