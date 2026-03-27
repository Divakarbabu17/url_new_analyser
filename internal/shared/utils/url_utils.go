package utils

import (
    "net/url"
    "strings"
)

// NormalizeURL converts relative URLs to absolute URLs
func NormalizeURL(href, base string) string {
    baseParsed, err := url.Parse(base)
    if err != nil {
        return href
    }
    hrefParsed, err := url.Parse(href)
    if err != nil {
        return href
    }
    return baseParsed.ResolveReference(hrefParsed).String()
}

// IsExternalLink returns true if href is external compared to base domain
func IsExternalLink(href, base string) bool {
    hrefURL, err1 := url.Parse(href)
    baseURL, err2 := url.Parse(base)
    if err1 != nil || err2 != nil {
        return true
    }
    return hrefURL.Host != baseURL.Host
}

// IsInternalLink returns true if href belongs to the same domain
func IsInternalLink(href, base string) bool {
    return !IsExternalLink(href, base)
}

// CleanURL trims spaces and fragments from URL
func CleanURL(href string) string {
    u, err := url.Parse(strings.TrimSpace(href))
    if err != nil {
        return href
    }
    u.Fragment = ""
    return u.String()
}