package handlers

import (
    "encoding/json"
    "io"
    "io/ioutil"
    "reflect"
    "strings"
    "testing"
)

func TestExtractSummary(t *testing.T) {
    pagePrologue := "<html><head>"
    pageEiplogue := "</head><body></body></html>"
    pageURL := "http://test.com/test.html"
    cases := []struct {
        name            string
        hint            string
        html            string
        expectedSummary *PageSummary
    }{
        {
            "Open Graph Type",
            `Make sure you are reading the <meta property="og:type" content="..."> element`,
            pagePrologue + `<meta property="og:type" content="test type">` + pageEiplogue,
            &PageSummary{
                Type: "test type",
            },
        },
        {
            "Open Graph URL",
            `Make sure you are reading the <meta property="og:url" content="..."> element`,
            pagePrologue + `<meta property="og:url" content="http://test.com">` + pageEiplogue,
            &PageSummary{
                URL: "http://test.com",
            },
        },
        {
            "Open Graph Title",
            `Make sure you are reading the <meta property="og:title" content="..."> element`,
            pagePrologue + `<meta property="og:title" content="test title">` + pageEiplogue,
            &PageSummary{
                Title: "test title",
            },
        },
        {
            "Open Graph Site name",
            `Make sure you are reading the <meta property="og:site_name" content="..."> element`,
            pagePrologue + `<meta property="og:site_name" content="test site name">` + pageEiplogue,
            &PageSummary{
                SiteName: "test site name",
            },
        },
        {
            "Open Graph Description",
            `Make sure you are reading the <meta property="og:description" content="..."> element`,
            pagePrologue + `<meta property="og:description" content="test description">` + pageEiplogue,
            &PageSummary{
                Description: "test description",
            },
        },
        {
            "Open Graph Image",
            `Make sure you are reading the <meta property="og:image" content="..."> element`,
            pagePrologue + `<meta property="og:image" content="http://test.com/test.png">` + pageEiplogue,
            &PageSummary{
                Images: []*PreviewImage{
                    {
                        URL: "http://test.com/test.png",
                    },
                },
            },
        },
        {
            "Open Graph Structured Image",
            `Make sure you are handling the image structured properties, as described in http://ogp.me/#structured`,
            pagePrologue + `
            <meta property="og:image" content="http://test.com/test.png">
            <meta property="og:image:secure_url" content="https://test.com/test.png">
            <meta property="og:image:type" content="image/png">
            <meta property="og:image:width" content="300">
            <meta property="og:image:height" content="300">
            <meta property="og:image:alt" content="test alt">
            ` + pageEiplogue,
            &PageSummary{
                Images: []*PreviewImage{
                    {
                        URL:       "http://test.com/test.png",
                        SecureURL: "https://test.com/test.png",
                        Type:      "image/png",
                        Width:     300,
                        Height:    300,
                        Alt:       "test alt",
                    },
                },
            },
        },
        {
            "Open Graph Structured Video",
            `Make sure you are handling the image structured properties, as described in http://ogp.me/#structured`,
            pagePrologue + `
            <meta property="og:video" content="http://test.com/test.mp4">
            <meta property="og:video:secure_url" content="https://test.com/test.mp4">
            <meta property="og:video:type" content="video/mp4">
            <meta property="og:video:width" content="300">
            <meta property="og:video:height" content="300">
            ` + pageEiplogue,
            &PageSummary{
                Videos: []*PreviewVideo{
                    {
                        URL:       "http://test.com/test.mp4",
                        SecureURL: "https://test.com/test.mp4",
                        Type:      "video/mp4",
                        Width:     300,
                        Height:    300,
                    },
                },
            },
        },
        {
            "Open Graph Structured Audio",
            `Make sure you are handling the image structured properties, as described in http://ogp.me/#structured`,
            pagePrologue + `
            <meta property="og:audio" content="http://test.com/test.mp3">
            <meta property="og:audio:secure_url" content="https://test.com/test.mp3">
            <meta property="og:audio:type" content="audio/mp3">
            ` + pageEiplogue,
            &PageSummary{
                Audio: []*PreviewAudio{
                    {
                        URL:       "http://test.com/test.mp3",
                        SecureURL: "https://test.com/test.mp3",
                        Type:      "audio/mp3",
                    },
                },
            },
        },
        {
            "Open Graph Multiple Images",
            `Make sure you are handling multiple images, as described in http://ogp.me/#array`,
            pagePrologue + `
            <meta property="og:image" content="http://test.com/test1.png">
            <meta property="og:image:width" content="100">
            <meta property="og:image:height" content="100">
            <meta property="og:image:alt" content="test alt 1">
            <meta property="og:image" content="http://test.com/test2.png">
            <meta property="og:image:alt" content="test alt 2">
            <meta property="og:image" content="http://test.com/test3.png">
            <meta property="og:image:alt" content="test alt 3">
            ` + pageEiplogue,
            &PageSummary{
                Images: []*PreviewImage{
                    {
                        URL:    "http://test.com/test1.png",
                        Width:  100,
                        Height: 100,
                        Alt:    "test alt 1",
                    },
                    {
                        URL: "http://test.com/test2.png",
                        Alt:    "test alt 2",
                    },
                    {
                        URL: "http://test.com/test3.png",
                        Alt: "test alt 3",
                    },
                },
            },
        },
        {
            "Open Graph Multiple Videos",
            `Make sure you are handling multiple images, as described in http://ogp.me/#array`,
            pagePrologue + `
            <meta property="og:video" content="http://test.com/test1.mp4">
            <meta property="og:video:width" content="100">
            <meta property="og:video:height" content="100">
            <meta property="og:video" content="http://test.com/test2.mp4">
            <meta property="og:video" content="http://test.com/test3.mp4">
            <meta property="og:video:height" content="100">
            ` + pageEiplogue,
            &PageSummary{
                Videos: []*PreviewVideo{
                    {
                        URL:    "http://test.com/test1.mp4",
                        Width:  100,
                        Height: 100,
                    },
                    {
                        URL: "http://test.com/test2.mp4",
                    },
                    {
                        URL: "http://test.com/test3.mp4",
                        Height: 100,
                    },
                },
            },
        },
        {
            "Open Graph Multiple Audio files",
            `Make sure you are handling multiple images, as described in http://ogp.me/#array`,
            pagePrologue + `
            <meta property="og:audio" content="http://test.com/test1.mp3">
            <meta property="og:audio:secure_url" content="https://test.com/test1.mp3">
            <meta property="og:audio:type" content="audio/mp3">
            <meta property="og:audio" content="http://test.com/test2.mp3">
            <meta property="og:audio" content="http://test.com/test3.mp3">
            <meta property="og:audio:type" content="audio/mp3">
            ` + pageEiplogue,
            &PageSummary{
                Audio: []*PreviewAudio{
                    {
                        URL:    "http://test.com/test1.mp3",
                        SecureURL:  "https://test.com/test1.mp3",
                        Type: "audio/mp3",
                    },
                    {
                        URL: "http://test.com/test2.mp3",
                    },
                    {
                        URL: "http://test.com/test3.mp3",
                        Type: "audio/mp3",
                    },
                },
            },
        },
        {
            "All Open Graph Props",
            "Make sure you are handling all of the Open Graph properties listed in the assignment",
            pagePrologue + `
            <meta property="og:type" content="test type">
            <meta property="og:url" content="http://test.com">
            <meta property="og:title" content="test title">
            <meta property="og:site_name" content="test site name">
            <meta property="og:description" content="test description">
            <meta property="og:image" content="http://test.com/test.png">
            ` + pageEiplogue,
            &PageSummary{
                Type:        "test type",
                URL:         "http://test.com",
                Title:       "test title",
                SiteName:    "test site name",
                Description: "test description",
                Images: []*PreviewImage{
                    {
                        URL: "http://test.com/test.png",
                    },
                },
            },
        },
        {
            "HTML Title",
            `Make sure you get the page title from the <title> element if not Open Graph title property is available`,
            pagePrologue + `<title>HTML Page Title</title>` + pageEiplogue,
            &PageSummary{
                Title: "HTML Page Title",
            },
        },
        {
            "HTML Description",
            `Make sure you get the page description from the <meta name="author" content="..."> tag if no Open Graph description is available`,
            pagePrologue + `<meta name="description" content="test description">` + pageEiplogue,
            &PageSummary{
                Description: "test description",
            },
        },
        {
            "HTML Author",
            `Make sure you get the page author from the <meta name="author" content="..."> tag`,
            pagePrologue + `<meta name="author" content="test author">` + pageEiplogue,
            &PageSummary{
                Author: "test author",
            },
        },
        {
            "HTML Keywords With Spaces",
            `Make sure you get the page keywords from the <meta name="keywords" content="..."> tag`,
            pagePrologue + `<meta name="keywords" content="one, two, three">` + pageEiplogue,
            &PageSummary{
                Keywords: []string{"one", "two", "three"},
            },
        },
        {
            "HTML Keywords With No Spaces",
            `Make sure you get the page keywords from the <meta name="keywords" content="..."> tag`,
            pagePrologue + `<meta name="keywords" content="one,two,three">` + pageEiplogue,
            &PageSummary{
                Keywords: []string{"one", "two", "three"},
            },
        },
        {
            "HTML Icon",
            `Make sure you get the page icon from the <link rel="icon" href="..."> tag`,
            pagePrologue + `<link rel="icon" href="http://test.com/test.png">` + pageEiplogue,
            &PageSummary{
                Icon: &PreviewImage{
                    URL: "http://test.com/test.png",
                },
            },
        },
        {
            "HTML Icon With Sizes",
            `Make sure you parse the 'sizes' attribute to get the icon height and width`,
            pagePrologue + `<link rel="icon" href="http://test.com/test.png" sizes="100x200">` + pageEiplogue,
            &PageSummary{
                Icon: &PreviewImage{
                    URL:    "http://test.com/test.png",
                    Height: 100,
                    Width:  200,
                },
            },
        },
        {
            "HTML Icon With Size Any",
            `The sizes attribute of the <link rel="icon"> tag may have the value "any" to indicate no size preference`,
            pagePrologue + `<link rel="icon" href="http://test.com/test.png" sizes="any">` + pageEiplogue,
            &PageSummary{
                Icon: &PreviewImage{
                    URL: "http://test.com/test.png",
                },
            },
        },
        {
            "HTML Icon With Type",
            `Make sure you read the 'type' attribute to get the icon type`,
            pagePrologue + `<link rel="icon" href="http://test.com/test.png" type="image/png">` + pageEiplogue,
            &PageSummary{
                Icon: &PreviewImage{
                    URL:  "http://test.com/test.png",
                    Type: "image/png",
                },
            },
        },
        {
            "Self-Closing Meta",
            "Make sure you are handling self-closing <meta ... /> tags",
            pagePrologue + `<meta property="og:title" content="Open Graph Title"/>` + pageEiplogue,
            &PageSummary{
                Title: "Open Graph Title",
            },
        },
        {
            "Attribute Order",
            "Attributes in HTML can be in any order; don't assume a particular order",
            pagePrologue + `<meta content="Open Graph Title" property="og:title"/>` + pageEiplogue,
            &PageSummary{
                Title: "Open Graph Title",
            },
        },
        {
            "HTML and Open Graph Title",
            `Make sure the <meta property="og:title"> overrides the HTML <title> element`,
            pagePrologue + `
            <meta property="og:title" content="Open Graph Title"/>
            <title>HTML Page Title</title>` + pageEiplogue,
            &PageSummary{
                Title: "Open Graph Title",
            },
        },
        {
            "HTML and Open Graph Description",
            `Make sure the <meta property="og:description"> overrides the HTML <meta name="description"> element`,
            pagePrologue + `
            <meta property="og:description" content="og description"/>
            <meta name="description" content="html description">` + pageEiplogue,
            &PageSummary{
                Description: "og description",
            },
        },
        {
            "Relative Image URL",
            "Remember to resolve relative image URLs to absolute ones using the page URL as a base",
            pagePrologue + `<meta property="og:image" content="/test.png"/>` + pageEiplogue,
            &PageSummary{
                Images: []*PreviewImage{
                    {
                        URL: "http://test.com/test.png",
                    },
                },
            },
        },
        {
            "Relative Icon URL",
            "Remember to resolve relative HTML Icon URLs to absolute ones using the page URL as a base",
            pagePrologue + `<link rel="icon" href="/test.png"/>` + pageEiplogue,
            &PageSummary{
                Icon: &PreviewImage{
                    URL: "http://test.com/test.png",
                },
            },
        },
        {
            "Twitter Scheme",
            "Check out the twitter documentation",
            pagePrologue + `
            <meta property="twitter:card" content="some_twitter_content"/>
            <meta property="twitter:image" content="http://test.com/test.png"/>
            <meta property="twitter:description" content="this should be the description"/>
            <meta property="twitter:title" content="twitter scheme is cool"/>
            ` + pageEiplogue,
            &PageSummary{
                Images: []*PreviewImage{
                    {
                        URL: "http://test.com/test.png",
                    },
                },
                Title: "twitter scheme is cool",
                Description: "this should be the description",
                Type: "some_twitter_content",
            },
        },
        {
            "Twitter Scheme with Open Graph content",
            "Check out the twitter documentation",
            pagePrologue + `
            <meta property="og:type" content="text/html"/>
            <meta property="og:description" content="open graph helps devs"/>
            <meta property="og:title" content="open graph is cool"/>
            <meta property="twitter:card" content="some_twitter_content"/>
            <meta property="twitter:image" content="http://test.com/test.png"/>
            <meta property="twitter:description" content="this should be the description"/>
            <meta property="twitter:title" content="twitter scheme is cool"/>
            ` + pageEiplogue,
            &PageSummary{
                Images: []*PreviewImage{
                    {
                        URL: "http://test.com/test.png",
                    },
                },
                Title: "open graph is cool",
                Description: "open graph helps devs",
                Type: "text/html",
            },
        },
        {
            "Empty Input",
            "A URL might return an empty page",
            "",
            &PageSummary{},
        },
    }

    for _, c := range cases {
        summary, err := extractSummary(pageURL, ioutil.NopCloser(strings.NewReader(c.html)))
        if err != nil && err != io.EOF {
            t.Errorf("case %s: unexpected error %v\nHINT: %s\n", c.name, err, c.hint)
        }
        if !reflect.DeepEqual(summary, c.expectedSummary) {
            expectedJSON, _ := json.MarshalIndent(c.expectedSummary, "", "  ")
            actualJSON, _ := json.MarshalIndent(summary, "", "  ")
            t.Errorf("case %s: incorrect result:\nEXPECTED: %s\nACTUAL: %s\nHINT: %s\n",
                c.name, string(expectedJSON), string(actualJSON), c.hint)
        }
    }
}

func TestFetchHTML(t *testing.T) {
    cases := []struct {
        name        string
        hint        string
        URL         string
        expectError bool
    }{
        {
            "Valid URL",
            "This is a valid HTML page, so this should work",
            "https://info344-a17.github.io/tests/ogall.html",
            false,
        },
        {
            "Not Found URL",
            "Remember to check the response status code",
            "https://info344-a17.github.io/tests/not-found.html",
            true,
        },
        {
            "Non-HTML URL",
            "Remember to check the response content-type to ensure it's an HTML page",
            "https://info344-a17.github.io/tests/test.png",
            true,
        },
    }

    for _, c := range cases {
        stream, err := fetchHTML(c.URL)

        if err != nil && !c.expectError {
            t.Errorf("case %s: unexpected error %v\nHINT: %s", c.name, err, c.hint)
        }
        if c.expectError && err == nil {
            t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
        }

        if stream != nil {
            stream.Close()
        }
    }

}
