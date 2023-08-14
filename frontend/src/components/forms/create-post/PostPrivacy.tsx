"use client"

import { useState } from "react"

import SelectUsers from "@/components/forms/SelectUsers"
import { useFollowers } from "@/api/followers/hooks"
import { Privacy } from "@/custom"

const PostPrivacy = () => {
  const { followers } = useFollowers()
  const [privacy, setPrivacy] = useState<Privacy>("Public")
  const [selected, setSelected] = useState<string[]>([])

  if (!followers) return null

  const handlePrivacyChange = (selectedPrivacy: Privacy) => {
    console.log("Selected Privacy:", selectedPrivacy)
    setPrivacy(selectedPrivacy)
  }

  return (
    <div className={"space-y-2 w-full"}>
      <div className={"flex flex-row items-center gap-1"}>
        <button
          onClick={() => handlePrivacyChange("Public")}
          name={"privacy"}
          value={"Public"}
          type={"button"}
          className={`btn btn-xs ${privacy === "Public" ? "btn-primary" : "btn-neutral"} `}
        >
          public
        </button>
        <button
          onClick={() => handlePrivacyChange("Private")}
          name={"privacy"}
          value={"Private"}
          type={"button"}
          className={`btn btn-xs ${privacy === "Private" ? "btn-primary" : "btn-neutral"} `}
        >
          private
        </button>
        <button
          onClick={() => handlePrivacyChange("SuperPrivate")}
          name={"privacy"}
          value={"select"}
          type={"button"}
          className={`btn btn-xs ${privacy === "SuperPrivate" ? "btn-primary" : "btn-neutral"}`}
        >
          select
        </button>
      </div>
      {privacy === "SuperPrivate" && (
        <div className={"w-full"}>
          <SelectUsers
            key={"select_" + selected.join(",")}
            name={"selected_users"}
            value={selected}
            userIds={followers}
            onChange={(newValue) => setSelected(newValue as never)}
            escapeClearsValue={true}
            noOptionsMessage={() =>
              followers.length > 0 ? "All followers are already invited" : "No followers to invite"
            }
          />
        </div>
      )}
    </div>
  )
}

export default PostPrivacy
