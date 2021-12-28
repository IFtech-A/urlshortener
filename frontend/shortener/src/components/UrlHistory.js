import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { SERVER_HOST } from "../api/consts";
import { FAILED, IDLE, SUCCEEDED } from "../store/consts";
import {
  fetchURLs,
  status as urlFetchStatus,
  urls as allURLs,
} from "../store/urls/urlSlice";

const UrlHistory = () => {
  const urls = useSelector(allURLs);
  const status = useSelector(urlFetchStatus);
  const dispatch = useDispatch();

  useEffect(() => {
    if (status === IDLE) {
      dispatch(fetchURLs());
    }
  }, [status, dispatch]);

  let urlHistory = null;
  switch (status) {
    case SUCCEEDED:
      if (urls?.length !== 0) {
        urlHistory = urls.map((url) => (
          <div key={url.shortened}>
            <a
              style={{ fontSize: 18 }}
              href={SERVER_HOST + "/" + url.shortened}
            >
              {SERVER_HOST + "/" + url.shortened}
            </a>
            <p style={{ fontSize: 12 }}>{url.real}</p>
          </div>
        ));
      } else {
        urlHistory = <p>You don't have any history of url shortenings</p>;
      }
      break;
    case FAILED:
      urlHistory = <p>Failed to load url history</p>;
      break;
    default:
      urlHistory = <div>Loading url history</div>;
  }

  return (
    <div style={{ fontFamily: "inherit" }}>
      <h2>Your URL history:</h2>
      {urlHistory}
    </div>
  );
};

export default UrlHistory;
