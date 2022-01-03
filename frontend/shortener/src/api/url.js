import { API_ENDPOINT, SERVER_HOST } from "./consts";

const URL_ENDPOINT = API_ENDPOINT + "/url";

export const getUrlHistory = async () => {
  try {
    const response = await fetch(URL_ENDPOINT, {
      method: "GET",
      credentials: "include",
    });
    const responseJSON = await response.json();
    return responseJSON;
  } catch (e) {
    console.error(e);
    throw e;
  }
};

export const createShortURL = async (url) => {
  try {
    const response = await fetch(URL_ENDPOINT, {
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      method: "POST",
      credentials: "include",
      body: JSON.stringify({
        real: url.url,
        short: url.shortURL,
      }),
    });
    const responseJSON = await response.json();
    console.log(responseJSON);
    console.log(SERVER_HOST + "/" + responseJSON.shortened);
    return SERVER_HOST + "/" + responseJSON.shortened;
  } catch (e) {
    alert(e);
    console.error(e);
    throw e;
  }
};
