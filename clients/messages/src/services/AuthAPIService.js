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
      return JSON.parse(response.body);
    })
    .catch((error) => {
      console.error(`Ran into error fetching data: ${error}`);
    })
};