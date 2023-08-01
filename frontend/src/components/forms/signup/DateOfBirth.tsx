import { FC } from "react"

const DateOfBirth: FC<{ dob: string }> = ({ dob }) => {
  return (
    <div className={"flex flex-wrap flex-row gap-3"}>
      <div className={"grow basis-1/3"}>
        <label className={"label"}>
          <span className={"label-text"}>Date of Birth</span>
        </label>
        <input
          type={"date"}
          min={"1900-01-01"}
          max={new Date().toISOString().split("T")[0]}
          className={"input input-bordered w-full text-sm"}
          name={"dob"}
          value={dob}
          onChange={() => void 0 /* handled by Form */}
          placeholder={"Date of Birth"}
          required
        />
      </div>
      <div className={"grow basis-1/3"}>
        <label className={"label"}>
          <span className={"label-text"}>Gender</span>
        </label>
        <div>
          <select
            title={"Your gender"}
            className={"input input-bordered w-full text-sm"}
            name={"gender"}
            placeholder={"Gender"}
            required
          >
            <option value={""}>Select</option>
            <option value={"male"}>Male</option>
            <option value={"female"}>Female</option>
            <option value={"other"}>Other</option>
          </select>
        </div>
      </div>
    </div>
  )
}

export default DateOfBirth
