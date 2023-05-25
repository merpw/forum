import { ActiveUser, InactiveUser } from "./types"

// This file is dedicated to sorting and displaying chat users in the sidebar.

const chatUsers = {
  active:
    [] as ActiveUser[] /* Active users are users you have a chat history with */,
  inactive:
    [] as InactiveUser[] /* Inactive users are users you don't have a chat history with */,
}

function getChatUsers() {
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
  const offlineList = document.getElementById("offline-users") as HTMLUListElement
  if (!onlineList || !offlineList) return
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
