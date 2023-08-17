import { notFound } from "next/navigation"
import { Metadata } from "next"

import { getGroupLocal } from "@/api/groups/edge"
import GroupPage from "@/app/group/[id]/group-page"

type Props = { params: { id: string } }

export const revalidate = 60

export const generateMetadata = async ({ params }: Props): Promise<Metadata> => {
  const id = +params.id
  if (isNaN(id)) {
    return notFound()
  }
  const group = await getGroupLocal(id).catch(notFound)

  return {
    title: `Events - ${group.title}`,
  }
}

const Page = async ({ params }: Props) => {
  const group = await getGroupLocal(+params.id).catch(notFound)

  return <GroupPage group={group} tab={"Events"} />
}

export default Page
