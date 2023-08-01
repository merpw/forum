import { FC } from "react"
import Image from "next/image"

import AvatarPlaceholder from "@/components/placeholders/avatar"

const AVATARS = [
  "0.jpg",
  "1.jpg",
  "2.jpg",
  "3.jpg",
  "4.jpg",
  "5.jpg",
  "6.jpg",
  "7.jpg",
  "8.jpg",
  "9.jpg",
]

const ChooseAvatar: FC<{
  avatar: string | undefined
  setAvatar: (newAvatar: string | undefined) => void
}> = ({ avatar, setAvatar }) => {
  return (
    <div className={"mx-auto"}>
      {avatar ? (
        <Image
          src={`/avatars/${avatar}`}
          alt={"avatar"}
          className={"w-52 mx-auto rounded-full"}
          width={96}
          height={96}
        />
      ) : (
        <AvatarPlaceholder className={"w-52 mx-auto"} />
      )}
      <div className={"flex gap-2 mt-2 p-1 overflow-auto bg-base-200"}>
        {[undefined, ...AVATARS].map((avatarName) => {
          const props = {
            className:
              "cursor-pointer" +
              " " +
              (avatar === avatarName ? "ring ring-primary" : "brightness-75"),
            onClick: () => setAvatar(avatarName),
          }

          const size = 60

          return avatarName ? (
            <Image
              {...props}
              src={`/avatars/${avatarName}`}
              alt={"avatar"}
              width={size}
              height={size}
              key={avatarName}
            />
          ) : (
            <span {...props} key={avatarName}>
              <AvatarPlaceholder width={size} height={size} />
            </span>
          )
        })}
      </div>
    </div>
  )
}

export default ChooseAvatar
