package handlers

import (
    "io"
    "net/http"
    "log"
    "encoding/json"

    "golang.org/x/net/html"
    "strings"
    "strconv"
    "fmt"
    "net/url"
)

type OpenGraphPreview struct {
}

//PreviewImage represents a preview image for a page
type PreviewImage struct {
    URL       string `json:"url,omitempty"`
    SecureURL string `json:"secureURL,omitempty"`
    Type      string `json:"type,omitempty"`
    Width     int    `json:"width,omitempty"`
    Height    int    `json:"height,omitempty"`
    Alt       string `json:"alt,omitempty"`
}

//PreviewVideo represents a preview video for a page
type PreviewVideo struct {
    URL       string `json:"url,omitempty"`
    SecureURL string `json:"secureURL,omitempty"`
    Type      string `json:"type,omitempty"`
    Width     int    `json:"width,omitempty"`
    Height    int    `json:"height,omitempty"`
}

//PreviewAudio represents a preview audio for a page
type PreviewAudio struct {
    URL       string `json:"url,omitempty"`
    SecureURL string `json:"secureURL,omitempty"`
    Type      string `json:"type,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
    Type        string          `json:"type,omitempty"`
    URL         string          `json:"url,omitempty"`
    Title       string          `json:"title,omitempty"`
    SiteName    string          `json:"siteName,omitempty"`
    Description string          `json:"description,omitempty"`
    Author      string          `json:"author,omitempty"`
    Keywords    []string        `json:"keywords,omitempty"`
    Icon        *PreviewImage   `json:"icon,omitempty"`
    Images      []*PreviewImage `json:"images,omitempty"`
    Videos      []*PreviewVideo `json:"videos,omitempty"`
    Audio       []*PreviewAudio `json:"audio,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Add(AccessControlAllowOriginKey, AccessControlAllowOriginVal)
    w.Header().Add(ContentTypeKey, ContentTypeJSONUTF8Val)

    URL := r.URL.Query().Get("url")

    reader, err := fetchHTML(URL)

    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    defer reader.Close()

    pageSummary, err := extractSummary(URL, reader)
    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    json.NewEncoder(w).Encode(pageSummary)
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
    // convert to http
    resp, err := http.Get(pageURL)

    //if there was an error, report it and exit
    if err != nil {
        return nil, fmt.Errorf("error fetching URL: %v\n", err)
    }

    if !(resp.StatusCode >= 200 && resp.StatusCode < 400) {
        return nil, fmt.Errorf("response status code was %d \n", resp.StatusCode)
    }

    contentType := resp.Header.Get(ContentTypeKey)
    if !strings.HasPrefix(contentType, ContentTypeHTMLVal) {
        return nil, fmt.Errorf("response content type was %s not %s\n", contentType, ContentTypeHTMLVal)
    }

    return resp.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
    tokenizer := html.NewTokenizer(htmlStream)

    pageSummary := &PageSummary{}
    tempPreviewImage := &PreviewImage{}
    tempPreviewVideo := &PreviewVideo{}
    tempPreviewAudio := &PreviewAudio{}

    seenTwitterMutuallyExclusiveOGTags := map[string]bool{}
    for {
        tokenType := tokenizer.Next()

        if tokenType == html.ErrorToken {
            return pageSummary, nil
        }

        t := tokenizer.Token()

        if tokenType == html.EndTagToken && "head" == t.Data {
            return pageSummary, nil
        }

        attributeMap := getAttributeMap(t.Attr)

        switch t.Data {
        case "link":
            {
                err := handleLinkTagData(attributeMap, pageURL, pageSummary)
                if err != nil {
                    return nil, err
                }
            }
        case "meta":
            {
                metaIDType := attributeMap["metaIdType"]
                metaIDTypeVal := attributeMap["metaIdVal"]
                content := attributeMap["content"]
                if metaIDTypeVal == "og:image" || metaIDTypeVal == "og:image:url" {
                    // Handle initialize images slice
                    newTempPreviewImage, err := initializeAndAppendImagesSlice(pageSummary, tempPreviewImage, pageURL, content)
                    if err != nil {
                        return nil, err
                    }
                    tempPreviewImage = newTempPreviewImage
                    seenTwitterMutuallyExclusiveOGTags["og:image"] = true
                } else if metaIDTypeVal == "og:video" {
                    // Handle initializing videos slice
                    absoluteURL, err := getAbsoluteURL(pageURL, content)
                    if err != nil {
                        return nil, fmt.Errorf("error while parsing URL: %v", err)
                    }
                    if pageSummary.Videos == nil {
                        pageSummary.Videos = []*PreviewVideo{}
                    }
                    if tempPreviewVideo.URL == "" {
                        tempPreviewVideo.URL = absoluteURL
                    } else {
                        // If new image entry, create new PreviewImage and update the PageSummary slice
                        tempPreviewVideo = &PreviewVideo{URL: absoluteURL}
                    }
                    pageSummary.Videos = append(pageSummary.Videos, tempPreviewVideo)
                } else if metaIDTypeVal == "og:audio" {
                    // Handle initializing audio slice
                    absoluteURL, err := getAbsoluteURL(pageURL, content)
                    if err != nil {
                        return nil, fmt.Errorf("error while parsing URL: %v", err)
                    }
                    if pageSummary.Audio == nil {
                        pageSummary.Audio = []*PreviewAudio{}
                    }
                    if tempPreviewAudio.URL == "" {
                        tempPreviewAudio.URL = absoluteURL
                    } else {
                        // If new image entry, create new PreviewImage and update the PageSummary slice
                        tempPreviewAudio = &PreviewAudio{URL: absoluteURL}
                    }
                    pageSummary.Audio = append(pageSummary.Audio, tempPreviewAudio)
                    seenTwitterMutuallyExclusiveOGTags["og:audio"] = true
                } else if strings.HasPrefix(metaIDTypeVal, "og:image") {
                    err := handleImagePrefixMetaData(tempPreviewImage, metaIDTypeVal, content)
                    if err != nil {
                        return nil, err
                    }
                } else if strings.HasPrefix(metaIDTypeVal, "og:video") {
                    err := handleVideoPrefixMetaData(tempPreviewVideo, metaIDTypeVal, content)
                    if err != nil {
                        return nil, err
                    }
                } else if strings.HasPrefix(metaIDTypeVal, "og:audio") {
                    handleAudioPrefixMetaData(tempPreviewAudio, metaIDTypeVal, content)
                } else if strings.HasPrefix(metaIDTypeVal, "twitter") {
                    switch metaIDTypeVal {
                    case "twitter:image":
                        {
                            hasTwitterImg := false
                            for _, image := range pageSummary.Images {
                                if image.URL == content {
                                    hasTwitterImg = true
                                }
                            }
                            if !hasTwitterImg {
                                newTempPreviewImage, err := initializeAndAppendImagesSlice(pageSummary, tempPreviewImage, pageURL, content)
                                if err != nil {
                                    return nil, err
                                }
                                tempPreviewImage = newTempPreviewImage
                            }
                        }
                    case "twitter:card": {
                        if !seenTwitterMutuallyExclusiveOGTags["og:type"] {
                            pageSummary.Type = content
                        }
                    }
                    case "twitter:title":
                        {
                            if !seenTwitterMutuallyExclusiveOGTags["og:title"] {
                                pageSummary.Title = content
                            }
                        }
                    case "twitter:description":
                        {
                            if !seenTwitterMutuallyExclusiveOGTags["og:description"] {
                                pageSummary.Description = content
                            }
                        }
                    }
                } else {
                    handleStandardMetaTagData(pageSummary, metaIDType, metaIDTypeVal, content, seenTwitterMutuallyExclusiveOGTags)
                }
            }
        case "title":
            handleTitle(&t, tokenizer, pageSummary)
        default: // Do nothing!
        }
    }

    return pageSummary, nil
}

