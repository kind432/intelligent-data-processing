package handler

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"intelligent-data-processing/internal/device"
	"intelligent-data-processing/internal/sensor"
	"intelligent-data-processing/pkg/logger"
	"intelligent-data-processing/pkg/utils"
	"strconv"
)

type Handler struct {
	Logger        logger.Logger
	TopicHandlers map[string]TopicHandler
}

func NewHandler(logger logger.Logger) Handler {
	return Handler{
		Logger:        logger,
		TopicHandlers: make(map[string]TopicHandler),
	}
}

func (h Handler) MessageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		h.Logger.Info.Printf("Received message from topic '%s': %s", msg.Topic(), msg.Payload())

		matched := false
		for pattern, topicHandler := range h.TopicHandlers {
			if utils.MatchesTopicPattern(pattern, msg.Topic()) {
				matched = true
				topicHandler.HandlerFunc(client, msg)
				break
			}
		}

		if !matched {
			h.Logger.Err.Printf("No handler found for topic '%s'", msg.Topic())
		}
	}
}

func (h Handler) HandlePowerRelayState(client mqtt.Client, msg mqtt.Message) {
	var rawValue bool

	serialNumber, err := utils.ParseSwitchTopic(msg.Topic())
	if err != nil {
		h.Logger.Err.Printf("Invalid topic format: '%s': %v", msg.Topic(), err)
		return
	}

	rawValueStr := string(msg.Payload())
	switch rawValueStr {
	case "ON":
		rawValue = true
	case "OFF":
		rawValue = false
	default:
		h.Logger.Err.Printf("Unexpected payload value: '%s'", rawValueStr)
		return
	}

	rawData := map[string]interface{}{
		"sensorType":  "default",
		"power_relay": rawValue,
	}

	handler := h.TopicHandlers[msg.Topic()]
	outputRawTopic := fmt.Sprintf("%s/%s", serialNumber, handler.OutputRawTopic)
	h.publishMessage(client, outputRawTopic, rawData)
}

func (h Handler) HandleSensorData(client mqtt.Client, msg mqtt.Message) {
	serialNumber, dataKey, err := utils.ParseSensorTopic(msg.Topic())
	if err != nil {
		h.Logger.Err.Printf("Error parsing topic '%s': %v", msg.Topic(), err)
		return
	}

	if !device.IsDeviceConnected(serialNumber) {
		h.SendConnectMessage(client, serialNumber)
	}

	rawValue, err := strconv.ParseFloat(string(msg.Payload()), 64)
	if err != nil {
		h.Logger.Err.Printf("Invalid value format: '%s': %v", msg.Payload(), err)
		return
	}

	rawData, procData, err := sensor.ProcessSensorData(serialNumber, dataKey, rawValue)
	if err != nil {
		h.Logger.Err.Printf("Error processing data: %v", err)
		return
	}

	handler := h.TopicHandlers[msg.Topic()]
	outputRawTopic := fmt.Sprintf("%s/%s", serialNumber, handler.OutputRawTopic)
	outputProcTopic := fmt.Sprintf("%s/%s", serialNumber, handler.OutputProcTopic)
	h.publishMessage(client, outputRawTopic, rawData)
	h.publishMessage(client, outputProcTopic, procData)
}

func (h Handler) HandlePowerRelayCommand(client mqtt.Client, msg mqtt.Message) {
	var command map[string]bool
	if err := json.Unmarshal(msg.Payload(), &command); err != nil {
		h.Logger.Err.Printf("Failed to parse command: %v", err)
		return
	}

	relayState, exists := command["power_relay"]
	if !exists {
		h.Logger.Err.Printf("Missing 'power_relay' in command payload")
		return
	}

	relayStateStr := "OFF"
	if relayState {
		relayStateStr = "ON"
	}

	serialNumber, err := utils.ParseSwitchTopic(msg.Topic())
	if err != nil {
		h.Logger.Err.Printf("Invalid topic format: '%s': %v", msg.Topic(), err)
		return
	}

	relayTopic := fmt.Sprintf("%s/switch/%s_power_relay/command", viper.GetString("mqtt_username"), serialNumber)
	h.publishMessage(client, relayTopic, relayStateStr)
}

func (h Handler) SendConnectMessage(client mqtt.Client, serialNumber string) {
	data := map[string]string{"sensorType": "default"}
	topic := fmt.Sprintf("%s/sensor/connect", serialNumber)
	device.AddOrUpdateDevice(serialNumber, true)
	h.publishMessage(client, topic, data)
}

func (h Handler) SendDisconnectMessage(client mqtt.Client, serialNumber string) {
	data := map[string]string{"sensorType": "default"}
	topic := fmt.Sprintf("%s/sensor/disconnect", serialNumber)
	device.AddOrUpdateDevice(serialNumber, false)
	h.publishMessage(client, topic, data)
}

func (h Handler) publishMessage(client mqtt.Client, topic string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		h.Logger.Err.Printf("Error marshalling data: %v", err)
		return
	}
	token := client.Publish(topic, 0, false, jsonData)
	token.Wait()
	if token.Error() != nil {
		h.Logger.Err.Printf("Error publishing message to '%s': %v", topic, token.Error())
		return
	}
	h.Logger.Info.Printf("Published message to '%s': %s", topic, jsonData)
}
