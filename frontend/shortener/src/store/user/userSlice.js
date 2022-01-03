import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { signup } from "../../api/user";
import { IDLE, LOADING, SUCCEEDED } from "../consts";

const initialState = {
  user: null,
  status: IDLE,
  error: null,
};

const signupAsyncThunk = createAsyncThunk("user/signup", async (creds) => {
  return await signup(creds);
});

const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    login(state, action) {
      state.user = action.payload;
      console.log("Logged in");
    },
    logout(state, action) {
      state.user = null;
      console.log("Logged out");
    },
  },
  extraReducers(builder) {
    builder
      .addCase(signupAsyncThunk.pending, (state, action) => {
        state.status = LOADING;
      })
      .addCase(signupAsyncThunk.fulfilled, (state, action) => {
        state.status = SUCCEEDED;
      });
  },
});

export const user = (state) => state.user.user;
export const {login, logout} = userSlice.actions
export default userSlice.reducer;