func initializeAndAppendImagesSlice(pageSummary *PageSummary, image *PreviewImage, pageURL, content string) (*PreviewImage, error) {
    absoluteURL, err := getAbsoluteURL(pageURL, content)
    if err != nil {
        return nil, fmt.Errorf("error while parsing URL: %v", err)
    }
    if pageSummary.Images == nil {
        pageSummary.Images = []*PreviewImage{}
    }
    image = &PreviewImage{URL: absoluteURL}
    pageSummary.Images = append(pageSummary.Images, image)
    return image, nil
}

// Handles populating the title into the PageSummary
func handleTitle(t *html.Token, tokenizer *html.Tokenizer, pageSummary *PageSummary) {
    if t.Type == html.StartTagToken {
        tokenizer.Next()
        titleToken := tokenizer.Token()
        if titleToken.Type == html.TextToken {
            if pageSummary.Title == "" {
                pageSummary.Title = titleToken.Data
            }
        }
    }
}

// Handles populating the supplementary `og:image:*` values
func handleImagePrefixMetaData(tempPreviewImage *PreviewImage, metaIDTypeVal, content string) (error) {
    // Add to image object
    err := handlePreviewImageMetaData(tempPreviewImage, metaIDTypeVal, content)
    if err != nil {
        if metaIDTypeVal == "og:image:height" {
            tempPreviewImage.Height = 0
        } else if metaIDTypeVal == "og:image:height" {
            tempPreviewImage.Width = 0
        } else {
            // Uncaught error
            return err
        }
    }

    return nil
}

// Handles populating the supplementary `og:video:*` values
func handleVideoPrefixMetaData(video *PreviewVideo, metaIDTypeVal, content string) (error) {
    err := handlePreviewVideoMetaData(video, metaIDTypeVal, content)
    if err != nil {
        if metaIDTypeVal == "og:video:height" {
            video.Height = 0
        } else if metaIDTypeVal == "og:video:height" {
            video.Width = 0
        } else {
            // Uncaught error
            return err
        }
    }

    return nil
}

