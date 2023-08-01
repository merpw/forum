import { FC } from "react"

const Bio: FC = () => {
  return (
    <div className={"form-control"}>
      <label className={"label"}>
        <span className={"label-text"}>Bio</span>
      </label>
      <textarea
        name={"bio"}
        className={"textarea textarea-bordered resize-none"}
        rows={2}
        placeholder={"about me"}
        maxLength={200}
      />
    </div>
  )
}

export default Bio
