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

func initTopicHandlers(h Handler) {
	baseSensorTopic := viper.GetString("mqtt_username") + "/sensor/robbo_protos_%02d_%s/state"

	sensors := []struct {
		Key        string
		RawSuffix  string
		ProcSuffix string
	}{
		{"motor_x_temperature", "sensor/motor_x_temperature/raw", "sensor/motor_x_temperature/proc"},
		{"motor_y_temperature", "sensor/motor_y_temperature/raw", "sensor/motor_y_temperature/proc"},
		{"motor_z_temperature", "sensor/motor_z_temperature/raw", "sensor/motor_z_temperature/proc"},
		{"motor_e0_temperature", "sensor/motor_e0_temperature/raw", "sensor/motor_e0_temperature/proc"},
		{"motor_e1_temperature", "sensor/motor_e1_temperature/raw", "sensor/motor_e1_temperature/proc"},
		{"table_temperature", "sensor/table_temperature/raw", "sensor/table_temperature/proc"},
		{"radiator_temperature", "sensor/radiator_temperature/raw", "sensor/radiator_temperature/proc"},
		{"air_temperature", "sensor/air_temperature/raw", "sensor/air_temperature/proc"},
		{"ina226_laser_current", "sensor/ina226_laser_current/raw", "sensor/ina226_laser_current/proc"},
		{"ina226_extruder_current", "sensor/ina226_extruder_current/raw", "sensor/ina226_extruder_current/proc"},
		{"mpu6500_carriage_accel_x", "sensor/mpu6500_carriage_accel_x/raw", "sensor/mpu6500_carriage_accel_x/proc"},
		{"mpu6500_carriage_accel_y", "sensor/mpu6500_carriage_accel_y/raw", "sensor/mpu6500_carriage_accel_y/proc"},
		{"mpu6500_carriage_accel_z", "sensor/mpu6500_carriage_accel_z/raw", "sensor/mpu6500_carriage_accel_z/proc"},
		{"mpu6500_carriage_gyro_x", "sensor/mpu6500_carriage_gyro_x/raw", "sensor/mpu6500_carriage_gyro_x/proc"},
		{"mpu6500_carriage_gyro_y", "sensor/mpu6500_carriage_gyro_y/raw", "sensor/mpu6500_carriage_gyro_y/proc"},
		{"mpu6500_carriage_gyro_z", "sensor/mpu6500_carriage_gyro_z/raw", "sensor/mpu6500_carriage_gyro_z/proc"},
		{"tachometer", "sensor/tachometer/raw", "sensor/tachometer/proc"},
		{"smoke_sensor", "sensor/smoke_sensor/raw", "sensor/smoke_sensor/proc"},
		{"input_current", "sensor/input_current/raw", "sensor/input_current/proc"},
	}

	for i := 1; i <= 8; i++ {
		for _, sensor := range sensors {
			topic := fmt.Sprintf(baseSensorTopic, i, sensor.Key)
			h.TopicHandlers[topic] = TopicHandler{
				OutputRawTopic:  sensor.RawSuffix,
				OutputProcTopic: sensor.ProcSuffix,
				HandlerFunc:     h.HandleFloatSensorData,
			}
		}

		doorStateTopic := fmt.Sprintf("%s/sensor/robbo_protos_%02d_door/state", viper.GetString("mqtt_username"), i)
		h.TopicHandlers[doorStateTopic] = TopicHandler{
			OutputRawTopic: "sensor/door/raw",
			HandlerFunc:    h.HandleDoorState,
		}
		prutok1StateTopic := fmt.Sprintf("%s/sensor/robbo_protos_%02d_prutok_1/state", viper.GetString("mqtt_username"), i)
		h.TopicHandlers[prutok1StateTopic] = TopicHandler{
			OutputRawTopic: "sensor/prutok_1/raw",
			HandlerFunc:    h.HandlePrutok1State,
		}
		prutok2StateTopic := fmt.Sprintf("%s/sensor/robbo_protos_%02d_prutok_2/state", viper.GetString("mqtt_username"), i)
		h.TopicHandlers[prutok2StateTopic] = TopicHandler{
			OutputRawTopic: "sensor/prutok_2/raw",
			HandlerFunc:    h.HandlePrutok2State,
		}
		redLedStateTopic := fmt.Sprintf("%s/switch/robbo_protos_%02d_red_led/state", viper.GetString("mqtt_username"), i)
		h.TopicHandlers[redLedStateTopic] = TopicHandler{
			OutputRawTopic: "switch/red_led/raw",
			HandlerFunc:    h.HandleRedLedState,
		}
		yellowLedStateTopic := fmt.Sprintf("%s/switch/robbo_protos_%02d_yellow_led/state", viper.GetString("mqtt_username"), i)
		h.TopicHandlers[yellowLedStateTopic] = TopicHandler{
			OutputRawTopic: "switch/yellow_led/raw",
			HandlerFunc:    h.HandleYellowLedState,
		}
		powerRelayStateTopic := fmt.Sprintf("%s/switch/robbo_protos_%02d_power_relay/state", viper.GetString("mqtt_username"), i)
		h.TopicHandlers[powerRelayStateTopic] = TopicHandler{
			OutputRawTopic: "switch/power_relay/raw",
			HandlerFunc:    h.HandlePowerRelayState,
		}
		mpuChooseStateTopic := fmt.Sprintf("%s/switch/robbo_protos_%02d_mpu_choose/state", viper.GetString("mqtt_username"), i)
		h.TopicHandlers[mpuChooseStateTopic] = TopicHandler{
			OutputRawTopic: "switch/mpu_choose/raw",
			HandlerFunc:    h.HandleMPUChooseState,
		}
	}

	h.TopicHandlers["+/switch/power_relay/command"] = TopicHandler{
		HandlerFunc: h.HandlePowerRelayCommand,
	}
	h.TopicHandlers["+/switch/red_led/command"] = TopicHandler{
		HandlerFunc: h.HandleRedLedCommand,
	}
	h.TopicHandlers["+/switch/yellow_led/command"] = TopicHandler{
		HandlerFunc: h.HandleYellowLedCommand,
	}
	h.TopicHandlers["+/switch/mpu_choose/command"] = TopicHandler{
		HandlerFunc: h.HandleMPUChooseCommand,
	}
}

func SubscribeToTopics(h Handler, client mqtt.Client) {
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
