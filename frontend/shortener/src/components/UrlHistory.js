import { Card, List, Typography } from "antd";
import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { FAILED, IDLE, SUCCEEDED } from "../store/consts";
import {
  fetchURLs,
  status as urlFetchStatus,
  urls as allURLs,
} from "../store/urls/urlSlice";
import UrlHistoryItem from "./UrlHistoryItem";

const UrlHistory = () => {
  const urls = useSelector(allURLs);
  const status = useSelector(urlFetchStatus);
  const dispatch = useDispatch();

  useEffect(() => {
    if (status === IDLE) {
      dispatch(fetchURLs());
    }
  }, [status, dispatch]);

  return (
    <Card
      style={{
        maxWidth: "100%",
      }}
    >
      {status === SUCCEEDED ? (
        <List
          size="large"
          header={<Typography.Title level={2}>URL history</Typography.Title>}
          dataSource={urls}
          renderItem={(item) => <UrlHistoryItem url={item} />}
        />
      ) : status === FAILED ? (
        <Typography.Title type="danger" level={2}>Failed to load url history</Typography.Title>
      ) : (
        <Typography.Title type="secondary" level={2}>Loading url history</Typography.Title>
      )}
    </Card>
  );
};

export default UrlHistory;
