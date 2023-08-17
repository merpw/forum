"use client"

import { useState } from "react"

import SelectUsers from "@/components/forms/SelectUsers"
import { useFollowers } from "@/api/followers/hooks"

const PostPrivacy = () => {
  const { followers } = useFollowers()
  const [privacy, setPrivacy] = useState(0)
  const audience = [].map(Number)

  if (!followers) return null

  const handlePrivacyChange = (selectedPrivacy: number) => {
    setPrivacy(selectedPrivacy)
  }

  return (
    <div className={"space-y-2 w-full"}>
      <input type={"hidden"} name={"privacy"} value={privacy} />
      <div className={"flex flex-row items-center gap-1"}>
        <button
          onClick={() => handlePrivacyChange(0)}
          value={0}
          type={"button"}
          className={`btn btn-xs ${privacy === 0 ? "btn-primary" : "btn-neutral"} `}
        >
          public
        </button>
        <button
          onClick={() => handlePrivacyChange(1)}
          value={1}
          type={"button"}
          className={`btn btn-xs ${privacy === 1 ? "btn-primary" : "btn-neutral"} `}
        >
          private
        </button>
        <button
          onClick={() => handlePrivacyChange(2)}
          value={2}
          type={"button"}
          className={`btn btn-xs ${privacy === 2 ? "btn-primary" : "btn-neutral"}`}
        >
          select
        </button>
      </div>
      {privacy === 2 && (
        <div className={"w-full"}>
          <SelectUsers
            key={"select_" + audience.join(",")}
            name={"audience"}
            userIds={followers}
            placeholder={"Select followers"}
            escapeClearsValue={true}
            noOptionsMessage={() =>
              followers.length > 0 ? "All followers are already invited" : "No followers to invite"
            }
            required
          />
        </div>
      )}
    </div>
  )
}

export default PostPrivacy
