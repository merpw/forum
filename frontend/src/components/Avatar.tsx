"use client"

import { FC } from "react"
import Image from "next/image"

import { useIsUserOnline } from "@/api/users/hooks"
import AvatarPlaceholder from "@/components/placeholders/avatar"
import { useMe } from "@/api/auth/hooks"
import { User } from "@/custom"

const Avatar: FC<{ user: User; className?: string; size: number }> = ({
  user,
  className,
  size,
}) => {
  const isOnline = useIsUserOnline(user.id)

  const { user: reqUser } = useMe()

  return (
    <>
      <div
        className={
          "rounded-full ring-2" +
          " " +
          (isOnline ? "ring-accent" : "ring-neutral") +
          " " +
          className
        }
      >
        {user?.avatar ? (
          <Image
            className={"rounded-full w-full h-full p-0.5"}
            src={`/avatars/${user.avatar}`}
            alt={user.username === reqUser?.username ? "Your avatar" : `${user.username}'s avatar`}
            width={size}
            height={size}
          />
        ) : (
          <AvatarPlaceholder className={"w-full h-full"} />
        )}
      </div>
    </>
  )
}

export default Avatar
