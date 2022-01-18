import { API_ENDPOINT } from "./consts";

class Endpoints {
  api = "";
  constructor(api) {
    this.api = api;
  }

  get user() {
    return this.api + "/users";
  }
  userById(id) {
    return `${this.user}/${id}`;
  }
  userLinks(id) {
    return `${this.user}/${id}/links`;
  }
  get link() {
    return this.api + "/links";
  }
  linkById(id) {
    return `${this.link}/${id}`;
  }
  get login() {
    return this.api + "/login";
  }
  get logout() {
    return this.api + "/logout";
  }
}

class SDK {
  user = {
    id: "",
    firstname: "",
    lastname: "",
    email: "",
    username: "",
    password: "",
  };
  #token = "";

  #headers = {
    "Content-Type": "application/json",
    Accept: "application/json",
  };

  api = new Endpoints();

  constructor(api) {
    this.api = new Endpoints(api);
  }

  async authenticate(username, password) {
    const response = await fetch(this.api.login, {
      method: "POST",
      headers: this.headersWithAuth(),
      body: JSON.stringify({
        username,
        password,
      }),
    });

    if (response != null && response.ok) {
      const data = await response.json();
      this.user = data.user;
      this.#token = data.token;
    } else {
      throw new Error("Authentication failed");
    }

    return this;
  }

  logout() {
    this.user = {
      id: "",
    };
    this.#token = "";
  }

  headersWithAuth() {
    if (typeof this.#token === "string" && this.#token.length > 0) {
      return {
        Authorization: "Bearer " + this.#token,
        ...this.#headers,
      };
    }
    return this.#headers;
  }

  apiFetch(url, method, body, headers) {
    return fetch(url, {
      method,
      headers: {
        ...headers,
        ...this.headersWithAuth(),
      },
      credentials: "include",
      body: body ? JSON.stringify(body) : null,
    })
      .then((response) => {
        if (response.ok) {
          const contentType = response.headers.get("Content-Type") || "";

          if (contentType.includes("application/json")) {
            return response.json().catch((error) => {
              return Promise.reject(
                new Error("Invalid JSON: " + error.message)
              );
            });
          }

          if (contentType.includes("text/html")) {
            return response
              .text()
              .then((html) => {
                return {
                  page_type: "generic",
                  html: html,
                };
              })
              .catch((error) => {
                return Promise.reject(
                  new Error("HTML error: " + error.message)
                );
              });
          }

          return Promise.reject(
            new Error("Invalid content type: " + contentType)
          );
        }

        if (response.status === 404) {
          return Promise.reject(new Error("Page not found: " + url));
        }

        if (response.status === 403) {
          this.#token = "";
          return Promise.reject(new Error("API authentication failed"));
        }

        return response.json().then((res) => {
          // if the response is ok but the server rejected the request, e.g. because of a wrong password, we want to display the reason
          // the information is contained in the json()
          // there may be more than one error
          let errors = [];
          Object.keys(res).forEach((key) => {
            errors.push(`${key}: ${res[key]}`);
          });
          return Promise.reject(new Error(errors));
        });
      })
      .catch((error) => {
        return Promise.reject(error.message);
      });
  }

  fetchGet(uri, headers = {}) {
    return this.apiFetch(uri, "GET", null, headers);
  }
  fetchPost(uri, bodyData, headers = {}) {
    return this.apiFetch(uri, "POST", bodyData, headers);
  }

  createLink(link, id = this.user.id) {
    if (id === "") {
      return this.fetchPost(this.api.link, link);
    } else {
      return this.fetchPost(this.api.userLinks(id), link);
    }
  }

  createUser(user) {
    return this.fetchPost(this.api.user, user);
  }

  readUser(id = this.user.id) {
    return this.fetchGet(this.api.userById(id));
  }

  readUserLinks(id = this.user.id) {
    if (id === "") {
      return this.fetchGet(this.api.link);
    } else {
      return this.fetchGet(this.api.userLinks(id));
    }
  }

  readLink(id) {
    return this.fetchGet(this.api.linkById(id));
  }
}

const sdk = new SDK(API_ENDPOINT);

export default sdk;
