package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"

	"intelligent-data-processing/pkg/logger"
)

func NewClient(log logger.Logger, h *Handler) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(viper.GetString("mqtt_broker_address"))
	opts.SetClientID("mvp_processor")
	opts.SetUsername(viper.GetString("mqtt_username"))
	opts.SetPassword(viper.GetString("mqtt_password"))

	opts.SetDefaultPublishHandler(h.UniversalHandler)

	opts.OnConnect = func(client mqtt.Client) {
		log.Info.Println("Connected to MQTT Broker")

		username := viper.GetString("mqtt_username")
		topics := []string{
			// Data from devices
			username + "/sensor/+/state",
			username + "/binary_sensor/+/state",
			username + "/switch/+/state",

			// Listen external commands
			"+/switch/+/command",
		}

		for _, topic := range topics {
			client.Subscribe(topic, 0, nil)
			log.Info.Printf("Subscribed to: %s", topic)
		}
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Err.Fatalf("Error connecting to MQTT broker: %s", token.Error())
	}
	return client
}
