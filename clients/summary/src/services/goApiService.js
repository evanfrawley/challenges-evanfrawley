import 'whatwg-fetch';

const LOCAL_GO_API_PATH = "http://localhost:1460";

export function getSummaryResourcePromise(urlToQuery) {
    let path = "/v1/summary";
    return fetch(`${LOCAL_GO_API_PATH}${path}?url=${urlToQuery}`).then((response) => {
        return response.json();
    });
}