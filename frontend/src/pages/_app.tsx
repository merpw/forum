import "../styles/globals.css"
import type { AppProps } from "next/app"
import React from "react"
import Layout from "../components/layout"
import Head from "next/head"

export default function App({ Component, pageProps }: AppProps) {
  return (
    <Layout>
      <Head>
        <meta property={"og:site_name"} content={"FORUM"} />
        <meta property={"og:type"} content={"website"} />
        <meta property={"og:url"} content={"https://forum.mer.pw"} />
      </Head>
      <Component {...pageProps} />
    </Layout>
  )
}
