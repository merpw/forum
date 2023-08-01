import { FC } from "react"

const Email: FC<{ email: string }> = ({ email }) => {
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
        value={email}
        onChange={() => void 0 /* handled by Form */}
        required
      />
    </div>
  )
}

export default Email
