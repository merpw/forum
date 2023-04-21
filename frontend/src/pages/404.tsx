import Link from "next/link"
import { NextSeo } from "next-seo"

export default function Custom404() {
  return (
    <>
      <NextSeo title={"Not Found"} />
      <div className={"text-center my-20"}>
        <div className={"text-5xl m-auto"}>
          <h1 className={"text-8xl opacity-80"}>404</h1>
          <h2 className={"text-2xl"}>Not Found</h2>
        </div>
        <p className={"text-xl my-16 font-light"}>
          <Link href={"/"} className={"button"}>
            Back to homepage
          </Link>
        </p>
      </div>
    </>
  )
}
