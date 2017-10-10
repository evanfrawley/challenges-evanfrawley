import React, {Component} from 'react';
import 'whatwg-fetch';
import './App.css';

import * as GoApiService from './services/GoApiService.js';
import { createEmbedVideoArray } from './services/YoutubeService.js';

class App extends Component {

    constructor(props) {
        super(props);
        this.state = {
            summary: null,
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
        let imageDiv = null;
        let videoDiv = null;
        if (this.state.summary && this.state.summary.images) {
            let images = this.state.summary.images.map((image) => {
                let height = image.height ? image.height : 'auto';
                let width = image.width ? image.width : 'auto';
                let type = image.type ? image.type : '';
                let alt = image.alt ? image.alt : '';
                return <img src={image.url} alt={alt} height={height} width={width} type={type}/>
            });
            imageDiv = (
                <div>
                    {images}
                </div>
            )
        }
        if (this.state.summary && this.state.summary.videos) {
            let videos = createEmbedVideoArray(this.state.summary.videos);
            videoDiv = (
                <div>
                    {videos}
                </div>
            )
        }
        return (
            <div className="App">
                <h1>Page Summary</h1>
                <div>
                    <form onSubmit={this.handleOnSubmit.bind(this)}>
                        URL to get summary for: <input type="text" placeholder="e.g.: https://google.com" name="url"/>
                        <input type="submit" value="Submit"/>
                    </form>
                </div>
                <div className="json-container">
                    {this.state.summary &&
                        <div className="align-left">
                            {this.state.summary.title &&
                            <p>Title: {this.state.summary.title}</p>
                            }
                            {this.state.summary.url &&
                                <p>URL: <a href={this.state.summary.url}>{this.state.summary.url}</a></p>
                            }
                            {this.state.summary.description &&
                                <p>Description: {this.state.summary.description}</p>
                            }
                            {imageDiv && imageDiv}
                            {videoDiv && videoDiv}
                        </div>
                    }
                </div>
            </div>
        );
    }
}


export default App;
