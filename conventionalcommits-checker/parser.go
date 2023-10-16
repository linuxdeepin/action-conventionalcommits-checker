package main

import (
	"regexp"
	"strings"
)

type SubjectData struct {
	Valid       bool
	Type        string
	Scope       string
	Description string
}

type Trailer struct {
	Key   string
	Value string
}

type TrailersData struct {
	Trailers []Trailer
}

func indexAtStringArray(arr []string, idx int) string {
	if idx < 0 || idx >= len(arr) {
		return ""
	}

	return arr[idx]
}

func ParseSubjectLine(subject string) SubjectData {
	// [skip test]type(scope): description subject line
	r := regexp.MustCompile(`(?s)(?:\[[\w ]+\])?(?P<type>[\w]+)(?:\((?P<scope>\w[\w\/]*)\))?!?: (?P<description>.*)`)
	matches := r.FindStringSubmatch(subject)

	return SubjectData{
		Valid:       len(matches) > 0,
		Type:        indexAtStringArray(matches, r.SubexpIndex("type")),
		Scope:       indexAtStringArray(matches, r.SubexpIndex("scope")),
		Description: indexAtStringArray(matches, r.SubexpIndex("description")),
	}
}

func ParseTrailers(trailers string) TrailersData {
	var result TrailersData

	trailersLines := strings.Split(trailers, "\n")
	r := regexp.MustCompile(`(?P<key>[\w]+):(?: (?P<value>.*))?`)

	for _, trailer := range trailersLines {
		matches := r.FindStringSubmatch(trailer)
		if len(matches) > 0 {
			result.Trailers = append(result.Trailers, Trailer{
				Key:   indexAtStringArray(matches, r.SubexpIndex("key")),
				Value: indexAtStringArray(matches, r.SubexpIndex("value")),
			})
		}
	}

	return result
}

func (data *TrailersData) Contains(key string) bool {
	for _, trailer := range data.Trailers {
		if trailer.Key == key {
			return true
		}
	}

	return false
}

func (data *TrailersData) StrictContains(key string) bool {
	for _, trailer := range data.Trailers {
		if trailer.Key == key && trailer.Value != "" {
			return true
		}
	}

	return false
}
