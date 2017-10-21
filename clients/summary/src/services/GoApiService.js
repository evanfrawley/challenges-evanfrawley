import 'whatwg-fetch';

const LOCAL_GO_API_PATH = "https://api.evan.gg";

export function getSummaryResourcePromise(urlToQuery) {
    let path = "/v1/summary";
    return fetch(`${LOCAL_GO_API_PATH}${path}?url=${urlToQuery}`)
        .then(function(response) {
            if (response.ok) {
                return response.json();
            } else {
                return response.text()
                    .then(function(errorMessage) {
                        throw new Error(errorMessage)
                    });
            }
        })
        .catch(function(err) {
            console.error(err);
            alert(`You were given the following error:\n${err}`)
        });
}
