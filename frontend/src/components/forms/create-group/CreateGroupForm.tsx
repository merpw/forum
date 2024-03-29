"use client"

import { FC, useState } from "react"
import ReactTextAreaAutosize from "react-textarea-autosize"
import { useRouter } from "next/navigation"

import { FormError } from "@/components/error"
import { trimInput } from "@/helpers/input"
import SelectUsers from "@/components/forms/SelectUsers"
import { useFollowers } from "@/api/followers/hooks"
import { createGroup } from "@/api/groups/create"

const CreateGroupForm: FC = () => {
  const [isSame, setIsSame] = useState(false)

  const [formError, setFormError] = useState<string | null>(null)

  const router = useRouter()

  const { followers } = useFollowers()

  return (
    <form
      onChange={() => setIsSame(false)}
      onSubmit={async (e) => {
        e.preventDefault()

        if (isSame) return

        const form = e.currentTarget as HTMLFormElement

        const formData = new FormData(form)

        const formFields = {
          title: formData.get("title") as string,
          description: formData.get("description") as string,
          invite: formData.getAll("invite") as string[],
        }

        if (formError != null) setFormError(null)
        setIsSame(true)

        createGroup(formFields.title, formFields.description, formFields.invite.map(Number))
          .then((groupId) => {
            router.push(`/group/${groupId}`)
          })
          .catch((error) => {
            setFormError(error.message)
            setIsSame(false)
          })
      }}
    >
      <div className={"max-w-3xl mx-auto flex-col"}>
        <div className={"card flex-shrink-0 bg-base-200 shadow-md"}>
          <div className={"card-body"}>
            <div className={"my-2 font-Yesteryear text-3xl text-primary opacity-50 text-center"}>
              <span className={"text-neutral-content"}> {"How about a new group?"}</span>
            </div>

            <div className={"form-control"}>
              <input
                onBlur={trimInput}
                type={"text"}
                name={"title"}
                className={"input input-bordered bg-base-100 text-sm"}
                placeholder={"Group name"}
                required
                maxLength={25}
              />
            </div>

            <ReactTextAreaAutosize
              name={"description"}
              required
              onBlur={trimInput}
              className={"textarea textarea-bordered"}
              rows={2}
              minRows={2}
              placeholder={"Group description"}
              maxLength={205}
            />

            {followers && (
              <SelectUsers
                userIds={followers}
                name={"invite"}
                placeholder={"Invite your followers"}
                noOptionsMessage={() => "No followers to invite"}
              />
            )}

            <FormError error={formError} />

            <div className={"form-control mt-3"}>
              <button type={"submit"} className={"button self-center"}>
                <svg
                  xmlns={"http://www.w3.org/2000/svg"}
                  fill={"none"}
                  viewBox={"0 0 24 24"}
                  strokeWidth={1.5}
                  stroke={"currentColor"}
                  className={"w-4 h-4"}
                >
                  <path
                    strokeLinecap={"round"}
                    strokeLinejoin={"round"}
                    d={
                      "M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"
                    }
                  />
                </svg>
                SUBMIT
              </button>
            </div>
          </div>
        </div>
      </div>
    </form>
  )
}

export default CreateGroupForm
