import { FC } from "react"

const Username: FC<{ username: string }> = ({ username }) => {
  return (
    <div className={"form-control mt-3"}>
      <label className={"label pt-0"}>
        <span className={"label-text"}>Username</span>
      </label>
      <input
        type={"text"}
        className={"input input-bordered"}
        name={"username"}
        minLength={3}
        maxLength={15}
        placeholder={"username"}
        value={username}
        onChange={() => void 0 /* handled by Form */}
        required
      />
    </div>
  )
}

export default Username
