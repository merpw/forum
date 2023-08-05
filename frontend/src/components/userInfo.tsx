"use client"

import Link from "next/link"

import { useLogOut, useMe } from "@/api/auth/hooks"
import Avatar from "@/components/Avatar"

const UserInfo = () => {
  const { user, isError, isLoading, isLoggedIn } = useMe()
  const logOut = useLogOut()

  if (isError || isLoading) {
    return null
  }

  if (!isLoggedIn || !user) {
    return (
      <Link className={"clickable flex font-light space-x-0.5"} href={"/login"}>
        <span>Login</span>
        <svg
          xmlns={"http://www.w3.org/2000/svg"}
          fill={"none"}
          viewBox={"0 0 24 24"}
          strokeWidth={2.5}
          stroke={"currentColor"}
          className={"w-6 h-6 text-primary"}
        >
          <title>Login</title>
          <path
            strokeLinecap={"round"}
            strokeLinejoin={"round"}
            d={
              "M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75"
            }
          />
        </svg>
      </Link>
    )
  }
  return (
    <div className={"flex gap-1 my-auto"}>
      <div className={"dropdown dropdown-end"}>
        <div tabIndex={0} className={"btn min-w-full btn-ghost btn-circle"}>
          <Avatar user={user} size={40} className={"w-9"} />
        </div>

        <ul
          tabIndex={0}
          className={
            "mt-3 p-2 shadow z-50 menu menu-compact dropdown-content bg-base-100 rounded-box w-52"
          }
        >
          <li className={"menu-title inline"}>
            <span className={"font-light"}>Hello, </span>
            <span className={"text-primary"}>{user.username}</span>
          </li>
          <hr className={"mx-3 mb-1 border-dotted border-t-0 border-b-4 border-info opacity-20"} />
          <li>
            <Link href={"/me"}>Profile</Link>
          </li>
          <li>
            <button type={"submit"} title={"Logout"} onClick={() => logOut()}>
              Logout
            </button>
          </li>
        </ul>
      </div>
    </div>
  )
}

export default UserInfo
