import { ChatUser } from "../types"
import { sendWsObject } from "./ws.js"
import { userInfo } from "./auth.js"
import { createElement, iterator } from "./utils.js"
import { getUserById, getUserIds } from "../api/get.js"
// This file is dedicated to sorting and displaying chat users in the sidebar.
export const chatList = {
  Ids: new Map<number, number>, //userId, chatId
  Users: [] as ChatUser[]
}

export const messages = {
  list: [] as object[]
}

async function getChatUsers() {
  const userIds: iterator = await getUserIds()
  if (userIds.length - 1 == chatList.Ids.size && chatList.Ids.size != 0){
    return
  }
  for (const id of Object.values(userIds)){
    if (id !== userInfo.Id){
      chatList.Users.push(await getUserById(id))

    }
  }
  setTimeout(() => {
    for (const user of chatList.Users){
      if (!chatList.Ids.has(user.Id)) {
        sendWsObject(
          {
            type: "post",
            item: {
              url: "/chat/create",
              data: {
                userId: user.Id
              }
            }
          }
        )
      }
    }

  }, 50)
}
  //   ws.send(JSON.stringify({
  //       type: "get",
  //       item: {
  //         url: `/users/online`,
  //         data: {
  //           content: message,
  //         }
  //       }
  //     })
  // ) TODO: Fix so that online users and UnreadMsg are updated in chatList

// Display the current chat based on chatId
const showChat = async (id: number) => {
  const chatId = chatList.Ids.get(id) as number
  console.log("ehm")
  const chatArea = createElement("div")
  const chat = createElement("div", "chat show-chat")
  const chatMessages = createElement("div", "chat-messages", `Chat${chatId}`)
  await getMessages(chatId)
  const chatFormContainer = createElement("div", "chat-form-container")
  const chatForm = createElement("form", "chat-form")
  chatArea.replaceChildren(chat)
}

const getMessages = async (chatId: number) => {
  sendWsObject(
    {
      type: "get",
      item: {
        url: `/chat/${chatId}/messages`
      }
    })
}

export const displayChatUsers = async () => {
  const onlineList = document.getElementById("online-users") as HTMLUListElement,
  offlineList = document.getElementById("offline-users") as HTMLUListElement,
  onlineTitle = document.getElementById("online-title") as HTMLElement,
  offlineTitle = document.getElementById("offline-title") as HTMLElement

  /* Add eventlisteners to show/hide users in chat list */
  onlineTitle.addEventListener("click", toggleOnline)
  offlineTitle.addEventListener("click", toggleOffline)
  await getChatUsers()

  chatList.Users.sort((a, b) => {
    const name1 = a.Name.toLowerCase()
    const name2 = b.Name.toLowerCase()
    if (name1 < name2) {
      return -1
    }
    if (name1 > name2) {
      return 1
    }
    return 0
  })
  for (const u of chatList.Users) {
    const user = createElement("li")
    const name = createElement("p", null, null, u.Name)
    user.appendChild(name)

    if (u.UnreadMsg) {
      const unread = createElement("i", "bx bx-message-dots")
      user.appendChild(unread)
    }

    user.addEventListener("click", () => {
      showChat(u.Id)
    })

    if (u.Online) {
      user.className = "online"
      onlineList.appendChild(user)
    } else {
      user.className = "offline"
      offlineList.appendChild(user)
    }
  }
  setTimeout(() => {
    sendWsObject(
      {
        type: "get",
        item: {
          url: "/chat/all"
        }
      }
    )

  }, 100)

}

const toggleOnline = () => {
  const onlineToggle = document.getElementById("online-toggle") as HTMLElement,
  onlineUsers = document.getElementById("online-users") as HTMLUListElement
  if (onlineToggle.className === "bx bx-chevron-down") {
    onlineToggle.className = "bx bx-chevron-right"
    onlineUsers.style.display = "none"
  } else {
    onlineToggle.className = "bx bx-chevron-down"
    onlineUsers.style.display = "block"
  }
}

const toggleOffline = () => {
  const offlineToggle = document.getElementById("offline-toggle") as HTMLElement
  const offlineUsers = document.getElementById("offline-users") as HTMLUListElement
  if (!offlineToggle) return

  if (offlineToggle.className === "bx bx-chevron-down") {
    offlineToggle.className = "bx bx-chevron-right"
    offlineUsers.style.display = "none"
  } else {
    offlineToggle.className = "bx bx-chevron-down"
    offlineUsers.style.display = "block"
  }
}
