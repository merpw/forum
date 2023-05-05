import "../styles/globals.css"
import type { AppProps } from "next/app"
import React, { ReactElement, ReactNode } from "react"
import { NextSeo } from "next-seo"
import { NextPage } from "next"

import Layout from "@/components/layout"

export type NextPageWithLayout<P = object, IP = P> = NextPage<P, IP> & {
  getLayout?: (page: ReactElement) => ReactNode
}

type AppPropsWithLayout = AppProps & {
  Component: NextPageWithLayout
}

export default function App({ Component, pageProps }: AppPropsWithLayout) {
  const getLayout = Component.getLayout ?? ((page) => <Layout>{page}</Layout>)

  return getLayout(
    <>
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
    </>
  )
}
