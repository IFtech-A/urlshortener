import { useState } from "react";

const URL = ({ onEnter }) => {
  const [urlValue, setUrlValue] = useState("");
  const [shortened, setShortened] = useState("");

  const onKeyDown = (e) => {
    if (e.keyCode === 13) {
      onSubmit();
    }
  };

  const onSubmit = () => {
    if (urlValue === "") {
      alert("empty field");
      return;
    }
    let short = onEnter(urlValue);
    if (short !== undefined) {
      setShortened(short);
    }
    setUrlValue("");
  };

  return (
    <div>
      <div style={{ display: "flex", justifyContent: "space-between" }}>
        <div>
          <label for="fullUrl">Input the URL here : </label>
          <input
            id="fullUrl"
            placeholder="Enter a url"
            type="text"
            value={urlValue}
            onChange={(e) => setUrlValue(e.target.value)}
            onKeyDown={onKeyDown}
          />
        </div>
        <button className="btn" onClick={() => onSubmit(urlValue)}>
          Shorten!
        </button>
      </div>
      {shortened && shortened.length > 0 && <h3>{shortened}</h3>}
    </div>
  );
};

export default URL;
