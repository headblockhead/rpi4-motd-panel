package main

import (
	"log/slog"
	"os"

	"github.com/a-h/character"
	//mqtt "github.com/eclipse/paho.mqtt.golang"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	_, err := host.Init()
	if err != nil {
		log.Error("error initializing periph", slog.Any("error", err))
		return
	}
	bus, err := i2creg.Open("")
	if err != nil {
		log.Error("error opening i2c", slog.Any("error", err))
		return
	}
	dev := &i2c.Dev{
		Bus:  bus,
		Addr: 0x27,
	}
	d := character.NewDisplay(dev, false)

	//textToDraw := make(chan string)

	/* options := mqtt.NewClientOptions()*/
	/*options.AddBroker(fmt.Sprintf("tcp://172.16.3.41:1883"))*/
	/*options.SetClientID("rpi4-motd-panel")*/
	/*// TODO*/
	/*options.SetUsername("user")*/
	/*options.SetPassword("pass")*/
	/*options.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {*/
	/*if msg.Topic() == "home-assistant/signage/control" {*/
	/*textToDraw <- string(msg.Payload())*/
	/*}*/
	/*})*/

	/*client := mqtt.NewClient(options)*/
	/*if token := client.Connect(); token.Wait() && token.Error() != nil {*/
	/*log.Error("error connecting to mqtt", slog.Any("error", token.Error()))*/
	/*return*/
	/*}*/
	/*if token := client.Subscribe("home-assistant/signage/control", 0, nil); token.Wait() && token.Error() != nil {*/
	/*log.Error("error subscribing to topic", slog.Any("error", token.Error()))*/
	/*return*/
	/*}*/
	/*if token := client.Publish(client, "home-assistant/signage/availability", 1, "online", false); token.Wait() && token.Error() != nil {*/
	/*log.Error("error publishing availability", slog.Any("error", token.Error()))*/
	/*return*/
	/*}*/

	for {
		//d.WriteInstruction()
		d.SetBacklight(false)
		d.Goto(0, 0)
		/*d.Print("SYSTEM ONLINE")*/
		/*time.Sleep(5 * time.Second)*/
		/*d.Clear()*/
		/*d.SetBacklight(false)*/
		/*time.Sleep(25 * time.Second)*/
	}
}
