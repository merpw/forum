import { ChangeEvent } from "react"

export const trimInput = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
  e.target.value = e.target.value.trim()
}
