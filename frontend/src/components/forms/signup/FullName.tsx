import { FC } from "react"

const FullName: FC<{ first_name: string; last_name: string }> = ({ first_name, last_name }) => {
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
          value={first_name}
          onChange={() => void 0 /* handled by Form */}
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
          value={last_name}
          onChange={() => void 0 /* handled by Form */}
          maxLength={15}
          required
        />
      </div>
    </div>
  )
}

export default FullName
