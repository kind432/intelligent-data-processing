package utils

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	brokerRegex   = regexp.MustCompile(`^([^/]+)/([^/]+)/robbo_protos(?:_(\d+))?_(.+)/(state|command)$`)
	externalRegex = regexp.MustCompile(`^robbo-protos(?:-(\d+))?/([^/]+)/([^/]+)/(state|command)$`)
)

type TopicInfo struct {
	IsExternal   bool
	Username     string
	Category     string
	SerialNumber string
	Key          string
	Action       string
}

func ParseTopic(topic string) (*TopicInfo, error) {
	if strings.HasPrefix(topic, "robbo-protos") {
		matches := externalRegex.FindStringSubmatch(topic)
		if len(matches) == 5 {
			info := &TopicInfo{
				IsExternal: true,
				Category:   matches[2],
				Key:        matches[3],
				Action:     matches[4],
			}
			if matches[1] != "" {
				info.SerialNumber = "robbo-protos-" + matches[1]
			} else {
				info.SerialNumber = "robbo-protos"
			}
			return info, nil
		}
	}

	matches := brokerRegex.FindStringSubmatch(topic)
	if len(matches) == 6 {
		info := &TopicInfo{
			IsExternal: false,
			Username:   matches[1],
			Category:   matches[2],
			Key:        matches[4],
			Action:     matches[5],
		}
		if matches[3] != "" {
			info.SerialNumber = "robbo-protos-" + matches[3]
		} else {
			info.SerialNumber = "robbo-protos"
		}
		return info, nil
	}

	matches = externalRegex.FindStringSubmatch(topic)
	if len(matches) == 5 {
		info := &TopicInfo{
			IsExternal: true,
			Category:   matches[2],
			Key:        matches[3],
			Action:     matches[4],
		}
		if matches[1] != "" {
			info.SerialNumber = "robbo-protos-" + matches[1]
		} else {
			info.SerialNumber = "robbo-protos"
		}
		return info, nil
	}

	return nil, fmt.Errorf("unknown topic structure: %s", topic)
}
