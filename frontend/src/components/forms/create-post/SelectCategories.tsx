import { FC, useId } from "react"
import Select from "react-select"

import { Capitalize } from "@/helpers/text"

// TODO: maybe replace with another select component

const SelectCategories: FC<{ categories: string[] }> = ({ categories }) => {
  return (
    <div className={"mb-3"}>
      <Select
        placeholder={"Categories"}
        instanceId={useId()}
        isMulti={true}
        name={"categories"}
        className={"react-select-container"}
        classNamePrefix={"react-select"}
        options={categories.map((name) => ({ label: Capitalize(name), value: name }))}
        required
      />
    </div>
  )
}

export default SelectCategories
