import { createSlice, PayloadAction } from "@reduxjs/toolkit"

const initialState: {
  status: "connecting" | "connected" | "disconnected"
  pendingRequests: string[]
} = {
  status: "disconnected",
  pendingRequests: [],
}

const wsConnectionSlice = createSlice({
  name: "wsConnection",
  initialState: initialState,
  reducers: {
    connectionStarted: (state) => {
      state.status = "connecting"
    },
    connectionEstablished: (state) => {
      state.status = "connected"
    },
    connectionClosed: (state) => {
      state.status = "disconnected"
    },
    requestPending: (state, action: PayloadAction<string>) => {
      state.pendingRequests.push(action.payload)
    },
    requestDone: (state, action: PayloadAction<string>) => {
      state.pendingRequests = state.pendingRequests.filter((url) => url !== action.payload)
    },
  },
})

export const wsConnectionActions = wsConnectionSlice.actions

export default wsConnectionSlice.reducer
