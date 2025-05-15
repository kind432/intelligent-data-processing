package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func ParseSwitchTopic(topic string) (serialNumber string, err error) {
	parts := strings.Split(topic, "/")
	if len(parts) < 4 {
		return "", fmt.Errorf("invalid topic format: too few segments")
	}

	if strings.Contains(parts[0], "@") {
		serialParts := strings.Split(parts[2], "_")
		serialNumber = strings.Join(serialParts[:3], "-")
	} else {
		serialParts := strings.Split(parts[0], "-")
		if len(serialParts) == 3 {
			serialNumber = strings.Join(serialParts, "_")
		} else {
			return "", fmt.Errorf("invalid serial number format: %s", parts[0])
		}
	}
	return
}

func ParseSensorTopic(topic string) (serialNumber, dataKey string, err error) {
	parts := strings.Split(topic, "/")
	if len(parts) < 4 {
		return "", "", fmt.Errorf("topic format is invalid")
	}

	raw := parts[2]
	subParts := strings.Split(raw, "_")

	baseName := "robbo-protos"
	switch len(subParts) {
	case 0:
		return "", "", fmt.Errorf("invalid device name in topic")
	case 1:
		serialNumber = baseName
		dataKey = subParts[0]
	case 2:
		serialNumber = baseName
		dataKey = strings.Join(subParts[1:], "_")
	case 3:
		if isNumber(subParts[2]) {
			serialNumber = fmt.Sprintf("%s-%s", baseName, subParts[2])
			dataKey = subParts[1]
		} else {
			serialNumber = baseName
			dataKey = strings.Join(subParts[1:], "_")
		}
	default:
		if isNumber(subParts[2]) {
			serialNumber = fmt.Sprintf("%s-%s", baseName, subParts[2])
			dataKey = strings.Join(subParts[3:], "_")
		} else {
			serialNumber = baseName
			dataKey = strings.Join(subParts[2:], "_")
		}
	}

	return serialNumber, dataKey, nil
}

func isNumber(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func MatchesTopicPattern(pattern, topic string) bool {
	regexPattern := strings.ReplaceAll(pattern, "+", "[^/]+")
	re := regexp.MustCompile("^" + regexPattern + "$")
	return re.MatchString(topic)
}
