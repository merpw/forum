import dayjs from "dayjs"
import utc from "dayjs/plugin/utc"
import relativeTime from "dayjs/plugin/relativeTime"
import { useEffect, useState } from "react"

const REFRESH_INTERVAL = 30 * 1000

const useDates = (date: string | number) => {
  dayjs.extend(utc)
  dayjs.extend(relativeTime)

  const [localDate, setLocalDate] = useState(dayjs(date).utc().format("YYYY-MM-DD HH:mm:ss UTC"))
  const [relativeDate, setRelativeDate] = useState(dayjs(date).fromNow())

  useEffect(() => {
    setLocalDate(dayjs(date).local().format("YYYY-MM-DD HH:mm:ss"))

    setRelativeDate(dayjs(date).fromNow())

    const updateRelativeDate = () => {
      if (!(relativeDate.includes("minute") || relativeDate.includes("second"))) {
        clearInterval(interval)
      }
      setRelativeDate(dayjs(date).fromNow())
    }

    const interval = setInterval(updateRelativeDate, REFRESH_INTERVAL)
    updateRelativeDate()
  }, [relativeDate, date])

  return { localDate, relativeDate }
}

export default useDates
