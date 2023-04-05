import "../styles/globals.css"
import type { AppProps } from "next/app"
import React from "react"
import Layout from "../components/layout"
import { NextSeo } from "next-seo"

export default function App({ Component, pageProps }: AppProps) {
  return (
    <Layout>
      <NextSeo
        titleTemplate={"%s - Forum"}
        description={"The friendliest forum"}
        openGraph={{ siteName: "Forum" }}
        additionalLinkTags={[
          { rel: "icon", href: "/faviconDark.ico" },
          { rel: "icon", href: "/favicon.ico", media: "(prefers-color-scheme: light)" },
        ]}
      />
      <Component {...pageProps} />
    </Layout>
  )
}
