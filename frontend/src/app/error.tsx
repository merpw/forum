"use client"

import Link from "next/link"
import { useEffect, useState } from "react"

// TODO: add error-global to handle layout errors
// TODO: add automatic error reporting

export default function Custom500({ reset }: { error: Error; reset: () => void }) {
  const [devMessage, setDevMessage] = useState("")
  useEffect(() => {
    setDevMessage(`Error 500 on ${document.URL}`)
  }, [])
  return (
    <>
      <title>Oops... Error 500 - Forum</title>
      <div className={"text-center"}>
        <h1 className={"text-5xl m-auto"}>500 Internal Server Error</h1>
        <p className={"text-2xl my-5 font-light"}>
          {"Please, try again later or go to the Home page"}
        </p>
        <div className={"flex gap-3 justify-center"}>
          <button className={"text-2xl button"} onClick={reset}>
            Try again
          </button>
          <Link href={"/"}>
            <button className={"text-2xl button"}>To the home page</button>
          </Link>
        </div>

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
