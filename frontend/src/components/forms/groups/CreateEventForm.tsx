import { FC, useState } from "react"
import ReactTextAreaAutosize from "react-textarea-autosize"
import dayjs from "dayjs"
import { useSWRConfig } from "swr"
import { useRouter } from "next/navigation"

import { FormError } from "@/components/error"
import createEvent from "@/api/events/create"

const CreateEventForm: FC<{ groupId: number }> = ({ groupId }) => {
  const [formError, setFormError] = useState<string | null>(null)

  const { mutate } = useSWRConfig()

  const router = useRouter()

  return (
    <form
      onReset={() => {
        if (formError != null) {
          setFormError(null)
        }
      }}
      onChange={() => setFormError(null)}
      className={"modal-box bg-base-200 overflow-visible h-fit"}
      onSubmit={async (e) => {
        e.preventDefault()
        const form = e.currentTarget as HTMLFormElement
        const formData = new FormData(form)

        const event = {
          title: formData.get("title") as string,
          description: formData.get("description") as string,
          time_and_date: formData.get("time_and_date") as string,
        }

        formError && setFormError(null)

        createEvent(event, groupId)
          .then(() => {
            form.reset()
            form.closest("dialog")?.close()
            router.push(`/group/${groupId}/events`)
            mutate(`/api/groups/${groupId}/events`)
          })
          .catch((error) => {
            setFormError(error.message)
          })
      }}
    >
      <h2 className={"text-xl mb-3"}>Create event</h2>

      <div className={"form-control"}>
        <label className={"label"}>
          <span className={"label-text"}>Title</span>
        </label>
        <input
          type={"text"}
          name={"title"}
          placeholder={"Event title"}
          className={"input input-bordered"}
          required
        />
      </div>

      <div className={"form-control"}>
        <label className={"label"}>
          <span className={"label-text"}>Description</span>
        </label>
        <ReactTextAreaAutosize
          name={"description"}
          placeholder={"Event description"}
          className={"textarea textarea-bordered"}
          minRows={2}
          maxRows={5}
          required
        />
      </div>

      <div className={"form-control"}>
        <label className={"label"}>
          <span className={"label-text"}>Time and date</span>
        </label>
        <input
          defaultValue={dayjs().add(1, "day").format("YYYY-MM-DDT10:00")}
          min={dayjs().format("YYYY-MM-DDTHH:mm")}
          type={"datetime-local"}
          name={"time_and_date"}
          placeholder={"Event time and date"}
          className={"input input-bordered"}
        />
      </div>

      <div className={"modal-action items-center justify-between"}>
        <FormError error={formError} />
        <button className={"ml-auto btn btn-primary"}>Create</button>
      </div>
    </form>
  )
}

export default CreateEventForm
