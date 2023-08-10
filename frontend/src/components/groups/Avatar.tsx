import { FC } from "react"
import { TiGroup } from "react-icons/ti"

import { Group } from "@/api/groups/hooks"

const GroupAvatar: FC<{ group: Group; className?: string; size: number }> = ({
  className,
  size,
}) => {
  return (
    <div className={className}>
      <Placeholder className={"w-full h-full"} size={size} />
    </div>
  )
}

const Placeholder: FC<{ className?: string; size: number }> = ({ className, size }) => {
  return <TiGroup size={size} className={"opacity-30" + " " + className} />
}

export default GroupAvatar
