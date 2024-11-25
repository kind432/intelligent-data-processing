package mqtt

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"intelligent-data-processing/pkg/logger"
	"strings"
)

type TopicHandler struct {
	SensorType string
	// returning: rawData, processedData, error
	ProcessFunc     func([]byte) (map[string]interface{}, map[string]interface{}, error)
	OutputRawTopic  string
	OutputProcTopic string
}

// Example:
// Topic: SN-001/sensor/temperature_2/state
// Then:
// TopicHandler:
// "+/sensor/temperature_2/state" : {
// 		SensorType:      "default",
//		ProcessFunc:     func([]byte) (map[string]interface{}, map[string]interface{}, error),
//		OutputRawTopic:  "sensor/temperature_2/raw",
//		OutputProcTopic: "sensor/temperature_2/proc",
//	}

var topicHandlers = map[string]TopicHandler{
	"+/sensor/temperature_1/state": {
		SensorType:      "default",
		ProcessFunc:     processTemperature1,
		OutputRawTopic:  "sensor/temperature_1/raw",
		OutputProcTopic: "sensor/temperature_1/proc",
	},
	"+/sensor/temperature_2/state": {
		SensorType:      "default",
		ProcessFunc:     processTemperature2,
		OutputRawTopic:  "sensor/temperature_2/raw",
		OutputProcTopic: "sensor/temperature_2/proc",
	},
	"+/sensor/motor_current/state": {
		SensorType:      "default",
		ProcessFunc:     processMotorCurrent,
		OutputRawTopic:  "sensor/motor_current/raw",
		OutputProcTopic: "sensor/motor_current/proc",
	},
	"+/sensor/gas_sensor/state": {
		SensorType:      "default",
		ProcessFunc:     processGasSensor,
		OutputRawTopic:  "sensor/gas_sensor/raw",
		OutputProcTopic: "sensor/gas_sensor/proc",
	},
	"+/sensor/mpu6050_temperature/state": {
		SensorType:      "default",
		ProcessFunc:     processMPU6050Temperature,
		OutputRawTopic:  "sensor/mpu6050_temperature/raw",
		OutputProcTopic: "sensor/mpu6050_temperature/proc",
	},
	"+/sensor/mpu6050_gyro_z/state": {
		SensorType:      "default",
		ProcessFunc:     processMPU6050GyroZ,
		OutputRawTopic:  "sensor/mpu6050_gyro_z/raw",
		OutputProcTopic: "sensor/mpu6050_gyro_z/proc",
	},
	"+/sensor/mpu6050_accel_z/state": {
		SensorType:      "default",
		ProcessFunc:     processMPU6050AccelZ,
		OutputRawTopic:  "sensor/mpu6050_accel_z/raw",
		OutputProcTopic: "sensor/mpu6050_accel_z/proc",
	},
	"+/sensor/mpu6050_gyro_y/state": {
		SensorType:      "default",
		ProcessFunc:     processMPU6050GyroY,
		OutputRawTopic:  "sensor/mpu6050_gyro_y/raw",
		OutputProcTopic: "sensor/mpu6050_gyro_y/proc",
	},
	"+/sensor/mpu6050_accel_y/state": {
		SensorType:      "default",
		ProcessFunc:     processMPU6050AccelY,
		OutputRawTopic:  "sensor/mpu6050_accel_y/raw",
		OutputProcTopic: "sensor/mpu6050_accel_y/proc",
	},
	"+/sensor/mpu6050_gyro_x/state": {
		SensorType:      "default",
		ProcessFunc:     processMPU6050GyroX,
		OutputRawTopic:  "sensor/mpu6050_gyro_x/raw",
		OutputProcTopic: "sensor/mpu6050_gyro_x/proc",
	},
	"+/sensor/mpu6050_accel_x/state": {
		SensorType:      "default",
		ProcessFunc:     processMPU6050AccelX,
		OutputRawTopic:  "sensor/mpu6050_accel_x/raw",
		OutputProcTopic: "sensor/mpu6050_accel_x/proc",
	},
	"+/sensor/ina226_power/state": {
		SensorType:      "default",
		ProcessFunc:     processINA226Power,
		OutputRawTopic:  "sensor/ina226_power/raw",
		OutputProcTopic: "sensor/ina226_power/proc",
	},
	"+/sensor/ina226_current/state": {
		SensorType:      "default",
		ProcessFunc:     processINA226Current,
		OutputRawTopic:  "sensor/ina226_current/raw",
		OutputProcTopic: "sensor/ina226_current/proc",
	},
	"+/sensor/ina226_shunt_voltage/state": {
		SensorType:      "default",
		ProcessFunc:     processINA226ShuntVoltage,
		OutputRawTopic:  "sensor/ina226_shunt_voltage/raw",
		OutputProcTopic: "sensor/ina226_shunt_voltage/proc",
	},
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

		serialNumber := parts[0]
		dataKey := parts[2]

		topicKey := fmt.Sprintf("+/sensor/%s/state", dataKey)
		handler, exists := topicHandlers[topicKey]
		if !exists {
			loggers.Err.Printf("No handler found for topic '%s'", msg.Topic())
			return
		}

		connectData := Connect{
			SensorType: handler.SensorType,
		}
		connectDeviceTopic := serialNumber + "/" + ConnectTopic
		if err := publishMessage(client, connectDeviceTopic, connectData, loggers); err != nil {
			loggers.Err.Printf("Failed to publish connect message: %v", err)
			return
		}

		rawData, processedData, err := handler.ProcessFunc(msg.Payload())
		if err != nil {
			loggers.Err.Printf("Error processing message: %v", err)
			return
		}

		rawData["sensorType"] = handler.SensorType
		outputRawTopic := serialNumber + "/" + handler.OutputRawTopic
		if err := publishMessage(client, outputRawTopic, rawData, loggers); err != nil {
			loggers.Err.Printf("Failed to publish raw data: %v", err)
			return
		}

		processedData["sensorType"] = handler.SensorType
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
