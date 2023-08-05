"use client"

import { useInvitations } from "@/api/invitations/hooks"
import InvitationCard from "@/components/invitations/Card"

const NavbarBell = () => {
  const { invitations, mutate } = useInvitations()

  return (
    <div className={"dropdown dropdown-end min-w-fit"}>
      <div className={"relative"}>
        <button className={"btn btn-ghost btn-circle"} onClick={() => mutate()}>
          <svg
            xmlns={"http://www.w3.org/2000/svg"}
            fill={"none"}
            viewBox={"0 0 24 24"}
            strokeWidth={1.5}
            stroke={"currentColor"}
            className={"w-full h-full p-2.5"}
          >
            <path
              strokeLinecap={"round"}
              strokeLinejoin={"round"}
              d={
                "M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75v-.7V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0"
              }
            />
          </svg>
        </button>
        {invitations && (
          <span className={"badge badge-sm badge-accent absolute top-0 right-0"}>
            {invitations.length}
          </span>
        )}
      </div>
      <ul
        tabIndex={0}
        className={
          "mt-3 p-2 shadow z-50 menu menu-compact dropdown-content bg-base-100 rounded-box"
        }
      >
        <li className={"menu-title inline"}>
          <span className={"font-light"}>Notifications</span>
        </li>
        <hr
          className={
            "mx-3 mb-1 cursor-default border-dotted border-t-0 border-b-4 border-info opacity-20"
          }
        />
        {invitations?.map((invitation) => (
          <li key={invitation}>
            <InvitationCard id={invitation} />
          </li>
        ))}
        {invitations?.length === 0 && (
          <li className={"menu-title w-44"}>
            <div className={"text-base-content"}>
              <span className={"font-light"}>No notifications yet</span>
            </div>
          </li>
        )}
      </ul>
    </div>
  )
}

export default NavbarBell
