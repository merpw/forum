import { ChatUser, Message } from "../types"
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
  list: [] as Message[]
}


async function getChatUsers() {
  Object.assign(chatList, {
    Ids: new Map<number, number>,
    Users: [],
  })

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
  console.log(chatList)
  const userIds: iterator = await getUserIds()
  // if (userIds.length - 1 == chatList.Ids.size && chatList.Ids.size != 0){
  //   return
  // }

  setTimeout(async () => {
    for (const id of Object.values(userIds)){
      if (id !== userInfo.Id){
        chatList.Users.push(await getUserById(id))

      }
    }
  }, 100)

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
  }, 150)


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
  Object.assign(messages, {list: []})
  const chatArea = document.getElementById("chat-area") as HTMLDivElement
  const chat = createElement("div", "chat show-chat")

  const chatId = chatList.Ids.get(id) as number
  console.log(chatId)

  const user = await getUserById(id.toString())

  const chatName = createElement("div", "chat-name", null, user.Name) as HTMLDivElement
  chatName.style.paddingBottom = "10px;"

  const chatMessages = createElement("div", "chat-messages", `Chat${chatId}`)
  await getMessages(chatId)


  setTimeout(() => {
    console.log(chatId)
    for (const msg of Object.values(messages.list)) {
      console.log(msg)
      const msgElement = createElement("div", "message", null, msg.content) as HTMLDivElement 
      chatMessages.appendChild(msgElement)
    }
  }, 100)
  const chatFormContainer = createElement("div", "chat-form-container")
  const chatForm = createElement("form", "chat-form")
  chatFormContainer.appendChild(chatForm)
  chat.append(chatName, chatMessages, chatFormContainer)
  chatArea.replaceChildren(chat)
}

const sendMessage = async (chatId: number) => {

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

let parsedUsers = false
export const displayChatUsers = async () => {

  const onlineList = document.getElementById("online-users") as HTMLUListElement,
  offlineList = document.getElementById("offline-users") as HTMLUListElement,
  onlineTitle = document.getElementById("online-title") as HTMLElement,
  offlineTitle = document.getElementById("offline-title") as HTMLElement

  /* Add eventlisteners to show/hide users in chat list */
  onlineTitle.addEventListener("click", toggleOnline)
  offlineTitle.addEventListener("click", toggleOffline)

  await getChatUsers()
  setTimeout(() => {
    sendWsObject(
      {
        type: "get",
        item: {
          url: "/chat/all"
        }
      }
    )
  }, 600)

  setTimeout(() => {
    // TODO: Remove this sort when sorted in backend.
    console.log(chatList)
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

    const chatUsers = document.querySelectorAll(".chat-user") as NodeListOf<HTMLElement>
    if (chatUsers.length < chatList.Users.length) {
      for (const u of chatList.Users) {
        const user = createElement("li", "chat-user", u.Name)
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
          user.classList.add("online")
          onlineList.appendChild(user)
        } else {
          user.classList.add("offline")
          offlineList.appendChild(user)
        }
      }
    }

  }, 200)

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
