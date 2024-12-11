package handler

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

type TopicHandler struct {
	OutputRawTopic  string
	OutputProcTopic string
	HandlerFunc     func(mqtt.Client, mqtt.Message)
}

func initTopicHandlers(h *Handler) {
	baseSensorTopic := viper.GetString("mqtt_username") + "/sensor/ROBBO_protos_%02d_%s/state"

	sensors := []struct {
		Key        string
		RawSuffix  string
		ProcSuffix string
	}{
		{"temperature_1", "sensor/temperature_1/raw", "sensor/temperature_1/proc"},
		{"temperature_2", "sensor/temperature_2/raw", "sensor/temperature_2/proc"},
		{"gas_sensor", "sensor/gas_sensor/raw", "sensor/gas_sensor/proc"},
		{"mpu6050_accel_x", "sensor/mpu6050_accel_x/raw", "sensor/mpu6050_accel_x/proc"},
		{"mpu6050_accel_y", "sensor/mpu6050_accel_y/raw", "sensor/mpu6050_accel_y/proc"},
		{"mpu6050_accel_z", "sensor/mpu6050_accel_z/raw", "sensor/mpu6050_accel_z/proc"},
		{"mpu6050_gyro_x", "sensor/mpu6050_gyro_x/raw", "sensor/mpu6050_gyro_x/proc"},
		{"mpu6050_gyro_y", "sensor/mpu6050_gyro_y/raw", "sensor/mpu6050_gyro_y/proc"},
		{"mpu6050_gyro_z", "sensor/mpu6050_gyro_z/raw", "sensor/mpu6050_gyro_z/proc"},
		{"ina226_current", "sensor/ina226_current/raw", "sensor/ina226_current/proc"},
		{"ina226_power", "sensor/ina226_power/raw", "sensor/ina226_power/proc"},
		{"ina226_shunt_voltage", "sensor/ina226_shunt_voltage/raw", "sensor/ina226_shunt_voltage/proc"},
		{"motor_current", "sensor/motor_current/raw", "sensor/motor_current/proc"},
	}

	for i := 1; i <= 8; i++ {
		for _, sensor := range sensors {
			topic := fmt.Sprintf(baseSensorTopic, i, sensor.Key)
			h.TopicHandlers[topic] = TopicHandler{
				OutputRawTopic:  sensor.RawSuffix,
				OutputProcTopic: sensor.ProcSuffix,
				HandlerFunc:     h.HandleSensorData,
			}
		}

		switchStateTopic := fmt.Sprintf("%s/switch/ROBBO_protos_%02d_power_relay/state", viper.GetString("mqtt_username"), i)
		h.TopicHandlers[switchStateTopic] = TopicHandler{
			OutputRawTopic: "switch/power_relay/raw",
			HandlerFunc:    h.HandlePowerRelayState,
		}
	}

	h.TopicHandlers["+/switch/power_relay/command"] = TopicHandler{
		HandlerFunc: h.HandlePowerRelayCommand,
	}
}

func SubscribeToTopics(h *Handler, client mqtt.Client) {
	initTopicHandlers(h)
	for topic := range h.TopicHandlers {
		token := client.Subscribe(topic, 0, nil)
		token.Wait()
		if token.Error() != nil {
			h.Logger.Err.Printf("Error subscribing to topic '%s': %v", topic, token.Error())
		} else {
			h.Logger.Info.Printf("Successfully subscribed to topic: '%s'", topic)
		}
	}
}
