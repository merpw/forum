import { ReactNode } from "react"
import { Metadata } from "next"
import { Alatsi, Yesteryear } from "next/font/google"

import "@/styles/globals.css"

import Layout from "@/components/layout"

const alatsi = Alatsi({
  subsets: ["latin"],
  variable: "--font-alatsi",
  display: "block",
  weight: "400",
  fallback: ["sans-serif"],
})

const yesteryear = Yesteryear({
  subsets: ["latin"],
  variable: "--font-yesteryear",
  display: "block",
  weight: "400",
  fallback: ["cursive"],
})

export const metadata: Metadata = {
  title: {
    default: "Forum",
    template: "%s - Forum",
  },
  description: "The friendliest forum",
  metadataBase: process.env.FORUM_BASE ? new URL(process.env.FORUM_BASE) : null,
  openGraph: {
    siteName: "Forum",
  },
  icons: [
    { rel: "icon", url: "/faviconDark.ico" },
    { rel: "icon", url: "/favicon.ico", media: "(prefers-color-scheme: light)" },
  ],
}

// TODO: test if all the necessary features are supported with the Edge runtime
export const runtime = "edge"

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang={"en"} className={alatsi.variable + " " + yesteryear.variable}>
      <body>
        <Layout>{children}</Layout>
      </body>
    </html>
  )
}
