import {useEffect, useState} from 'react'
import {getUrlHistory, SERVER_HOST} from '../api/api'


const UrlHistory = () => {
    // const [historyList, SetHistoryList] = useState([])
    // const [loaded, setLoaded] = useState(false)
    // const [loading, setLoading] = useState(false)
    const [state, setState] = useState({
        loaded: false,
        historyList: [],
    })

    const fetchHistory = async ()=>{
        let urls = await getUrlHistory()
        if (urls !== undefined) {
            console.log(urls)
            // SetHistoryList([...urls])
            setState({loaded: true, historyList: [...urls.reverse()]})
        } else {
            setState({loaded: true, historyList: []})
        }
    }
    
    useEffect(()=> fetchHistory(), [])
        
    return (
        <div style={{fontFamily:'inherit'}}>
            <h2>Your URL history:</h2>
            {!state.loaded ? <p>Loading</p> : 
            state.historyList.length === 0 ?  <p>You don't have any history of url shortenings</p> :
            state.historyList.map(url => 
                <div key={url.shortened}>
                    <a style={{fontSize:18}} href={SERVER_HOST + "/" +url.shortened}>
                    {SERVER_HOST + "/" +url.shortened}
                    </a>
                    <p style={{fontSize:12}}>{url.real}</p>
                </div>
            )} 
        </div>
    )
}

export default UrlHistory
