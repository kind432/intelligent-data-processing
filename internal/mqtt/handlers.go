package mqtt

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"intelligent-data-processing/pkg/logger"
	"regexp"
	"strings"
)

type TopicHandler struct {
	OutputRawTopic  string
	OutputProcTopic string
}

var topicHandlers = make(map[string]TopicHandler)

func initTopicHandlers() {
	sensors := []struct {
		Key        string
		RawSuffix  string
		ProcSuffix string
	}{
		{"temperature_1", "sensor/temperature_1/raw", "sensor/temperature_1/proc"},
		{"temperature_2", "sensor/temperature_2/raw", "sensor/temperature_2/proc"},
		{"motor_current", "sensor/motor_current/raw", "sensor/motor_current/proc"},
		{"gas_sensor", "sensor/gas_sensor/raw", "sensor/gas_sensor/proc"},
		{"mpu6050_temperature", "sensor/mpu6050_temperature/raw", "sensor/mpu6050_temperature/proc"},
		{"mpu6050_gyro_z", "sensor/mpu6050_gyro_z/raw", "sensor/mpu6050_gyro_z/proc"},
		{"mpu6050_accel_z", "sensor/mpu6050_accel_z/raw", "sensor/mpu6050_accel_z/proc"},
		{"mpu6050_gyro_y", "sensor/mpu6050_gyro_y/raw", "sensor/mpu6050_gyro_y/proc"},
		{"mpu6050_accel_y", "sensor/mpu6050_accel_y/raw", "sensor/mpu6050_accel_y/proc"},
		{"mpu6050_gyro_x", "sensor/mpu6050_gyro_x/raw", "sensor/mpu6050_gyro_x/proc"},
		{"mpu6050_accel_x", "sensor/mpu6050_accel_x/raw", "sensor/mpu6050_accel_x/proc"},
		{"ina226_power", "sensor/ina226_power/raw", "sensor/ina226_power/proc"},
		{"ina226_current", "sensor/ina226_current/raw", "sensor/ina226_current/proc"},
		{"ina226_shunt_voltage", "sensor/ina226_shunt_voltage/raw", "sensor/ina226_shunt_voltage/proc"},
	}

	baseTopic := "yarila682@yandex.ru/sensor/ROBBO_protos_%02d_%s/state"
	for i := 1; i <= 8; i++ {
		for _, sensor := range sensors {
			topic := fmt.Sprintf(baseTopic, i, sensor.Key)
			topicHandlers[topic] = TopicHandler{
				OutputRawTopic:  sensor.RawSuffix,
				OutputProcTopic: sensor.ProcSuffix,
			}
		}
	}
}

const ConnectTopic = "sensor/connect"
const DisconnectTopic = "sensor/disconnect"

type Connect struct {
	SensorType string `json:"sensorType"`
}

func messageHandler(loggers logger.Loggers) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		loggers.Info.Printf("Received message from topic '%s': %s", msg.Topic(), msg.Payload())

		parts := strings.Split(msg.Topic(), "/")
		if len(parts) != 4 || parts[1] != "sensor" || parts[3] != "state" {
			loggers.Err.Printf("Invalid topic format: '%s'", msg.Topic())
			return
		}

		re := regexp.MustCompile(`([^/]+)@[^/]+/sensor/(.+)/state`)
		matches := re.FindStringSubmatch(msg.Topic())
		if len(matches) != 3 {
			loggers.Err.Printf("Failed to parse topic: '%s'", msg.Topic())
			return
		}

		serialNumber := matches[1]
		dataKey := matches[2]

		loggers.Info.Printf("Parsed serialNumber: %s, dataKey: %s", serialNumber, dataKey)

		topicKey := fmt.Sprintf("+/sensor/%s/state", dataKey)
		handler, exists := topicHandlers[topicKey]
		if !exists {
			loggers.Err.Printf("No handler found for topic '%s'", msg.Topic())
			return
		}

		connectData := Connect{
			SensorType: "default",
		}
		connectDeviceTopic := serialNumber + "/" + ConnectTopic
		if err := publishMessage(client, connectDeviceTopic, connectData, loggers); err != nil {
			loggers.Err.Printf("Failed to publish connect message: %v", err)
			return
		}

		rawData, processedData, err := processSensorData(serialNumber, dataKey, msg.Payload())
		if err != nil {
			loggers.Err.Printf("Error processing message: %v", err)
			return
		}

		outputRawTopic := serialNumber + "/" + handler.OutputRawTopic
		if err := publishMessage(client, outputRawTopic, rawData, loggers); err != nil {
			loggers.Err.Printf("Failed to publish raw data: %v", err)
			return
		}

		outputProcTopic := serialNumber + "/" + handler.OutputProcTopic
		if err := publishMessage(client, outputProcTopic, processedData, loggers); err != nil {
			loggers.Err.Printf("Failed to publish processed data: %v", err)
			return
		}

		disconnectTopic := serialNumber + "/" + DisconnectTopic
		if err := publishMessage(client, disconnectTopic, nil, loggers); err != nil {
			loggers.Err.Printf("Failed to publish disconnect message: %v", err)
		}
	}
}

func publishMessage(client mqtt.Client, topic string, data interface{}, loggers logger.Loggers) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data: %v", err)
	}
	token := client.Publish(topic, 0, false, jsonData)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("error publishing message to topic '%s': %v", topic, token.Error())
	}
	loggers.Info.Printf("Published message to topic '%s': %s", topic, jsonData)
	return nil
}
