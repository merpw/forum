"use client"

import { FC, ReactNode } from "react"
import Link from "next/link"

import { useMe } from "@/api/auth/hooks"

const UserLink: FC<{ userId: number; children?: ReactNode; className?: string }> = ({
  userId,
  children,
  className,
}) => {
  const { user: meUser } = useMe()

  return (
    <Link href={meUser && meUser.id === userId ? "/me" : `/user/${userId}`} className={className}>
      {children}
    </Link>
  )
}

export default UserLink
