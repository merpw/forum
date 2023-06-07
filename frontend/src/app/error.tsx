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
      <div className={"text-center my-20"}>
        <h1 className={"text-8xl text-base-100 font-Alatsi"}>500</h1>
        <h2 className={"text-2xl font-light"}>Internal Server Error</h2>
        <div className={"my-16 font-light flex flex-col gap-3"}>
          <p className={"font-light flex flex-col gap-1"}>
            {"Please, try to update page or visit "}
            <span className={"my-1"}>
              <Link href={"/"} className={"button"}>
                Back to homepage
              </Link>
            </span>
          </p>
          <p className={"font-light"}>
            {"If error persists, please, "}

            <Link
              href={`mailto:dev@mer.pw?subject=${devMessage}`}
              className={"clickable text-lg font-normal"}
            >
              contact site administrator
            </Link>
          </p>
        </div>
      </div>
    </>
  )
}
