import { Head, Html, Main, NextScript } from "next/document"

export default function Document() {
  return (
    <Html lang={"en"}>
      <Head>
        <link href={"/favicon.ico"} rel={"icon"} media={"(prefers-color-scheme: light)"} />
        <link href={"/faviconDark.png"} rel={"icon"} media={"(prefers-color-scheme: dark)"} />
      </Head>
      <body>
        <Main />
        <NextScript />
      </body>
    </Html>
  )
}
