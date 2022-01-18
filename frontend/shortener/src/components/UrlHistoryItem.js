import React from "react";
import { List, Space } from "antd";
import { SERVER_HOST } from "../api/consts";
import { Typography, Row } from "antd";

const { Text, Link } = Typography;
const UrlHistoryItem = ({ url }) => {
  const shortURL = SERVER_HOST + "/" + url.short;
  return (
    <List.Item style={{ width: "100%" }}>
      <Space style={{ width: "100%" }} direction="vertical" key={url.short}>
        <Link
          copyable
          ellipsis
          style={{ fontSize: 24, display:'block' }}
          href={shortURL}
          target="_blank"
        >
          {SERVER_HOST + "/" + url.short}
        </Link>
        <Row justify="space-between">
          <Text ellipsis style={{ fontSize: 14 }}>
            {url.real}
          </Text>
          <Text ellipsis style={{ fontSize: 14 }}>
            {url.created_at}
          </Text>
        </Row>
      </Space>
    </List.Item>
  );
};

export default UrlHistoryItem;
