import { configureStore } from "@reduxjs/toolkit"

import chatsReducer from "@/store/chats/chats"
import wsConnectionMiddleware from "@/store/wsMiddleware"
import wsConnectionReducer from "@/store/ws"

const store = configureStore({
  reducer: {
    chats: chatsReducer,
    wsConnection: wsConnectionReducer,
  },
  middleware: [wsConnectionMiddleware],
})

export default store

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>
// Inferred type: {posts: PostsState, comments: CommentsState, users: UsersState}
export type AppDispatch = typeof store.dispatch
