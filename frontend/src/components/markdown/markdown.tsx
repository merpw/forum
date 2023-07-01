"use client"

import { FC, useEffect, useState } from "react"

import "highlight.js/styles/github-dark.css"

import RenderMarkdown from "@/components/markdown/render"

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
