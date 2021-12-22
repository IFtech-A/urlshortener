import {configureStore} from '@reduxjs/toolkit'
import urlReducer from './urls/urlSlice'

export default configureStore({
    reducer: {
        urls: urlReducer
    }
})