import { Metadata } from "next"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"

import checkSession from "@/api/auth/edge"
import CreateGroupForm from "@/components/forms/create-group/CreateGroupForm"

export const metadata: Metadata = {
  title: "Create a new group",
}

const Page = async () => {
  const token = cookies().get("forum-token")?.value
  if (!token) {
    return redirect("/login")
  }
  try {
    await checkSession(token)
  } catch (error) {
    return redirect("/login")
  }

  return <CreateGroupForm />
}

export default Page
