import {
  createAsyncThunk,
  createSlice,
} from "@reduxjs/toolkit";
import { createShortURL, getUrlHistory } from "../../api/api";
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
    const url = await createShortURL(initialURL);

    return url;
  }
);

export const fetchURLs = createAsyncThunk("urls/fetchURLs", async () => {
  const urls = await getUrlHistory();

  return urls;
});

const urlSlice = createSlice({
  name: "urls",
  initialState,
  reducers: {
    urlAdded(state, action) {
      state.status = IDLE
      state.createStatus = IDLE
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
        console.log({state,action})
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
export default urlSlice.reducer;
