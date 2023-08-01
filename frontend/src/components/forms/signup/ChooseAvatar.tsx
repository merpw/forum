import { FC, useState } from "react"
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

const ChooseAvatar: FC = () => {
  const [avatar, setAvatar] = useState<string | undefined>(undefined)

  return (
    <div className={"mx-auto"}>
      {avatar && <input type={"hidden"} name={"avatar"} value={avatar} />}
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
        {[undefined, ...AVATARS].map((avatarName, key) => {
          const size = 60

          return (
            <button
              key={key}
              type={"button"}
              onClick={() => setAvatar(avatarName)}
              title={avatarName ? `Avatar option ${key} of ${AVATARS.length}` : "No avatar"}
              className={
                "cursor-pointer min-w-fit" +
                " " +
                (avatar === avatarName ? "ring ring-primary" : "brightness-75")
              }
            >
              {avatarName ? (
                <Image
                  src={`/avatars/${avatarName}`}
                  alt={"Avatar option " + key}
                  width={size}
                  height={size}
                  key={avatarName}
                  // ignore this in accessibility
                  tabIndex={-1}
                />
              ) : (
                <AvatarPlaceholder width={size} height={size} />
              )}
            </button>
          )
        })}
      </div>
    </div>
  )
}

export default ChooseAvatar
