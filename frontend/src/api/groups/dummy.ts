/* 0-not member, 1-member, 2-requested */
import { Group } from "@/api/groups/hooks"

export const DUMMY_GROUPS: Group[] = [
  {
    id: 1,
    title: "Google",
    description: "Google Inc.",
    member_status: 0,
    member_count: 120,
    creator_id: 2,
  },
  {
    id: 2,
    title: "Gritlab",
    description: "The best coding school in the world",
    member_status: 1,
    member_count: 60,
    creator_id: 2,
  },
  {
    id: 3,
    title: "Facebook",
    description: "Facebook Inc.",
    member_status: 2,
    member_count: 1,
    creator_id: 2,
  },
]
