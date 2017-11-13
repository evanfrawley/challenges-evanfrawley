// export const API_PATH = "";
console.log("env", process.env.NODE_ENV);
export const API_PATH = process.env.NODE_ENV === 'development' ? "http://localhost:4000" : "https://api.evan.gg";
export const USERS_PATH = "/v1/users";
export const USERS_ME_PATH = `${USERS_PATH}/me`;
export const SUMMARY_PATH = "/v1/summary";
export const SESSIONS_PATH = "/v1/sessions";
export const SESSIONS_PATH_MINE = `${SESSIONS_PATH}/mine`;

export const METHOD_POST = "POST";
export const METHOD_PATCH = "PATCH";
export const METHOD_DELETE = "DELETE";

export const TOKEN_KEY = "344_messaging_token";
export const TOKEN_KEY_CREATED = "344_messaging_token_created_at";