"use client"

import { FC } from "react"

import { useIsUserOnline } from "@/api/users/hooks"
import AvatarPlaceholder from "@/components/placeholders/avatar"

const Avatar: FC<{ userId: number; className?: string }> = ({ userId, className }) => {
  const isOnline = useIsUserOnline(userId)

  return (
    <>
      <div
        className={
          "font-Alatsi text-xl text-primary rounded-full ring-2" +
          " " +
          (isOnline ? "ring-accent" : "ring-neutral") +
          " " +
          className
        }
      >
        <AvatarPlaceholder />
        {/* TODO: add Avatar*/}
      </div>
    </>
  )
}

export default Avatar
