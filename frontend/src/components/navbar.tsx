import { NextComponentType } from "next"
import Link from "next/link"

import UserInfo from "@/components/userInfo"

const Navbar: NextComponentType = () => {
  return (
    <div className={"border-b mb-5 pb-2"}>
      <nav className={"flex justify-between flex-wrap"}>
        <div>
          <Link href={"/"} className={"clickable text-3xl"}>
            FORUM
          </Link>
        </div>
        <div className={"ml-auto"}>
          <UserInfo />
        </div>
      </nav>
    </div>
  )
}

export default Navbar
