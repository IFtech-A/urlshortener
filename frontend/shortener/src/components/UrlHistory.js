import { List } from "antd";
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
    <div style={{ fontFamily: "inherit" }}>
      {status === SUCCEEDED ? (
        <List
          size="large"
          header={<h2>URL history</h2>}
          dataSource={urls}
          renderItem={(item) => <UrlHistoryItem url={item} />}
        />
      ) : status === FAILED ? (
        <p>Failed to load url history</p>
      ) : (
        <p>Loading url history</p>
      )}
    </div>
  );
};

export default UrlHistory;
