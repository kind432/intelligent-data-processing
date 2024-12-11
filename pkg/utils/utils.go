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
	serialParts := strings.Split(parts[2], "_")
	serialNumber = strings.Join(serialParts[:3], "-")
	dataKey = strings.Join(serialParts[3:], "_")
	return
}

func MatchesTopicPattern(pattern, topic string) bool {
	regexPattern := strings.ReplaceAll(pattern, "+", "[^/]+")
	re := regexp.MustCompile("^" + regexPattern + "$")
	return re.MatchString(topic)
}
