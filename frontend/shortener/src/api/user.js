import { API_ENDPOINT } from "./consts";

const USER_ENDPOINT = API_ENDPOINT + "/user";

const [LOGIN_ACTION, SIGNUP_ACTION] = ["login", "signup"];

export const login = async (creds) => {
  creds.action = LOGIN_ACTION;
  try {
    const response = await fetch(API_ENDPOINT + '/login', {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify(creds),
    });
    const jsonData = await response.json();
    return jsonData;
  } catch (e) {
    console.error(e);
    throw e;
  }
};

export const signup = async (creds) => {
  creds.action = SIGNUP_ACTION;
  try {
    const response = await fetch(USER_ENDPOINT, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify(creds),
    });
    const jsonData = await response.json();
    return jsonData;
  } catch (e) {
    console.error(e);
    throw e;
  }
};
