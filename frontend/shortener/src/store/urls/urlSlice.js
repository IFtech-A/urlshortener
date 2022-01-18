import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import sdk from "../../api";
import { IDLE, LOADING, FAILED, SUCCEEDED } from "../consts";

const initialState = {
  urls: [],
  status: IDLE,
  createStatus: IDLE,
  error: null,
};

export const addNewURL = createAsyncThunk(
  "urls/addNewURL",
  async (initialURL, thunkAPI) => {
    const url = await sdk.createLink(initialURL);
    console.log("url addnewurl thunk", url);

    return url;
  }
);

export const fetchURLs = createAsyncThunk("urls/fetchURLs", async () => {
  const urls = await sdk.readUserLinks();

  return urls;
});

const urlSlice = createSlice({
  name: "urls",
  initialState,
  reducers: {
    urlAdded(state, action) {
      state.createStatus = IDLE;
    },
    refreshUrls(state, action) {
      state.status = IDLE;
    },
  },
  extraReducers(builder) {
    builder
      .addCase(addNewURL.pending, (state, action) => {
        state.createStatus = LOADING;
      })
      .addCase(addNewURL.fulfilled, (state, action) => {
        state.createStatus = SUCCEEDED;
        state.status = IDLE;
      })
      .addCase(addNewURL.rejected, (state, action) => {
        state.createStatus = FAILED;
        state.error = action.error.message;
      })
      .addCase(fetchURLs.fulfilled, (state, action) => {
        state.status = SUCCEEDED;
        state.urls = action.payload;
      })
      .addCase(fetchURLs.pending, (state, action) => {
        state.status = LOADING;
      })
      .addCase(fetchURLs.rejected, (state, action) => {
        state.status = FAILED;
        state.error = action.error.message;
      });
  },
});

export const status = (state) => state.urls.status;
export const urls = (state) => state.urls.urls;
export const error = (state) => state.urls.error;
export const createStatus = (state) => state.urls.createStatus;
export const { urlAdded, refreshUrls } = urlSlice.actions;
export default urlSlice.reducer;
