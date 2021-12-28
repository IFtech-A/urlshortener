import {configureStore} from '@reduxjs/toolkit'
import urlReducer from './urls/urlSlice'
import userReducer from './user/userSlice'

export default configureStore({
    reducer: {
        urls: urlReducer,
        user: userReducer
    }
})