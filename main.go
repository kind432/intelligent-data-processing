package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	//if err := config.Init(); err != nil {
	//	log.Fatalf("%s", err.Error())
	//}
	//
	//loggers := logger.InitLogger()
	//client := mqtt.InitMQTTClient(loggers)
	//defer client.Disconnect(250)
	//
	//sigChan := make(chan os.Signal, 1)
	//signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	//<-sigChan
	topic := "test@mail.ru/sensor/ROBBO_protos_02_temperature_1/state"

	// Регулярное выражение для разбора топика
	re := regexp.MustCompile(`([^/]+)@[^/]+/sensor/([^/]+)/state`)

	// Применяем регулярное выражение
	matches := re.FindStringSubmatch(topic)
	if len(matches) != 3 {
		fmt.Printf("Failed to parse topic: '%s'. Matches length: %d\n", topic, len(matches))
		return
	}

	// Преобразуем serialNumber в формат ROBBO-protos-02
	serialParts := strings.Split(matches[2], "_")
	serialNumber := fmt.Sprintf("%s-%s-%s", serialParts[0], serialParts[1], serialParts[2])

	// Получаем dataKey как часть после последнего подчеркивания
	dataKey := strings.Join(serialParts[3:], "_")

	// Печатаем результат
	fmt.Println("serialNumber:", serialNumber)
	fmt.Println("dataKey:", dataKey)
}
