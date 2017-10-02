import 'whatwg-fetch';

const LOCAL_GO_API_PATH = "http://localhost:1460";

export function fetchNameResource(name) {
    let str = `http://localhost:4000/hello?name=${name}`;
    fetch(str).then((response) => {
        return response.text();
    }).then((text) => {
        this.setState({phrase: text})
    })
}

export function fetchMemResource() {
    let str = "http://localhost:4000/meme";
    fetch(str).then((response) => {
        return response.json();
    }).then((json) => {
        this.setState({mem: json["Alloc"]})
    })
}

export function getSummaryResourcePromise(urlToQuery) {
    let path = "/v1/summary";
    return fetch(`${LOCAL_GO_API_PATH}${path}?url=${urlToQuery}`).then((response) => {
        return response.json();
    });
}