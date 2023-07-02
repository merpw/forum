import { redirect } from "next/navigation"

import ChatPage from "@/app/chat/[id]/chat-page"

type Props = { params: { id: string } }

const Page = ({ params }: Props) => {
  const id = +params.id
  if (isNaN(id) || id < 0) {
    return redirect("/chat")
  }

  return <ChatPage id={id} />
}

export default Page
