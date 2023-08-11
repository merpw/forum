import { FC, useState } from "react"

import { useFollowers } from "@/api/followers/hooks"
import SelectUsers from "@/components/forms/SelectUsers"
import { inviteToGroup } from "@/api/groups/invite"
import { FormError } from "@/components/error"
import { useGroupMembers } from "@/api/groups/hooks"

const InviteFollowersForm: FC<{ groupId: number }> = ({ groupId }) => {
  const { followers } = useFollowers()
  const { members: alreadyMembers } = useGroupMembers(groupId, true)

  const [formError, setFormError] = useState<string | null>(null)
  const [members, setMembers] = useState<string[]>([])

  if (!followers || !alreadyMembers) return null

  const userIds = followers.filter((followerId) => !alreadyMembers.includes(followerId))

  return (
    <form
      onReset={() => {
        if (formError != null) {
          setFormError(null)
        }
        setMembers([])
      }}
      onBlur={(e) => e.currentTarget.reset()}
      onChange={() => setFormError(null)}
      className={"modal-box bg-base-200 overflow-visible h-fit"}
      onSubmit={async (e) => {
        e.preventDefault()
        const form = e.currentTarget as HTMLFormElement
        const members = new FormData(form).getAll("members") as string[]

        formError && setFormError(null)

        Promise.all(members.map(Number).map((userId) => inviteToGroup(groupId, userId)))
          .then(() => {
            form.reset()
            form.closest("dialog")?.close()
          })
          .catch((error) => {
            console.error(error)
            setFormError(error.message)
          })
      }}
    >
      <h2 className={"text-xl mb-3"}>Invite your followers</h2>

      <SelectUsers
        key={"select_" + members.join(",")}
        name={"members"}
        value={members}
        userIds={userIds}
        onChange={(newValue) => setMembers(newValue as never)}
        escapeClearsValue={true}
        noOptionsMessage={() =>
          followers.length > 0 ? "All followers are already invited" : "No followers to invite"
        }
      />

      <div className={"modal-action items-center justify-between"}>
        <FormError error={formError} />
        <button className={"ml-auto btn btn-primary"}>Invite</button>
      </div>
    </form>
  )
}

export default InviteFollowersForm
