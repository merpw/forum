import { components, MultiValueProps, OptionProps, Props as ReactSelectProps } from "react-select"

import Select from "@/components/Select"
import { User } from "@/custom"
import { useUserList } from "@/api/users/hooks"
import Avatar from "@/components/Avatar"

export const SelectUsersWithData = (props: ReactSelectProps & { users: User[] }) => {
  const cleanProps = { ...props, users: undefined }

  const { users } = props

  const options: UserOption[] = users.map((user) => ({
    value: user.id,
    label: user.username,
    user,
  }))

  // TODO: improve typings
  return (
    <Select
      {...cleanProps}
      isMulti={true}
      closeMenuOnSelect={false}
      options={options}
      components={{ Option: Option as never, MultiValueLabel: MultiValueLabel as never }}
    />
  )
}

type UserOption = { value: number; label: string; user: User }

const Option = (props: OptionProps<UserOption>) => {
  const { user } = props.data

  return (
    <components.Option {...props}>
      <div className={"flex items-center gap-2"}>
        <Avatar user={user} size={10} className={"h-10"} />
        <span>{user.username}</span>
      </div>
    </components.Option>
  )
}

const MultiValueLabel = (props: MultiValueProps<UserOption>) => {
  const { user } = props.data

  return (
    <div className={"flex items-center gap-2 p-1.5"}>
      <Avatar user={user} size={10} className={"h-7"} />
      <span>{user.username}</span>
    </div>
  )
}

const SelectUsers = (
  props: ReactSelectProps & {
    userIds: number[]
  }
) => {
  const { users } = useUserList(props.userIds)

  if (!users) return null

  return <SelectUsersWithData {...props} users={users} />
}

export default SelectUsers
