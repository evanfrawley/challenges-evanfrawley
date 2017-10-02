package handlers

import (
    "io"
    "net/http"
    "log"
    "encoding/json"

    "golang.org/x/net/html"
    "strings"
    "strconv"
    "regexp"
    "fmt"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
    URL       string `json:"url,omitempty"`
    SecureURL string `json:"secureURL,omitempty"`
    Type      string `json:"type,omitempty"`
    Width     int    `json:"width,omitempty"`
    Height    int    `json:"height,omitempty"`
    Alt       string `json:"alt,omitempty"`
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
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Access-Control-Allow-Origin", "*")
    w.Header().Add("Content-Type", "application/json")

    URL := r.URL.Query().Get("url")

    reader, err := fetchHTML(URL)
    if err != nil {
        log.Fatalf("some error occurred when fetching HTML: %s", err)
    }

    pageSummary, err := extractSummary(URL, reader)
    if err != nil {
        log.Fatalf("some error occurred when extracting HTML: %s", err)
    }

    reader.Close()

    json.NewEncoder(w).Encode(pageSummary)
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
    resp, err := http.Get(pageURL)

    //if there was an error, report it and exit
    if err != nil {
        return nil, fmt.Errorf("error fetching URL: %v\n", err)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("response status code was %d \n", resp.StatusCode)
    }

    contentType := resp.Header.Get("Content-Type")
    if !strings.HasPrefix(contentType, "text/html") {
        return nil, fmt.Errorf("response content type was %s not text/html\n", contentType)
    }

    return resp.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
    tokenizer := html.NewTokenizer(htmlStream)

    pageSummary := &PageSummary{}
    tempPreviewImage := &PreviewImage{}
    for {
        nextTokenType := tokenizer.Next()

        if nextTokenType == html.ErrorToken {
            log.Printf("error tokenizing HTML: %v \n", tokenizer.Err())
            break
        }

        nextToken := tokenizer.Token()

        if nextTokenType == html.EndTagToken && "head" == nextToken.Data {
            break
        }

        switch nextToken.Data {
        case "link": {
            attributes := nextToken.Attr
            rel := findAndGetValueForAttribute(attributes, "rel")
            if rel == "icon" {
                href := findAndGetValueForAttribute(attributes, "href")
                // href is required
                if href == "" {
                    log.Fatal("The href attribute is required on a link")
                }
                href = getAbsoluteURL(pageURL, href)
                linkType := findAndGetValueForAttribute(attributes, "type")
                width, height := parseLinkImageSizes(findAndGetValueForAttribute(attributes, "sizes"))

                iconPreviewImage := &PreviewImage{
                    URL: href,
                    Width: width,
                    Height: height,
                    Type: linkType,
                }

                pageSummary.Icon = iconPreviewImage
            }
        }
        case "meta": {
            metaIDType := getMetaIDType(nextToken.Attr)
            metaIDTypeVal := getMetaIDTypeVal(nextToken.Attr)
            content := getMetaContent(nextToken.Attr)
            if metaIDTypeVal == "og:image" {
                absoluteURL := getAbsoluteURL(pageURL, content)
                if pageSummary.Images == nil {
                    pageSummary.Images = []*PreviewImage{}
                }
                if tempPreviewImage.URL != "" {
                    // If new image entry, create new PreviewImage and update the PageSummary slice
                    tempPreviewImage = &PreviewImage{ URL: absoluteURL }
                } else {
                    tempPreviewImage.URL = absoluteURL
                }
                pageSummary.Images = append(pageSummary.Images, tempPreviewImage)
            } else if strings.HasPrefix(metaIDTypeVal, "og:image") {
                // Add to image object
                handlePreviewImageData(tempPreviewImage, metaIDTypeVal, content)
            } else {
                handleMetaTagData(pageSummary, metaIDType, metaIDTypeVal, content)
            }
        }
        case "title": {
            if nextToken.Type == html.StartTagToken {
                tokenizer.Next()
                titleToken := tokenizer.Token()
                if titleToken.Type == html.TextToken {
                    if pageSummary.Title == "" {
                        pageSummary.Title = titleToken.Data
                    }
                }
            }
        }
        default: // Do nothing!
        }
    }

    return pageSummary, nil
}


