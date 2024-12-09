package mqtt

import (
	"github.com/spf13/viper"
	"intelligent-data-processing/pkg/logger"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func InitMQTTClient(loggers logger.Loggers) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(viper.GetString("mqtt_broker_address"))
	opts.SetClientID("mvp_processor")
	opts.SetUsername(viper.GetString("mqtt_username"))
	opts.SetPassword(viper.GetString("mqtt_password"))
	opts.SetDefaultPublishHandler(messageHandler(loggers))
	opts.OnConnect = func(client mqtt.Client) {
		loggers.Info.Println("Connected to MQTT Broker")
		SubscribeToTopics(client, loggers)
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
	if token.Error() != nil {
		loggers.Err.Printf("Error connecting to MQTT broker: %s", token.Error())
	}
	return client
}

func SubscribeToTopics(client mqtt.Client, loggers logger.Loggers) {
	initTopicHandlers()
	for topic := range topicHandlers {
		token := client.Subscribe(topic, 0, nil)
		token.Wait()
		if token.Error() != nil {
			loggers.Err.Printf("Error subscribing to topic '%s': %v", topic, token.Error())
		} else {
			loggers.Info.Printf("Successfully subscribed to topic: '%s'", topic)
		}
	}
}
