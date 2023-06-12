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

export let currentChat = -1

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
  currentChat = chatId

  const user = await getUserById(id.toString())

  const chatName = createElement("div", "chat-name", null, user.Name) as HTMLDivElement
  chatName.style.paddingBottom = "10px"

  const chatMessages = createElement("div", "chat-messages", `Chat${chatId}`)

  const chatFormContainer = createElement("div", "chat-form-container")
  const chatForm = createElement("form", "chat-form")
  const messageField = createElement("input", null, "chat-text")
  messageField.setAttribute("maxlength", "150")
  chatMessages.addEventListener("messageEvent", () => {
    updateChat(chatId)
  })

  const messageSend = createElement("button", null, "chat-send", "Send")
  messageSend.addEventListener("click", () => {
    sendMessage(chatId)
  })
  await getMessages(chatId)

  setTimeout(() => {
    messages.list.sort((a, b) => {
      return a.id - b.id
    }) 
  }, 220)

  chatForm.append(messageField, messageSend)
  chatFormContainer.appendChild(chatForm)

  setTimeout(() => {
    if (messages.list.length === 0) {
      return
    }
    for (const message of Object.values(messages.list.reverse())) {
      const msgElement = createElement("div", "message") as HTMLDivElement
      const date = new Date(message.timestamp)
      const formatDate = date.toLocaleString("en-GB", { timeZone: "EET" }).slice(0, 10)
      if (message.authorId === userInfo.Id){
        msgElement.classList.add("send")
        msgElement.textContent = message.content + '\n' + formatDate
      } else if (message.authorId === -1){
        msgElement.classList.add("status")
        msgElement.textContent = message.content + '\n'  + formatDate
      } else {
        msgElement.classList.add("recieve")
        msgElement.textContent = message.content + '\n' + formatDate
      }
      chatMessages.appendChild(msgElement)
    }
  }, 250)
  chat.append(chatName, chatMessages, chatFormContainer)
  chatArea.replaceChildren(chat)
}

const sendMessage = async (chatId: number) => {
  const content = document.getElementById("chat-text") as HTMLInputElement
  let message = content.value.toString().trim()
  if (message.length > 0) {
    content.value = ""
    if (message.length > 150) {
      message = message.slice(0, 150)
    }
    sendWsObject(
      {
        type: "post",
        item: {
          url: `/chat/${chatId}/message`,
          data: {
            content: message        }
        }
      }
    )
      getMessages(chatId)
  }
}

export const updateChat = (chatId: number) => {
  const message = messages.list.reverse()[0]
  const chatMessages = document.getElementById(`Chat${chatId}`) as HTMLDivElement
  const msgElement = createElement("div", "message", null, message.content) as HTMLDivElement
  console.log(message, messages.list[0].id)
  if (message.id < messages.list[messages.list.length - 1].id || message.authorId == -1) {
    return
  }
      const date = new Date(message.timestamp)
      const formatDate = date.toLocaleString("en-GB", { timeZone: "EET" }).slice(0, 10)
      if (message.authorId === userInfo.Id){
        msgElement.classList.add("send")
        msgElement.textContent = message.content + "\n" + formatDate
      } else if (message.authorId === -1){
        msgElement.classList.add("status")
        msgElement.textContent = message.content + "\n"  + formatDate
      } else {
        msgElement.textContent = message.content + "\n" + formatDate
        msgElement.classList.add("recieve")
      }
  if (chatId === currentChat) {
    chatMessages.prepend(msgElement)
  }
}


const getMessages = async (chatId: number) => {
  Object.assign(messages, {list: []})
  setTimeout(() => {
    sendWsObject(
      {
        type: "get",
        item: {
          url: `/chat/${chatId}/messages`
        }
      })
  }, 150)
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
