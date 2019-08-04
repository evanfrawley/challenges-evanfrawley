// export const API_PATH = "";
console.log("env", process.env.NODE_ENV);
export const API_PATH = process.env.NODE_ENV === 'development' ? "http://localhost:4000" : "https://api.evan.gg";
export const USERS_PATH = "/v1/users";
export const USERS_ME_PATH = `${USERS_PATH}/me`;
export const SUMMARY_PATH = "/v1/summary";
export const SESSIONS_PATH = "/v1/sessions";
export const SESSIONS_PATH_MINE = `${SESSIONS_PATH}/mine`;

export const METHOD_GET = "GET";
export const METHOD_POST = "POST";
export const METHOD_PATCH = "PATCH";
export const METHOD_DELETE = "DELETE";

export const TOKEN_KEY = "344_messaging_token";
export const TOKEN_KEY_CREATED = "344_messaging_token_created_at";

export const CHANNELS_PATH = "/v1/channels";
export const MESSAGES_PATH = "/v1/messages";

// Headers
const headerContentType = "Content-Type";
const contentTypeJSON = "application/json";
const headerAuthorization = "Authorization";

export const createAndSendGet = (url, headers) => {
  // do a for loop over headers
  return createAndSendRequest(url, null, headers, METHOD_GET);
};

export const createAndSendPost = (url, body, headers) => {
  // do a for loop over headers
  return createAndSendRequest(url, body, headers, METHOD_POST);
};

export const createAndSendPatch = (url, updates, headers) => {
  // do a for loop over headers
  return createAndSendRequest(url, updates, headers, METHOD_PATCH);
};

export const createAndSendDelete = (url, headers) => {
  // do a for loop over headers
  return createAndSendRequest(url, null, headers, METHOD_DELETE);
};

const createAndSendRequest = (url, body, headers, method) => {
  let requestHeaders = {};
  requestHeaders[headerContentType] = contentTypeJSON;
  requestHeaders[headerAuthorization] = localStorage.getItem(TOKEN_KEY);

  // set necessary request headers
  if (headers) {
    headers.keys().forEach((key) => {
      requestHeaders[key] = headers[key];
    });
  }

  // create package to send in fetch
  let requestPackage = {
    method,
    headers: requestHeaders
  };

  // attach body if it exists
  if (body) {
    requestPackage.body = JSON.stringify(body);
  }

  return fetch(url, requestPackage)
    .then(res => {
      return res.json()
    });
};


