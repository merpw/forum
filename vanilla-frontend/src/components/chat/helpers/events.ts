const startChat = new Event("startChat")
const chatCreated = new Event("chatCreated")
const updateClientState = new Event("updateClientState")
const renderChatList = new Event("renderChatList")
const renderChatMessages = new Event("renderChatMessages")
const renderNewMessages = new Event("renderNewMessages")

export {
  startChat,
  chatCreated,
  updateClientState,
  renderChatList,
  renderChatMessages,
  renderNewMessages,
}