// Handles populating the supplementary `og:audio:*` values
func handleAudioPrefixMetaData(audio *PreviewAudio, metaIDTypeVal, content string) {
    switch metaIDTypeVal {
    case "og:audio:secure_url":
        audio.SecureURL = content
    case "og:audio:type":
        audio.Type = content
    default: // Do nothing!
    }
}

func handleLinkTagData(attributeMap map[string]string, pageURL string, pageSummary *PageSummary) (error) {
    rel := attributeMap["rel"]
    if rel == "icon" {
        href := attributeMap["href"]
        // href is required
        if href == "" {
            return fmt.Errorf("the href attribute is required on a link")
        }
        href, err := getAbsoluteURL(pageURL, href)
        if err != nil {
            return fmt.Errorf("error while parsing URL: %v", err)
        }
        linkType := attributeMap["type"]
        width, height, err := parseLinkImageSizes(attributeMap["sizes"])
        if err != nil {
            log.Printf("error while parsing image size: %v \n defaulting height and width to be 0 \n", err)
        }

        iconPreviewImage := &PreviewImage{
            URL:    href,
            Width:  width,
            Height: height,
            Type:   linkType,
        }

        pageSummary.Icon = iconPreviewImage
    }

    return nil
}

func handleStandardMetaTagData(pageSummary *PageSummary, tagType, tagValue, content string, seenTags map[string]bool) {
    switch tagType {
    case "name":
        switch tagValue {
        case "author":
            pageSummary.Author = content
        case "description":
            {
                if pageSummary.Description == "" {
                    pageSummary.Description = content
                }
            }
        case "keywords":
            {
                keywords := strings.Split(content, ",")
                for index, keyword := range keywords {
                    keywords[index] = strings.TrimSpace(keyword)
                }
                pageSummary.Keywords = keywords
            }
        default: // Do nothing!
        }
    case "property":
        switch tagValue {
        case "og:type":
            pageSummary.Type = content
            seenTags["og:type"] = true
        case "og:url":
            pageSummary.URL = content
        case "og:title":
            pageSummary.Title = content
            seenTags["og:title"] = true
        case "og:site_name":
            pageSummary.SiteName = content
        case "og:description":
            pageSummary.Description = content
            seenTags["og:description"] = true
        case "og:image": // Do nothing! case handled above
        default: // Do nothing!
        }
    }
}

func handlePreviewImageMetaData(image *PreviewImage, imageAttribute, content string) (error) {
    switch imageAttribute {
    case "og:image:secure_url":
        {
            image.SecureURL = content
        }
    case "og:image:type":
        image.Type = content
    case "og:image:width":
        {
            widthInt, err := strconv.Atoi(content)
            if err != nil {
                return fmt.Errorf("an error occurred parsing the width: %v", err)
            } else {
                image.Width = widthInt
            }
        }
    case "og:image:height":
        {
            heightInt, err := strconv.Atoi(content)
            if err != nil {
                return fmt.Errorf("an error occurred parsing the height: %v", err)
            } else {
                image.Height = heightInt
            }
        }
    case "og:image:alt":
        image.Alt = content
    default: // Do nothing!
    }

    return nil
}

// Much duplicated code as above, but I couldn't figure out how to reuse it without inheritance
func handlePreviewVideoMetaData(video *PreviewVideo, attribute, content string) (error) {
    switch attribute {
    case "og:video:secure_url":
        video.SecureURL = content
    case "og:video:type":
        video.Type = content
    case "og:video:width":
        {
            widthInt, err := strconv.Atoi(content)
            if err != nil {
                return fmt.Errorf("an error occurred parsing the width: %v", err)
            } else {
                video.Width = widthInt
            }
        }
    case "og:video:height":
        {
            heightInt, err := strconv.Atoi(content)
            if err != nil {
                return fmt.Errorf("an error occurred parsing the height: %v", err)
            } else {
                video.Height = heightInt
            }
        }
    default: // Do nothing!
    }

    return nil
}

func getAttributeMap(attributes []html.Attribute) (map[string]string) {
    attributeMap := map[string]string{}

    for _, element := range attributes {
        if element.Key == "name" || element.Key == "property" {
            attributeMap["metaIdType"] = element.Key
            attributeMap["metaIdVal"] = element.Val
        } else {
            attributeMap[element.Key] = element.Val
        }
    }

    return attributeMap
}

func parseLinkImageSizes(sizes string) (int, int, error) {
    if sizes != "" && sizes != "any" {
        parsedSizes := strings.Split(sizes, "x")
        height, err := strconv.Atoi(parsedSizes[0])

        if err != nil {
            return 0, 0, err
        }

        width, err := strconv.Atoi(parsedSizes[1])
        if err != nil {
            return 0, 0, err
        }
        return width, height, nil
    } else {
        return 0, 0, nil
    }
}

// Returns absolute URL of a particular resource or reference
func getAbsoluteURL(parentURL, relativeUrl string) (string, error) {
    parsedURL, err := url.Parse(parentURL)
    if err != nil {
        return "", fmt.Errorf("illegal url: %s", parentURL)
    }
    childURL, err := url.Parse(relativeUrl)
    return parsedURL.ResolveReference(childURL).String(), nil
}
