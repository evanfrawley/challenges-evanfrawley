import 'whatwg-fetch';

const LOCAL_GO_API_PATH = "http://localhost:4000";

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
        });
}
