import dayjs from "dayjs"

export const formatDate = (timestamp: string) => {
  const date = dayjs(timestamp)
  const today = dayjs()
  const tomorrow = dayjs().add(1, "day")

  if (date.format("YYYY-MM-DD") === tomorrow.format("YYYY-MM-DD")) {
    return "Tomorrow"
  }

  if (date.format("YYYY-MM-DD") === today.format("YYYY-MM-DD")) {
    return "Today"
  }

  if (date.add(1, "day").format("YYYY-MM-DD") === today.format("YYYY-MM-DD")) {
    return "Yesterday"
  }

  return date.format(date.year() === dayjs().year() ? "MMMM D" : "MMMM D, YYYY")
}
