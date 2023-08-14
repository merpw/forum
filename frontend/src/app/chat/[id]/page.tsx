import { redirect } from "next/navigation"
import { cookies } from "next/headers"

import ChatIdPage from "@/components/chats/page/ChatIdPage"
import checkSession from "@/api/auth/edge"
import { getUserLocal } from "@/api/users/edge"
import { AssociatedIdChatPage } from "@/components/chats/page/AssociatedIdChatPage"
import { getGroupLocal } from "@/api/groups/edge"

type Props = { params: { id: string } }

export const generateMetadata = async ({ params }: Props) => {
  if (params.id.startsWith("u")) {
    const userId = +params.id.slice(1)
    if (isNaN(userId) || userId < 0) {
      return redirect("/chat")
    }

    try {
      const user = await getUserLocal(userId)
      return {
        title: `Chat with ${user.username}`,
      }
    } catch (error) {
      return redirect("/chat")
    }
  }

  if (params.id.startsWith("g")) {
    const groupId = +params.id.slice(1)
    if (isNaN(groupId) || groupId < 0) {
      return redirect("/chat")
    }

    try {
      const group = await getGroupLocal(groupId)
      return {
        title: `Chat - ${group.title}`,
      }
    } catch (error) {
      return redirect("/chat")
    }
  }

  return {
    title: `Chats`,
  }
}

const Page = async ({ params }: Props) => {
  const token = cookies().get("forum-token")?.value
  if (!token) {
    return redirect("/login")
  }
  const requestUserId = await checkSession(token)
  if (!requestUserId) {
    return redirect("/login")
  }

  if (params.id.startsWith("u")) {
    try {
      const userId = +params.id.slice(1)
      const companion = await getUserLocal(userId)
      return <AssociatedIdChatPage companionId={companion.id} />
    } catch (error) {
      return redirect("/chat")
    }
  }

  if (params.id.startsWith("g")) {
    try {
      const groupId = +params.id.slice(1)
      const group = await getGroupLocal(groupId)
      return <AssociatedIdChatPage groupId={group.id} />
    } catch (error) {
      return redirect("/chat")
    }
  }

  const id = +params.id
  if (isNaN(id) || id < 0) {
    return redirect("/chat")
  }

  return <ChatIdPage id={id} />
}

export default Page