func handleMetaTagData(pageSummary *PageSummary, tagType, tagValue, content string) {
    switch tagType {
    case "name":
        switch tagValue {
        case "author": pageSummary.Author = content
        case "description": {
            if pageSummary.Description == "" {
                pageSummary.Description = content
            }
        }
        case "keywords": {
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
        case "og:type": pageSummary.Type = content
        case "og:url": pageSummary.URL = content
        case "og:title": pageSummary.Title = content
        case "og:site_name": pageSummary.SiteName = content
        case "og:description": pageSummary.Description = content
        case "og:image": // Do nothing! case handled above
        default: // Do nothing!
        }
    }
}

func handlePreviewImageData(image *PreviewImage, imageAttribute, content string) {
    switch imageAttribute {
    case "og:image:secure_url": {
        image.SecureURL = content
    }
    case "og:image:type": image.Type = content
    case "og:image:width": {
        widthInt, err := strconv.Atoi(content)
        if err != nil {
            log.Printf("an error occurred parsing the width: %s", err)
            image.Width = 0
        } else {
            image.Width = widthInt
        }
    }
    case "og:image:height": {
        heightInt, err := strconv.Atoi(content)
        if err != nil {
            log.Printf("an error occurred parsing the height: %s", err)
            image.Height = 0
        } else {
            image.Height = heightInt
        }
    }
    case "og:image:alt": image.Alt = content
    default: // Do nothing!
    }
}

func getMetaIDTypeVal(attributes []html.Attribute) (string) {
    name := findAndGetValueForAttribute(attributes, "name")
    property := findAndGetValueForAttribute(attributes, "property")
    if property != "" {
        return property
    } else {
        return name
    }
}

func getMetaIDType(attributes []html.Attribute) (string) {
    metaIDType := ""
    for _, element := range attributes {
        if element.Key == "name" || element.Key == "property" {
            metaIDType = element.Key
        }
    }

    return metaIDType
}

func getMetaContent(attributes []html.Attribute) (string) {
    return findAndGetValueForAttribute(attributes, "content")
}

func findAndGetValueForAttribute(attributes []html.Attribute, targetAttribute string) (string) {
    metaIDType := ""
    for _, element := range attributes {
        if element.Key == targetAttribute {
            metaIDType = element.Val
        }
    }

    return metaIDType
}

func parseLinkImageSizes(sizes string) (int, int) {
    if sizes != "" && sizes != "any" {
        parsedSizes := strings.Split(sizes, "x")
        height, err := strconv.Atoi(parsedSizes[0])

        if err != nil {
            log.Printf("Failure parsing height: %s \n", err)
            height = 0
        }

        width, err := strconv.Atoi(parsedSizes[1])
        if err != nil {
            log.Printf("Failure parsing width: %s \n", err)
            width = 0
        }
        return width, height
    } else {
        return 0, 0
    }
}

func getAbsoluteURL(parentURL, relativeUrl string) (string) {
    httpsPrefix := "^https?://*"
    matched, err := regexp.MatchString(httpsPrefix, relativeUrl)
    if err != nil {
        log.Printf("regex was syntactically incorrect: %s", httpsPrefix)
    }

    returnURL := ""
    if matched {
        // is absolute URL
        returnURL = relativeUrl
    } else {
        // is relative URL
        relativeDirBackCount := 0

        // remove any and all `../`, `./`, `/` from the relative img path and count how many times they were removed
        for {
            if strings.HasPrefix(relativeUrl, "../") {
                relativeUrl = strings.Replace(relativeUrl, "../", "", 1)
                relativeDirBackCount++
            } else if strings.HasPrefix(relativeUrl, "./")  {
                relativeUrl = strings.Replace(relativeUrl, "./", "", 1)
                relativeDirBackCount++
            } else if strings.HasPrefix(relativeUrl, "/") {
                relativeUrl = strings.Replace(relativeUrl, "/", "", 1)
                relativeDirBackCount++
            } else {
                break
            }
        }

        // sanitize parent URL to have no trailing `/`
        if string(parentURL[len(parentURL) - 1]) == "/" {
            parentURL = string(parentURL[:len(parentURL) - 1])
        }

        // separate protocol and URL body
        protocolAndRestOfLink := strings.Split(parentURL, "://")

        // split host and resources
        urlPieces := strings.Split(protocolAndRestOfLink[1], "/")

        // select only necessary resource paths, traversing up directories if the page provided was relative with `../` `./` `/`
        relativePathUrlPieces := urlPieces[:len(urlPieces) - relativeDirBackCount]

        // append sanitized relative path
        relativePathUrlPieces = append(relativePathUrlPieces, relativeUrl)

        // construct url body
        linkBody := strings.Join(relativePathUrlPieces[:],"/")

        // prepend protocol to joined link body
        finalAbsoluteURLSlice := []string {protocolAndRestOfLink[0], linkBody}

        // join protocol and link body
        finalAbsoluteURL := strings.Join(finalAbsoluteURLSlice, "://")

        returnURL = finalAbsoluteURL
    }

    return returnURL
}