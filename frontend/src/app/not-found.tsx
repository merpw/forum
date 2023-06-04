import Link from "next/link"

// unexpected error with edge runtime, the page is static
export const runtime = "nodejs"

export default function Custom404() {
  return (
    <>
      <title>Not Found - Forum</title>
      <div className={"text-center my-20"}>
        <h1 className={"text-8xl text-base-100 font-Alatsi"}>404</h1>
        <h2 className={"text-2xl font-extralight"}>Not Found</h2>
        <p className={"text-xl my-16 font-light"}>
          <Link href={"/"} className={"button"}>
            Back to homepage
          </Link>
        </p>
      </div>
    </>
  )
}
