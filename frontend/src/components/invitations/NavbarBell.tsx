"use client"

import { useEffect, useRef } from "react"

import { useInvitations } from "@/api/invitations/hooks"
import InvitationCard from "@/components/invitations/Card"

// TODO: improve accessibility

const NavbarBell = () => {
  const { invitations, mutate } = useInvitations()

  const detailsRef = useRef<HTMLDetailsElement>(null)

  useEffect(() => {
    if (detailsRef.current?.open) {
      const summaryElement = detailsRef.current.firstChild as HTMLElement
      summaryElement.focus()
    }
  }, [invitations])

  return (
    <details
      className={"dropdown dropdown-end min-w-fit"}
      onBlur={() => {
        const detailsElement = detailsRef.current as HTMLDetailsElement

        requestAnimationFrame(() => {
          if (!document.activeElement || !detailsElement.contains(document.activeElement)) {
            // close only if focuses outside
            detailsElement.open = false
          }
        })
      }}
      ref={detailsRef}
    >
      <summary className={"btn btn-ghost btn-circle relative"} onClick={() => mutate()}>
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
        {invitations && invitations.length > 0 && (
          <span className={"badge badge-sm badge-accent absolute top-0 right-0"}>
            {invitations.length}
          </span>
        )}
      </summary>

      <ul
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
    </details>
  )
}

export default NavbarBell
