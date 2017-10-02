import React, { Component } from 'react';
import 'whatwg-fetch';
import './App.css';

import * as GoApiService from './services/goApiService.js';

class App extends Component {

    constructor(props) {
        super(props);
        this.state = {
            summary: { data: 'empty' },
        }
    }

    componentWillMount() {
    }

    handleOnSubmit(e) {
        e.preventDefault();
        console.log('event target', e.target.url.value);
        GoApiService.getSummaryResourcePromise(e.target.url.value).then((json) => {
            this.setState({summary: json})
        })
    }

    render() {
        return (
            <div className="App">
                <h1>This is a generic web client for INFO344</h1>
                <div>
                    <form onSubmit={this.handleOnSubmit.bind(this)}>
                        URL to get summary for: <input type="text" placeholder="e.g.: https://google.com" name="url"/>
                        <input type="submit" value="Submit"/>
                    </form>
                </div>
                <div className="json-container">
                    <pre className="align-left">{JSON.stringify(this.state.summary, null, 2)}</pre>
                </div>
            </div>
        );
    }
}

export default App;
