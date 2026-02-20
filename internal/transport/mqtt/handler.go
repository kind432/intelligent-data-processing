package mqtt

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"

	"intelligent-data-processing/internal/repository"
	"intelligent-data-processing/internal/service"
	"intelligent-data-processing/pkg/logger"
	"intelligent-data-processing/pkg/utils"
)

type Handler struct {
	log  logger.Logger
	repo *repository.Storage
	proc *service.Processor
}

func NewHandler(log logger.Logger, repo *repository.Storage, proc *service.Processor) *Handler {
	return &Handler{log: log, repo: repo, proc: proc}
}

func (h *Handler) UniversalHandler(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payloadStr := string(msg.Payload())

	info, err := utils.ParseTopic(topic)
	if err != nil {
		h.log.Err.Printf("Skipping unhandled topic %s: %v", topic, err)
		return
	}

	if info.Action == "command" && !info.IsExternal {
		return
	}

	if info.Action == "command" && info.IsExternal {
		h.handleCommand(client, info, msg.Payload())
		return
	}

	if !info.IsExternal && !h.repo.IsConnected(info.SerialNumber) {
		h.publishDeviceStatus(client, info.SerialNumber, true)
	}

	h.handleState(client, info, payloadStr)
}

func (h *Handler) handleState(client mqtt.Client, info *utils.TopicInfo, payload string) {
	if val, err := strconv.ParseFloat(payload, 64); err == nil {
		avg := h.proc.ProcessFloat(info.SerialNumber, info.Key, val)

		rawTopic := fmt.Sprintf("%s/%s/%s/raw", info.SerialNumber, info.Category, info.Key)
		procTopic := fmt.Sprintf("%s/%s/%s/proc", info.SerialNumber, info.Category, info.Key)

		h.publishData(client, rawTopic, map[string]interface{}{
			"sensorType": "default",
			info.Key:     val,
		})

		h.publishData(client, procTopic, map[string]interface{}{
			"sensorType":       "default",
			"proc_" + info.Key: avg,
		})
		return
	}

	var boolVal bool
	if payload == "ON" {
		boolVal = true
	} else if payload == "OFF" {
		boolVal = false
	} else {
		h.log.Err.Printf("Unknown payload for state: %s", payload)
		return
	}

	rawTopic := fmt.Sprintf("%s/%s/%s/raw", info.SerialNumber, info.Category, info.Key)
	h.publishData(client, rawTopic, map[string]interface{}{
		"sensorType": "default",
		info.Key:     boolVal,
	})
}

func (h *Handler) handleCommand(client mqtt.Client, info *utils.TopicInfo, payload []byte) {
	var command map[string]bool
	if err := json.Unmarshal(payload, &command); err != nil {
		h.log.Err.Printf("Failed to parse command payload: %v", err)
		return
	}

	state, exists := command[info.Key]
	if !exists {
		h.log.Err.Printf("Missing key '%s' in command payload", info.Key)
		return
	}

	stateStr := "OFF"
	if state {
		stateStr = "ON"
	}

	id := ""
	reID := regexp.MustCompile(`\d+`)
	id = reID.FindString(info.SerialNumber)

	var brokerTopic string
	if id != "" {
		brokerTopic = fmt.Sprintf("%s/%s/robbo_protos_%s_%s/command",
			viper.GetString("mqtt_username"), info.Category, id, info.Key)
	} else {
		brokerTopic = fmt.Sprintf("%s/%s/robbo_protos_%s/command",
			viper.GetString("mqtt_username"), info.Category, info.Key)
	}

	h.publishData(client, brokerTopic, stateStr)
}

func (h *Handler) publishDeviceStatus(client mqtt.Client, serialNumber string, connected bool) {
	h.repo.UpdateDeviceStatus(serialNumber, connected)
	status := "disconnect"
	if connected {
		status = "connect"
	}
	topic := fmt.Sprintf("%s/sensor/%s", serialNumber, status)
	h.publishData(client, topic, map[string]string{"sensorType": "default"})
}

func (h *Handler) publishData(client mqtt.Client, topic string, data interface{}) {
	var payload []byte
	switch v := data.(type) {
	case string:
		payload = []byte(v)
	default:
		payload, _ = json.Marshal(v)
	}

	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		h.log.Err.Printf("Error publishing to '%s': %v", topic, token.Error())
	} else {
		h.log.Info.Printf("Published to '%s'", topic)
	}
}
