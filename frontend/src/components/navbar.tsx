import { NextComponentType } from "next"
import Link from "next/link"

import UserInfo from "@/components/userInfo"

// TODO: fix hydration error

const Navbar: NextComponentType = () => {
  return (
    <header
      className={
        "sticky top-0 z-50 border-b border-base-100 py-2 px-3 bg-base-100 dark:bg-base-300"
      }
    >
      <nav className={"container m-auto flex justify-between flex-wrap"}>
        <div className={"self-center"}>
          <Link href={"/"} className={"clickable text-4xl font-Alatsi"}>
            <span className={"font-black  gradient-text"}>{"{"}</span>
            <span className={"text-3xl"}>{"FORUM"}</span>
            <span className={"font-black  gradient-text"}>{"}"}</span>
          </Link>
        </div>
        <div className={"my-auto"}>
          <UserInfo />
        </div>
      </nav>
    </header>
  )
}

export default Navbar
