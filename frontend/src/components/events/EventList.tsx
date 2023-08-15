import { FC } from "react"
import dayjs from "dayjs"

import { Event, useEvent } from "@/api/events/hooks"
import { formatDate } from "@/helpers/dates"
import { joinEvent, leaveEvent } from "@/api/events/respond"

const EventList: FC<{ events: number[] }> = ({ events }) => {
  if (events.length == 0) {
    return <div className={"text-info text-center mt-5 mb-7"}>There are no events yet...</div>
  }

  return (
    <>
      {events.map((eventId) => (
        <EventCard eventId={eventId} key={eventId} />
      ))}
    </>
  )
}

const EventCard: FC<{ eventId: number }> = ({ eventId }) => {
  const { event } = useEvent(eventId)

  if (!event) return null

  return (
    <div
      className={"max-w-3xl mx-auto rounded-lg border border-base-100 shadow-base-100 shadow-sm"}
    >
      <div className={"mx-5 my-3"}>
        <h3 className={"font-Alatsi text-base start-dot end-dot mb-3"}>{event.title}</h3>
        <p className={"text-sm text-base-content"}>{event.description}</p>
      </div>
      <div
        className={
          "bg-base-100 m-1 rounded-lg flex flex-wrap justify-between items-center py-1.5 px-3 gap-x-1"
        }
      >
        <span>{formatDate(event.time_and_date)}</span>
        <EventResponse event={event} />
      </div>
    </div>
  )
}

const EventResponse: FC<{ event: Event }> = ({ event }) => {
  const { mutate } = useEvent(event.id)

  if (event.status === undefined) return null

  if (dayjs(event.time_and_date).isBefore(dayjs())) {
    return <span className={"self-center ml-auto font-light text-sm text-info"}>Finished</span>
  }

  return (
    <span className={"flex gap-2"}>
      <button
        onClick={() => {
          leaveEvent(event.id)
            .catch(() => null)
            .finally(() =>
              mutate({
                ...event,
                status: 0,
              })
            )
        }}
        className={
          "btn btn-sm normal-case text-white font-light" +
          " " +
          (event.status === 2
            ? "btn-primary"
            : event.status === 0
            ? "btn-secondary"
            : "btn-neutral")
        }
      >
        {"I'm out!"}
      </button>
      <button
        onClick={() => {
          joinEvent(event.id)
            .catch(() => null)
            .finally(() =>
              mutate({
                ...event,
                status: 1,
              })
            )
        }}
        className={
          "btn btn-sm normal-case text-white font-light" +
          " " +
          (event.status === 2
            ? "btn-primary"
            : event.status === 1
            ? "btn-secondary"
            : "btn-neutral")
        }
      >
        {"I'm in!"}
      </button>
    </span>
  )
}

export default EventList
