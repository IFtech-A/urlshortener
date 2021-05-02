import React from 'react'

const URL = ({onEnter}) => {
    const onKeyDown = (e) => {
        if (e.keyCode === 13) {
            onEnter();
            document.getElementById("fullUrl").value = "";
        }
        
    }
    return (
        <div>
            <label for="fullUrl">Input the URL here : </label>
            <input id="fullUrl" placeholder="Enter a url" type="text" onKeyDown={onKeyDown} />
        </div>
    )
}

export default URL
