import { ChatUser, Message } from "../types"
import { sendWsObject } from "./ws.js"
import { userInfo } from "./auth.js"
import { createElement, iterator } from "./utils.js"
import { getUserById, getUserIds } from "../api/get.js"

// This file is dedicated to sorting and displaying chat users in the sidebar.
export const chatList = {
  Ids: new Map<number, number>(), //companionId, chatId
  Users: [] as ChatUser[],
}

export const messages = {
  list: [] as Message[],
}


export const currentChat = {
  username: "",
  userId: -1,
  chatId: -1
}


const getChatUsers = async () => {
  Object.assign(chatList, {
    Ids: new Map<number, number>(),
    Users: [],
  })

  await sendWsObject({
    type: "get",
    item: {
      url: "/chat/all",
    },
  })
  const userIds: iterator = await getUserIds()

  for (const id of Object.values(userIds)) {
    if (id !== userInfo.Id) {
      chatList.Users.push(await getUserById(id))
    }
  }

  for (const user of chatList.Users) {
    if (!chatList.Ids.has(user.Id)) {
      await sendWsObject({
        type: "post",
        item: {
          url: "/chat/create",
          data: {
            userId: user.Id,
          },
        },
      }) 
    }
  }
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
  const chatArea = document.getElementById("chat-area") as HTMLDivElement
  const chat = createElement("div", "chat show-chat")

  const chatId = chatList.Ids.get(id) as number
  await getMessages(chatId)
  await getUserById(id.toString())
  .then((user) => {
    Object.assign(currentChat, {
      username: user.Name,
      userId: user.Id,
      chatId: chatId,
    })
    console.log(currentChat)

    const chatName = createElement(
      "div",
      "chat-name",
      null,
      user.Name
    ) as HTMLDivElement

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

    chatForm.append(messageField, messageSend)
    chatFormContainer.appendChild(chatForm)
    // if (messages.list.length === 0) {
    //   reject()
    // }
    setTimeout(() => {
      for (const message of Object.values(messages.list)) {
        const msgElement = createElement("div", "message") as HTMLDivElement
        const date = new Date(message.timestamp)
        const formatDate = date
        .toLocaleString("en-GB", { timeZone: "EET" })
        .slice(0, 10)
        if (message.authorId === userInfo.Id) {
          msgElement.classList.add("send")
          msgElement.textContent = userInfo.Name + ":\n" + message.content + "\n" + formatDate
        } else if (message.authorId === -1) {
          msgElement.classList.add("status")
          msgElement.textContent = message.content + "\n" + formatDate
        } else {
          msgElement.classList.add("recieve")
          msgElement.textContent = currentChat.username + ":\n" + message.content + "\n" + formatDate
        }
        chatMessages.appendChild(msgElement)
      }

      chat.append(chatName, chatMessages, chatFormContainer)
      chatArea.replaceChildren(chat)
    }, 100)

  })
}

const sendMessage = async (chatId: number) => {
  const content = document.getElementById("chat-text") as HTMLInputElement
  let message = content.value.toString().trim()
  if (message.length > 0) {
    content.value = ""
    if (message.length > 150) {
      message = message.slice(0, 150)
    }
    await sendWsObject({
      type: "post",
      item: {
        url: `/chat/${chatId}/message`,
        data: {
          content: message,
        },
      },
    })
  }
}

export const updateChat = async (chatId: number) => {
  const message = messages.list[0]
  const chatMessages = document.getElementById(
    `Chat${chatId}`
  ) as HTMLDivElement
  const msgElement = createElement(
    "div",
    "message",
    null,
    message.content
  ) as HTMLDivElement
    console.log(message.id, messages.list[messages.list.length - 1].id)

  const date = new Date(message.timestamp)
  const formatDate = date
  .toLocaleString("en-GB", { timeZone: "EET" })
  .slice(0, 10)
  if (message.authorId === userInfo.Id) {
    msgElement.classList.add("send")
    msgElement.textContent = userInfo.Name + ":\n" + message.content + "\n" + formatDate
  } else if (message.authorId === -1) {
    msgElement.classList.add("status")
    msgElement.textContent = message.content + "\n" + formatDate
  } else {
    msgElement.textContent = currentChat.username + ":\n" + message.content + "\n" + formatDate
    msgElement.classList.add("recieve")
  }
  if (chatId === currentChat.chatId) {
    chatMessages.prepend(msgElement)
  }
}

const getMessages = async (chatId: number): Promise<void> => {
  Object.assign(messages, { list: [] })
  await sendWsObject({
    type: "get",
    item: {
      url: `/chat/${chatId}/messages`,
    },
  })
}

export const displayChatUsers = async () => {
  const onlineList = document.getElementById(
    "online-users"
  ) as HTMLUListElement,
  offlineList = document.getElementById("offline-users") as HTMLUListElement,
  onlineTitle = document.getElementById("online-title") as HTMLElement,
  offlineTitle = document.getElementById("offline-title") as HTMLElement

  /* Add eventlisteners to show/hide users in chat list */
  onlineTitle.addEventListener("click", toggleOnline)
  offlineTitle.addEventListener("click", toggleOffline)

  await getChatUsers()

  await sendWsObject({
    type: "get",
    item: {
      url: "/chat/all",
    },
  })

  const chatUsers = document.querySelectorAll(
    ".chat-user"
  ) as NodeListOf<HTMLElement>
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
  const offlineUsers = document.getElementById(
    "offline-users"
  ) as HTMLUListElement
  if (!offlineToggle) return

  if (offlineToggle.className === "bx bx-chevron-down") {
    offlineToggle.className = "bx bx-chevron-right"
    offlineUsers.style.display = "none"
  } else {
    offlineToggle.className = "bx bx-chevron-down"
    offlineUsers.style.display = "block"
  }
}
