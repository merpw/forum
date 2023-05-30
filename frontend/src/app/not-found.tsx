import Link from "next/link"

// unexpected error with edge runtime, the page is static
export const runtime = "nodejs"

export default function Custom404() {
  return (
    <>
      <title>Not Found - Forum</title>
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
