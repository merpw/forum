import { FC } from "react"

const Username: FC = () => {
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
        onInput={(e) => {
          const updatedValue = e.currentTarget.value.replace(/\W/g, "")
          if (updatedValue !== e.currentTarget.value) {
            console.log(
              "updatedValue",
              updatedValue,
              "e.currentTarget.value",
              e.currentTarget.value
            )
            e.currentTarget.value = updatedValue
            e.currentTarget.setCustomValidity(
              "Only latin letters, numbers, and underscores are allowed."
            )
            e.currentTarget.reportValidity()
          } else {
            e.currentTarget.setCustomValidity("")
          }
        }}
        required
      />
    </div>
  )
}

export default Username
