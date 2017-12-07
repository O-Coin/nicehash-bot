package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

const (
	OrderId  = 5156753
	MaxPrice = 0.1950
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Se cambio la configuración")
	})

	if err != nil {
		log.Fatalln("No se pudo cargar el archivo de configuración")
	}

	nh := NewNiceHash(viper.GetString("ApiId"), viper.GetString("ApiKey"))

	for {
		orders := nh.GetAllOrders()
		order := orders.Get(OrderId)
		orderPrice, _ := strconv.ParseFloat(order.Price, 32)

		activeOrders, _ := orders.FilteredOrders()
		minPrice := activeOrders.MinPrice()
		maxPrice := activeOrders.MaxPrice()

		if order.Workers < 1 {
			inc := minPrice + 0.0002
			nh.IncreaseBid(OrderId, inc)
			fmt.Printf("Incremento a %1.4f\n", inc)
			time.Sleep(time.Minute * 5)
		} else if x := float32(orderPrice); x-0.0010 > minPrice {
			nh.DecreaseBid(OrderId)
			fmt.Printf("Decrementó a %1.4f\n", x)
			time.Sleep(time.Minute * 11)
		} else {
			fmt.Printf("Precio actual: %1.4f\n", orderPrice)
			fmt.Printf("Precio mínimo: %1.4f\n", minPrice)
			fmt.Printf("Precio máximo: %1.4f\n", maxPrice)
			print("**************************************")
			time.Sleep(time.Minute * 1)
		}

	}
}
