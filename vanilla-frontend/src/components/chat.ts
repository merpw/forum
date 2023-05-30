import { ActiveUser, InactiveUser } from "./types"
import { ws } from "./ws.js"
// This file is dedicated to sorting and displaying chat users in the sidebar.

const chatUsers = {
  active:
    [] as ActiveUser[] /* Active users are users you have a chat history with */,
  inactive:
    [] as InactiveUser[] /* Inactive users are users you don't have a chat history with */,
}

function getChatUsers() {
  /* Reset the state of the chatUsers upon relogging */
  Object.assign(chatUsers, { active: [], inactive: [] })

  const testUserList: ActiveUser[] = []
  const testUserListInactive: InactiveUser[] = []
  const activeUser1: ActiveUser = {
    Name: "Test",
    ID: 1,
    Online: true,
    UnreadMSG: true,
  }

  const activeUser2: ActiveUser = {
    Name: "Test2",
    ID: 2,
    Online: false,
    UnreadMSG: false,
  }

  const inactiveUser1: InactiveUser = {
    Name: "InactiveTest",
    ID: 1,
    Online: true,
  }

  const inactiveUser2: InactiveUser = {
    Name: "InactiveTest2",
    ID: 2,
    Online: false,
  }

  testUserList.push(activeUser1, activeUser2)
  testUserListInactive.push(inactiveUser1, inactiveUser2)

  for (const user of testUserList) {
    chatUsers.active.push(user)
  }

  for (const user of testUserListInactive) {
    chatUsers.inactive.push(user)
  }
}

export const displayChatUsers = () => {
  const onlineList = document.getElementById("online-users") as HTMLUListElement
  const offlineList = document.getElementById(
    "offline-users"
  ) as HTMLUListElement
  const onlineTitle = document.getElementById("online-title") as HTMLElement
  const offlineTitle = document.getElementById("offline-title") as HTMLElement
  if (!onlineList || !offlineList || !onlineTitle || !offlineTitle) return
  /* Add eventlisteners to show/hide users in chat list */
  onlineTitle.addEventListener("click", toggleOnline)
  offlineTitle.addEventListener("click", toggleOffline)
  getChatUsers()

  /* This loop gets all the ACTIVE users and appends them to chatlist */
  for (const user of chatUsers.active) {
    const newElement = document.createElement("li")
    const userName = document.createElement("p")
    userName.id = user.ID.toString()
    userName.textContent = `${user.Name} `
    newElement.appendChild(userName)

    if (user.UnreadMSG) {
      const unreadElement = document.createElement("i")
      unreadElement.className = "bx bx-message-dots"
      newElement.appendChild(unreadElement)
    }

    if (user.Online) {
      newElement.className = "online"
      onlineList.appendChild(newElement)
    } else {
      newElement.className = "offline"
      offlineList.appendChild(newElement)
    }
  }

  /* This loop gets all the INACTIVE users and appends them to chatlist */
  for (const user of chatUsers.inactive) {
    const newElement = document.createElement("li")
    const userName = document.createElement("p")
    userName.id = user.ID.toString()
    userName.textContent = `${user.Name} `
    newElement.appendChild(userName)
    if (user.Online) {
      newElement.className = "online"
      onlineList.appendChild(newElement)
    } else {
      newElement.className = "offline"
      offlineList.appendChild(newElement)
    }
  }
}

const toggleOnline = () => {
  const onlineToggle = document.getElementById("online-toggle") as HTMLElement
  const onlineUsers = document.getElementById(
    "online-users"
  ) as HTMLUListElement
  if (!onlineToggle || !onlineUsers) return

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

export const sendMessage = () => {
  const chatMsg = document.getElementById("chat-text") as HTMLInputElement
  ws.send(JSON.stringify(chatMsg.value))
  chatMsg.value = ""
  return
}

export const displayMessage = () => {
  const messageDisplay = document.getElementById("chat-messages") as HTMLDivElement
  const message = [] as HTMLDivElement []
  for (const message of messages) {
    
  }
}


const createMessage = async (author: string, userId: number, content: string, sender: boolean,): Promise<HTMLDivElement> => {
  const message = document.createElement("div")
  if(sender) {
    message.className = "message sender"
  } else {
    message.className = "message reciever"
  }
  const info = document.createElement("div")
  info.id = `MUID${userId.toString()}` // Message User Id
  info.textContent = `${author}`
  const contentElement = `${content}`
  message.append(info, contentElement)
  return message
}
