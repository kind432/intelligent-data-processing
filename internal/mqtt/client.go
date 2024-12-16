package mqtt

import (
	"github.com/spf13/viper"
	"intelligent-data-processing/internal/handler"
	"intelligent-data-processing/pkg/logger"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func InitMQTTClient(h handler.Handler, logger logger.Logger) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(viper.GetString("mqtt_broker_address"))
	opts.SetClientID("mvp_processor")
	opts.SetUsername(viper.GetString("mqtt_username"))
	opts.SetPassword(viper.GetString("mqtt_password"))
	opts.SetDefaultPublishHandler(h.MessageHandler())
	opts.OnConnect = func(client mqtt.Client) {
		logger.Info.Println("Connected to MQTT Broker")
		handler.SubscribeToTopics(h, client)
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
	if token.Error() != nil {
		logger.Err.Printf("Error connecting to MQTT broker: %s", token.Error())
	}
	return client
}
