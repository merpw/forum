import { ChatUser, Message } from "../../types"
import { getUserById, getUserIds } from "../../api/get.js"

import { userInfo } from "../authorization/auth.js"

import { 
  sendWsObject,
  getMessageList, 
  toggleOfflineSection, 
  toggleOnlineSection, 
  lazyLoading, 
  hideChat,
  sendMessage,
  createMessage } from "./helpers.js"

  import { createElement, iterator } from "../utils.js"

// chatList keeps track of all the users on the database
export const chatList = {
  Ids: new Map<number, number>(), //companionId, chatId
  Users: [] as ChatUser[],
}

// messages is used for message-state handling
export const messages = {
  list: [] as Message[],
  ids: [] as number[],
  current: {} as Message 
}

// currentChat keeps track of the state of the current chat
export const currentChat = {
  username: "",
  userId: -1,
  chatId: -1,
  range: {
    min: 0,
    max: 10
  }
}


// Gets all chat users and assigns them both as IDs as users
// Also creates ALL chats if not already created
const getChatUsers = async () => {
  Object.assign(chatList, {
    Ids: new Map<number, number>(),
    Users: [],
  })

  const userIds: iterator = await getUserIds()
  for (const id of Object.values(userIds)) {
    if (id !== userInfo.Id) {
      chatList.Users.push(await getUserById(id))
    }
  }

  sendWsObject({
    type: "get",
    item: {
      url: "/users/online"
    }
  })

  for (const user of chatList.Users.reverse()) {
    if (!chatList.Ids.has(user.Id)) {
      sendWsObject({
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

// Display the current chat based on chatId
// It automatically removes message notification
const showChat = async (id: number) => {
  const chatArea = document.getElementById("chat-area") as HTMLDivElement
  const chat = createElement("div", "chat show-chat")

  for (const user of chatList.Users) {
    if (user.Id === id && user.UnreadMsg) {
      user.UnreadMsg = false
      renderChatList()
      break
    }
  }

  const chatId = chatList.Ids.get(id) as number
  await getMessages(chatId)
  getMessageList(messages.ids.reverse())
  const user = await getUserById(id.toString())

  Object.assign(currentChat, {
    username: user.Name,
    userId: user.Id,
    chatId: chatId,
    range: {
      min: 0,
      max: 10
    }
  })

  const chatHeader = createElement(
    "div",
    "chat-header"
  )


  const chatName   = createElement(
    "div",
    "chat-name",
    null,
    "Chatting with " + user.Name
  ) as HTMLDivElement
  
  const closeBtn = createElement("a", "closebtn", "chat-window-close", null, "<i class='bx bx-x'>")
  closeBtn.addEventListener("click", (e) => {
    e.preventDefault()
    hideChat()
  })
  chatHeader.append(chatName, closeBtn)

  const chatMessages = createElement("div", "chat-messages", `Chat${chatId}`)

  const chatFormContainer = createElement("div", "chat-form-container")
  const chatForm = createElement("form", "chat-form")
  chatForm.setAttribute("autocomplete", "off")
  const messageField = createElement("input", null, "chat-text")
  messageField.setAttribute("maxlength", "150")

  // Adding event listeners for messageEvents and lazy loading
  chatMessages.addEventListener("messageEvent", () => {
    updateChat(chatId)
  })

  chatMessages.addEventListener("scroll", () => {
    if(Math.abs(chatMessages.scrollTop) === (chatMessages.scrollHeight - chatMessages.clientHeight)) {
     lazyLoading(chatId) 
    }
  })

  const messageSend = createElement("button", null, "chat-send", "Send")
  messageSend.style.marginLeft = "4px"
 
  messageSend.addEventListener("click", (e) => {
    sendMessage(chatId)
    e.preventDefault()
  })

  chatForm.append(messageField, messageSend)
  chatFormContainer.appendChild(chatForm)

  setTimeout(() => {
    for (const message of Object.values(messages.list.reverse().splice(0, 10))) {
      chatMessages.append(createMessage(message))
    }

    chat.append(chatHeader, chatMessages, chatFormContainer)
    chatArea.replaceChildren(chat)
  }, 200)
}



// Sends a message to another chat user through the WebSocket server




// Updates the chat if you send or recieve a chat message
export const updateChat = (chatId: number): void => {
  if (currentChat.chatId !== messages.current.chatId) {
    return
  }
  const [message] = messages.list
  const chatMessages = document.getElementById(
    `Chat${chatId}`
  ) as HTMLDivElement
  chatMessages.prepend(createMessage(message))
  currentChat.range.min += 1
  currentChat.range.max += 1
}

// Gets the message list of provided chatId and populates the messages.list array with messages
const getMessages = async (chatId: number): Promise<void> => {
  return new Promise<void>((resolve) => {
    // Reset the state of the messages
    Object.assign(messages, { list: [] })
    setTimeout(() => {
      sendWsObject({
        type: "get",
        item: {
          url: `/chat/${chatId}/messages`,
        },
      })
    }, 5)
    setTimeout(resolve, 40)
  })
}

// Displays chat users in the chat list when opening the chat list
export const displayChatUsers = async () => {
  const onlineTitle = document.getElementById("online-title") as HTMLElement
  const offlineTitle = document.getElementById("offline-title") as HTMLElement

  /* Add eventlisteners to show/hide users in chat list */
  onlineTitle.addEventListener("click", toggleOnlineSection)
  offlineTitle.addEventListener("click", toggleOfflineSection)

  await getChatUsers()
  renderChatList()
}

// Renders the chat list with all the requirements
// (alphabetically, latest message sent, online/offline, unread message)
export function renderChatList(): void {
    sendWsObject({
      type: "get",
      item: {
        url: "/chat/all",
      },
    }) 

    setTimeout(() => {
      chatList.Users.sort((a, b) => {
        if (a.LastMessageId > b.LastMessageId) {
          return -1
        } else if (a.LastMessageId < b.LastMessageId) {
          return 1
        } else {
          return 0
        }
      })
    }, 110)

    setTimeout(() => {
      const onlineList = document.getElementById("online-users") as HTMLUListElement,
      offlineList = document.getElementById("offline-users") as HTMLUListElement

      // Resets the list if user logging in or out
      onlineList.replaceChildren()
      offlineList.replaceChildren()

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
    }, 120)
}


