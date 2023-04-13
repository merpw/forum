import Link from "next/link"
import { useEffect, useState } from "react"
import { NextSeo } from "next-seo"

export default function Custom500() {
  const [devMessage, setDevMessage] = useState("")
  useEffect(() => {
    setDevMessage(`Error 500 on ${document.URL}`)
  }, [])
  return (
    <>
      <NextSeo title={"Oops... Error 500"} />
      <div className={"text-center"}>
        <h1 className={"text-5xl m-auto"}>500 Internal Server Error</h1>
        <p className={"text-2xl my-5 font-light"}>
          {"Please, try to update page or visit "}
          <Link
            href={"/"}
            className={
              "font-normal text-white bg-blue-600 rounded-md pt-1.5 pb-2 px-3 hover:opacity-80"
            }
          >
            Back to homepage
          </Link>
        </p>
        <p className={"text-2xl my-5 font-light"}>
          {"If error persists, please, "}
          <Link
            href={`mailto:dev@mer.pw?subject=${devMessage}`}
            className={"font-normal hover:underline text-bo"}
          >
            contact site administrator
          </Link>
        </p>
      </div>
    </>
  )
}
