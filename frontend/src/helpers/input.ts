import { ChangeEvent } from "react"

export const trimInput = (e: ChangeEvent<HTMLInputElement>) => {
  e.target.value = e.target.value.trim()
}
