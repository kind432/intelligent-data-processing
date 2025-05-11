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

func (h Handler) HandleFloatSensorData(client mqtt.Client, msg mqtt.Message) {
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

func (h Handler) HandlePrutok1State(client mqtt.Client, msg mqtt.Message) {
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
		"sensorType": "default",
		"prutok_1":   rawValue,
	}

	handler := h.TopicHandlers[msg.Topic()]
	outputRawTopic := fmt.Sprintf("%s/%s", serialNumber, handler.OutputRawTopic)
	h.publishMessage(client, outputRawTopic, rawData)
}

func (h Handler) HandlePrutok2State(client mqtt.Client, msg mqtt.Message) {
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
		"sensorType": "default",
		"prutok_2":   rawValue,
	}

	handler := h.TopicHandlers[msg.Topic()]
	outputRawTopic := fmt.Sprintf("%s/%s", serialNumber, handler.OutputRawTopic)
	h.publishMessage(client, outputRawTopic, rawData)
}

func (h Handler) HandleDoorState(client mqtt.Client, msg mqtt.Message) {
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
		"sensorType": "default",
		"door":       rawValue,
	}

	handler := h.TopicHandlers[msg.Topic()]
	outputRawTopic := fmt.Sprintf("%s/%s", serialNumber, handler.OutputRawTopic)
	h.publishMessage(client, outputRawTopic, rawData)
}

func (h Handler) HandleSmokeState(client mqtt.Client, msg mqtt.Message) {
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
		"sensorType": "default",
		"smoke":      rawValue,
	}

	handler := h.TopicHandlers[msg.Topic()]
	outputRawTopic := fmt.Sprintf("%s/%s", serialNumber, handler.OutputRawTopic)
	h.publishMessage(client, outputRawTopic, rawData)
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

func (h Handler) HandleYellowLedState(client mqtt.Client, msg mqtt.Message) {
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
		"sensorType": "default",
		"yellow_led": rawValue,
	}

	handler := h.TopicHandlers[msg.Topic()]
	outputRawTopic := fmt.Sprintf("%s/%s", serialNumber, handler.OutputRawTopic)
	h.publishMessage(client, outputRawTopic, rawData)
}

func (h Handler) HandleYellowLedCommand(client mqtt.Client, msg mqtt.Message) {
	var command map[string]bool
	if err := json.Unmarshal(msg.Payload(), &command); err != nil {
		h.Logger.Err.Printf("Failed to parse command: %v", err)
		return
	}

	yellowLedState, exists := command["yellow_led"]
	if !exists {
		h.Logger.Err.Printf("Missing 'yellow_led' in command payload")
		return
	}

	yellowLedStateStr := "OFF"
	if yellowLedState {
		yellowLedStateStr = "ON"
	}

	serialNumber, err := utils.ParseSwitchTopic(msg.Topic())
	if err != nil {
		h.Logger.Err.Printf("Invalid topic format: '%s': %v", msg.Topic(), err)
		return
	}

	relayTopic := fmt.Sprintf("%s/switch/%s_yellow_led/command", viper.GetString("mqtt_username"), serialNumber)
	h.publishMessage(client, relayTopic, yellowLedStateStr)
}

func (h Handler) HandleRedLedState(client mqtt.Client, msg mqtt.Message) {
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
		"sensorType": "default",
		"red_led":    rawValue,
	}

	handler := h.TopicHandlers[msg.Topic()]
	outputRawTopic := fmt.Sprintf("%s/%s", serialNumber, handler.OutputRawTopic)
	h.publishMessage(client, outputRawTopic, rawData)
}

func (h Handler) HandleRedLedCommand(client mqtt.Client, msg mqtt.Message) {
	var command map[string]bool
	if err := json.Unmarshal(msg.Payload(), &command); err != nil {
		h.Logger.Err.Printf("Failed to parse command: %v", err)
		return
	}

	redLedState, exists := command["red_led"]
	if !exists {
		h.Logger.Err.Printf("Missing 'red_led' in command payload")
		return
	}

	redLedStateStr := "OFF"
	if redLedState {
		redLedStateStr = "ON"
	}

	serialNumber, err := utils.ParseSwitchTopic(msg.Topic())
	if err != nil {
		h.Logger.Err.Printf("Invalid topic format: '%s': %v", msg.Topic(), err)
		return
	}

	relayTopic := fmt.Sprintf("%s/switch/%s_red_led/command", viper.GetString("mqtt_username"), serialNumber)
	h.publishMessage(client, relayTopic, redLedStateStr)
}

func (h Handler) HandleMPUChooseState(client mqtt.Client, msg mqtt.Message) {
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
		"sensorType": "default",
		"mpu_choose": rawValue,
	}

	handler := h.TopicHandlers[msg.Topic()]
	outputRawTopic := fmt.Sprintf("%s/%s", serialNumber, handler.OutputRawTopic)
	h.publishMessage(client, outputRawTopic, rawData)
}

func (h Handler) HandleMPUChooseCommand(client mqtt.Client, msg mqtt.Message) {
	var command map[string]bool
	if err := json.Unmarshal(msg.Payload(), &command); err != nil {
		h.Logger.Err.Printf("Failed to parse command: %v", err)
		return
	}

	mpuChooseState, exists := command["mpu_choose"]
	if !exists {
		h.Logger.Err.Printf("Missing 'mpu_choose' in command payload")
		return
	}

	mpuChooseStateStr := "OFF"
	if mpuChooseState {
		mpuChooseStateStr = "ON"
	}

	serialNumber, err := utils.ParseSwitchTopic(msg.Topic())
	if err != nil {
		h.Logger.Err.Printf("Invalid topic format: '%s': %v", msg.Topic(), err)
		return
	}

	relayTopic := fmt.Sprintf("%s/switch/%s_mpu_choose/command", viper.GetString("mqtt_username"), serialNumber)
	h.publishMessage(client, relayTopic, mpuChooseStateStr)
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
	var payload []byte

	switch v := data.(type) {
	case string:
		payload = []byte(v)
	default:
		var err error
		payload, err = json.Marshal(v)
		if err != nil {
			h.Logger.Err.Printf("Error marshalling data: %v", err)
			return
		}
	}

	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		h.Logger.Err.Printf("Error publishing message to '%s': %v", topic, token.Error())
		return
	}
	h.Logger.Info.Printf("Published message to '%s': %s", topic, payload)
}
