import 'whatwg-fetch';
import * as Helpers from './APIHelpers';

/* NewUser should be a object that has:
  - Email
  - Password
  - PasswordConf
  - UserName
  - FirstName
  - LastName
 */
export const signUpNewUser = (newUser) => {
    let signUpUserPath = `${Helpers.API_PATH}${Helpers.USERS_PATH}`;
    let requestPackage = {
        method: Helpers.METHOD_POST,
        body: JSON.stringify(newUser),
    };
    return fetch(signUpUserPath, requestPackage)
        .then((response) => {
            if (response.status < 300) {
                return response
            }
            throw new Error(`error signing up new account with status code ${response.status}`)
        })
        .catch((error) => {
            console.error(`Ran into error signing up: ${error}`);
        })
};

export const signInUser = (credentials) => {
    let signInUserPath = `${Helpers.API_PATH}${Helpers.SESSIONS_PATH}`;
    let requestPackage = {
        method: Helpers.METHOD_POST,
        body: JSON.stringify(credentials),
    };
    return fetch(signInUserPath, requestPackage)
        .then((response) => {
            let token = '';
            if (response.status < 300) {
                token = response.headers.get("Authorization");
                localStorage.setItem(Helpers.TOKEN_KEY, token);
                localStorage.setItem(Helpers.TOKEN_KEY_CREATED, new Date().toString());
            }
            return token;
        })
        .catch((error) => {
            console.error(`Ran into error signing in: ${error}`);
        })
};

export const signOutUser = () => {
    localStorage.removeItem(Helpers.TOKEN_KEY_CREATED);
    localStorage.removeItem(Helpers.TOKEN_KEY);
    let signOutPath = `${Helpers.API_PATH}${Helpers.SESSIONS_PATH_MINE}`;
    let requestPackage = {
        method: Helpers.METHOD_DELETE,
        headers: {
            "Authorization": localStorage.getItem(Helpers.TOKEN_KEY)
        }
    };
    return fetch(signOutPath, requestPackage)
        .then((response) => {
            if (response.status < 300) {
                return true;
            }
            throw new Error(`HTTP status was ${response.status}`)
        })
        .catch((error) => {
            console.error(`Ran into error signing out: ${error}`);
        })
};

export const getUser = () => {
    let userPath = `${Helpers.API_PATH}${Helpers.USERS_ME_PATH}`;
    let requestPackage = {
        method: "GET",
        headers: {
            "Authorization": localStorage.getItem(Helpers.TOKEN_KEY)
        }
    };
    return fetch(userPath, requestPackage)
        .then((response) => {
            if (response.status < 300) {
                return response.json()
            }
            throw new Error(`Fetch failed: status code recieved was: ${response.status}`)
        })
        .catch((error) => {
            console.error(`Ran into error signing out: ${error}`);
        })
};

export const updateUser = (userSettingsUpdates) => {
    let userPath = `${Helpers.API_PATH}${Helpers.USERS_ME_PATH}`;
    let requestPackage = {
        method: "PATCH",
        body: JSON.stringify(userSettingsUpdates),
        headers: {
            "Authorization": localStorage.getItem(Helpers.TOKEN_KEY)
        }
    };
    return fetch(userPath, requestPackage)
        .then((response) => {
            if (response.status < 300) {
                return response.json();
            }
            throw new Error(`Fetch failed. Error updating user. Status code was: ${response.status}`);
        })
        .catch((error) => {
            console.log(`error updating user: ${error}`);
        });
};

export const getUsersFromPrefix = (prefix) => {
    let userPath = `${Helpers.API_PATH}${Helpers.USERS_PATH}?q=${prefix}`;
    let requestPackage = {
        method: "GET",
        headers: {
            "Authorization": localStorage.getItem(Helpers.TOKEN_KEY)
        }
    };
    return fetch(userPath, requestPackage)
        .then((response) => {
            if (response.status < 300) {
                return response.json();
            }
            throw new Error(`Fetch failed. Error updating user. Status code was: ${response.status}`);
        })
        .catch((error) => {
            console.log(`error updating user: ${error}`);
        });
};