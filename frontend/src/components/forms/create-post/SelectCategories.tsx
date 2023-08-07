import { FC, useId } from "react"

import { Capitalize } from "@/helpers/text"
import Select from "@/components/Select"

const SelectCategories: FC<{ categories: string[] }> = ({ categories }) => {
  return (
    <div className={"mb-3"}>
      <Select
        name={"categories"}
        placeholder={"Categories"}
        instanceId={useId()}
        isMulti={true}
        options={categories.map((name) => ({ label: Capitalize(name), value: name }))}
        required
      />
    </div>
  )
}

export default SelectCategories
