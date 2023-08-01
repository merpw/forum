import { FC } from "react"

import { trimInput } from "@/helpers/input"

const Email: FC = () => {
  return (
    <div className={"form-control"}>
      <label className={"label"}>
        <span className={"label-text"}>Email</span>
      </label>
      <input
        type={"email"}
        className={"input input-bordered"}
        name={"email"}
        placeholder={"email"}
        onBlur={trimInput}
        required
      />
    </div>
  )
}

export default Email
