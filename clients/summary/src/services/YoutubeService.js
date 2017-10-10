import React from 'react';

const YOUTUBE_EMBED_URL_REGEX = "^https?://(www.)?youtube.com/embed/.+";
const YOUTUBE_URL_WATCH_QP_REGEX = "^https?://(www.)?youtube.com/watch?.+v=.+";
const YOUTUBE_URL_V_RESOURCE = "^https?://(www.)?youtube.com/v/.+";
const YOUTUBE_URL = "^https?://(www.)?youtube.com.+";

export function getEmbedYoutubeUrl(rawUrl) {
    console.log('raw', rawUrl);
    if (rawUrl.match(YOUTUBE_URL)) {
        if (rawUrl.match(YOUTUBE_EMBED_URL_REGEX)) {
            // yay! we're already there
            return rawUrl;
        } else {
            // snag from watch
            if (rawUrl.match(YOUTUBE_URL_WATCH_QP_REGEX)) {
                let queryParams = rawUrl.split("?")[1].split("&");
                let idParam = queryParams.filter((param) => {
                    return param.startsWith("v=");
                });
                console.log('idparam', idParam);
                let id = idParam.replace("v=", "");
                console.log('id', id);
                return `http://youtube.com/embed/${id}`;
            } else if (rawUrl.match(YOUTUBE_URL_V_RESOURCE)) {
                // snag from youtube.com/v/{id}
                return rawUrl.replace("/v/", "/embed/");
            } else {
                return rawUrl;
            }
        }
    }
}

export function createEmbedVideoArray(videos) {
    return videos.map((video) => {
        console.log(video);
        let cleanVideoUrl = getEmbedYoutubeUrl(video.url);
        console.log('clean', cleanVideoUrl);
        if (video.type) {
            if (video.type === "text/html") {
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
}