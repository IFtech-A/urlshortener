import { useState } from "react";
import { SERVER_HOST } from "../api/consts";
import { Button, Card, Col, Input, Row } from "antd";
import { status as selectURLStatus } from "../store/urls/urlSlice";
import { useSelector } from "react-redux";
import { LOADING } from "../store/consts";

const URL = ({ onEnter }) => {
  const [urlValue, setUrlValue] = useState("");
  const [shortURL, setShortURL] = useState("");

  const status = useSelector(selectURLStatus)

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
    onEnter({
      url: urlValue,
      shortURL: shortURL,
    });
    setUrlValue("");
  };

  return (
    <Card style={{
      display: "flex",
      justifyContent: "center",
      flexDirection: "column",
      alignContent:"center",
      maxWidth: "500px",
      width:'100%'
    }}>
      <Row gutter={[0, 24]}>
        
        <Col span={24}>
          <label htmlFor="fullUrl">Your URL</label>
          <Input
            autoFocus
            size="large"
            id="fullUrl"
            placeholder="Enter a url"
            type="text"
            value={urlValue}
            onChange={(e) => setUrlValue(e.target.value)}
            onKeyDown={onKeyDown}
            width={'100%'}
          />
        </Col>
        <Col span={24}>
          <label htmlFor="shortURL">Short URL</label>
          <Input
            id="shortURL"
            size="large"
            width={'100%'}
            addonBefore={SERVER_HOST + "/"}
            onChange={(e) => setShortURL(e.target.value)}
          />
        </Col>
        <Button type="primary" loading={status === LOADING} size="large" onClick={()=> onSubmit(urlValue)}>
          Shorten!
        </Button>
      </Row>

    </Card>
  );
};

export default URL;
