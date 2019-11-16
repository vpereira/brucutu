package main

import "net/url"

func parseURL(u string) (parsedURL *url.URL, err error) {
	myURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	return myURL, nil
}
