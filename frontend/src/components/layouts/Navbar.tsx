import { NextComponentType } from "next"
import Link from "next/link"

import NavUserInfo from "@/components/layouts/NavUserInfo"

const Navbar: NextComponentType = () => {
  return (
    <div className={"border-b border-base-100 pb-3 mb-3"}>
      <nav className={"flex justify-between flex-wrap"}>
        <div className={"self-center"}>
          <Link href={"/"} className={"clickable text-4xl font-Alatsi"}>
            <span className={"font-black  gradient-text"}>{"{"}</span>
            <span className={"text-3xl"}>{"FORUM"}</span>
            <span className={"font-black  gradient-text"}>{"}"}</span>
          </Link>
        </div>
        <div className={"my-auto"}>
          <NavUserInfo />
        </div>
      </nav>
    </div>
  )
}

export default Navbar
