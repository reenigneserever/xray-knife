package xray

import (
	"github.com/reenigneserever/xray-knife/utils"
	"github.com/reenigneserever/xray-knife/utils/customlog"
	"io"
	"log"
	"net/url"
	"strings"
)

// TODO: Make a database to store subscriptions
type Subscription struct {
	Remark      string
	Url         string
	UserAgent   string
	Method      string
	ConfigLinks []string
}

func (s *Subscription) FetchAll() ([]string, error) {
	ll, _ := url.Parse(s.Url)
	if s.Method == "" {
		s.Method = "GET"
	}
	response, err := utils.SendHTTPRequest(ll, s.UserAgent, s.Method)
	if err != nil {
		return nil, err
	}
	bytes, _ := io.ReadAll(response.Body)
	decoded, err2 := utils.Base64Decode(string(bytes))
	if err2 != nil {
		// Probably It's not base64 encoded!, let's try parsing without decoding
		customlog.Printf(customlog.Processing, "Couldn't decode the body! let's try parsing without decoding...\n")
		links := strings.Split(string(bytes), "\n")
		_, err = ParseXrayConfig(links[0])
		if err != nil {
			customlog.Printf(customlog.Failure, "Failed to parse the first config from the subscription!\n")
			return nil, err2
		}
		s.ConfigLinks = links
		return links, nil
	}
	// Configs are separated by newline char
	links := strings.Split(string(decoded), "\n")
	s.ConfigLinks = links
	return links, nil
}

func (s *Subscription) RemoveDuplicate(verbose bool) {
	// Remove duplicates using hashmap (hashed keys)
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range s.ConfigLinks {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	if verbose {
		log.Printf("Removed %d duplicate configs!\n", len(s.ConfigLinks)-len(list))
	}
	s.ConfigLinks = list
}
