import { FC } from "react"

import { trimInput } from "@/helpers/input"

const FullName: FC = () => {
  return (
    <div className={"flex flex-wrap flex-row gap-3"}>
      <div className={"grow basis-1/3"}>
        <label className={"label"}>
          <span className={"label-text"}>First Name </span>
        </label>
        <input
          type={"text"}
          className={"input input-bordered w-full"}
          name={"first_name"}
          placeholder={"first name"}
          onBlur={trimInput}
          maxLength={15}
          required
        />
      </div>
      <div className={"grow basis-1/3"}>
        <label className={"label"}>
          <span className={"label-text"}>Last Name </span>
        </label>
        <input
          type={"text"}
          className={"input input-bordered w-full"}
          name={"last_name"}
          placeholder={"last name"}
          onBlur={trimInput}
          maxLength={15}
          required
        />
      </div>
    </div>
  )
}

export default FullName
