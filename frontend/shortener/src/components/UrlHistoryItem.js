import React from "react";
import { List } from "antd";
import { SERVER_HOST } from "../api/consts";

const UrlHistoryItem = ({url}) => {
  return (
    <List.Item>
      <div key={url.shortened}>
        <a style={{ fontSize: 18 }} href={SERVER_HOST + "/" + url.shortened}>
          {SERVER_HOST + "/" + url.shortened}
        </a>
        <p style={{ fontSize: 12 }}>{url.real}</p>
      </div>
    </List.Item>
  );
};

export default UrlHistoryItem;
