import { ReactNode } from "react"
import Link from "next/link"

import Navbar from "./navbar"

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <div className={"container m-auto p-5 min-h-screen break-words"}>
      <header>
        <Navbar />
      </header>
      <main>{children}</main>
      <footer>
        <div
          className={
            "text-center text-sm font-light text-neutral border-t border-base-100 pt-3 mt-3"
          }
        >
          <p>
            {"Made with "}
            <Link href={"https://github.com/merpw/forum"} className={"gradient-text clickable"}>
              {"❤"}
            </Link>
            {" by "}
            <span className={"font-light"}>{"basket of kittens"}</span>
          </p>
        </div>
      </footer>
    </div>
  )
}
