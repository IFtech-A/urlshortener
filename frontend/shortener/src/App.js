import Header from './components/Header'
import URL from './components/URL'


function App() {
  const onUrlShortenClick = () => {
    
    console.log("shorten API call")
  }

  return (
    <div className="container">
      <Header onShorten={onUrlShortenClick} />
      <URL onEnter={onUrlShortenClick} />
    </div>
  );
}

export default App;
