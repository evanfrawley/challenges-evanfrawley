import React, { Component } from 'react';
import 'whatwg-fetch';
import './App.css';

import * as GoApiService from './services/goApiService.js';
import * as YoutubeService from './services/YoutubeURLCleaner.js';

class App extends Component {

    constructor(props) {
        super(props);
        this.state = {
            summary: { data: 'empty' },
        }
    }

    handleOnSubmit(e) {
        e.preventDefault();
        console.log('event target', e.target.url.value);
        GoApiService.getSummaryResourcePromise(e.target.url.value).then((json) => {
            this.setState({summary: json})
        })
    }

    render() {
        let videoDiv = null;
        if (this.state.summary && this.state.summary.videos) {
            let videos = this.state.summary.videos.map((video) => {
                console.log(video);
                let cleanVideoUrl = YoutubeService.getEmbedYoutubeUrl(video.url);
                console.log('clean', cleanVideoUrl);
                if (video.type) {
                    if(video.type === "text/html") {
                        return (
                            <div key={cleanVideoUrl}>
                                <iframe src={cleanVideoUrl}></iframe>
                            </div>
                        );
                    } else if (video.type.startsWith("video/")) {
                        return (
                            <div key={cleanVideoUrl}>
                                <video src={cleanVideoUrl}></video>
                            </div>
                        );
                    }
                } else {
                    return (
                        <div key={cleanVideoUrl}>
                            <iframe src={cleanVideoUrl}></iframe>
                        </div>
                    );
                }
            });
            videoDiv = (
                <div>
                    {videos}
                </div>
            )
        }
        return (
            <div className="App">
                <h1>This is a generic web client for INFO344</h1>
                <div>
                    <p>to test video, try using this URL: https://www.keithandthegirl.com/vip/bonus/episode/9/40/this-is-40</p>
                </div>
                <div>
                    <form onSubmit={this.handleOnSubmit.bind(this)}>
                        URL to get summary for: <input type="text" placeholder="e.g.: https://google.com" name="url"/>
                        <input type="submit" value="Submit"/>
                    </form>
                </div>
                <div className="json-container">
                    <pre className="align-left">{JSON.stringify(this.state.summary, null, 2)}</pre>
                </div>
                { videoDiv && videoDiv }
            </div>
        );
    }
}

export default App;
