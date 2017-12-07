package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	BaseURL        = "https://api.nicehash.com/api"
	AlgoSha256     = "1"
	LocationEurope = 0
	LocationUSA    = "1"
)

// NiceHash es una envoltura para el API de NiceHash
type NiceHash struct {
	ApiId  string
	ApiKey string
}

type Response struct {
	Method string
	Result struct {
		Orders Orders
	}
}

type Orders []Order

func (or *Orders) Get(id int) Order {
	for _, o := range *or {
		if o.ID == id {
			return o
		}
	}
	return Order{}
}

// FilteredOrders retorna las ordenes que tienen workers
// y que el tipo de orden es Standard
func (or *Orders) FilteredOrders() (Orders, int) {
	orders := *or
	filtered := Orders{}

	for _, order := range orders {
		if order.Workers > 0 && order.Type == 0 {
			filtered = append(filtered, order)
		}
	}

	return filtered, len(filtered)
}

func parsePrice(p string) float32 {
	n, _ := strconv.ParseFloat(p, 32)
	return float32(n)
}

// MaxPrice retorna el precio de la oferta más alta
func (or *Orders) MaxPrice() float32 {
	filtered, _ := or.FilteredOrders()
	price := parsePrice(filtered[0].Price)
	// price := filtered[0].Price
	return float32(price)
}

// MinPrice retorna el precio de oferta más baja
func (or *Orders) MinPrice() float32 {
	filtered, n := or.FilteredOrders()
	price := parsePrice(filtered[n-1].Price)
	// price := filtered[n-1].Price
	return float32(price)
}

// Order representa una orden de minería
type Order struct {
	LimitSpeed    string `json:"limit_speed"`
	Alive         bool
	Price         string
	ID            int
	Type          int
	Workers       int
	Algo          int
	AcceptedSpeed string `json:"accepted_speed"`
}

// NewNiceHash retorna un nuevo NH
func NewNiceHash(id, key string) *NiceHash {
	return &NiceHash{id, key}
}

func (n *NiceHash) GetApiVersion() {
	r, _ := http.Get(BaseURL)
	io.Copy(os.Stdout, r.Body)
}

func (n *NiceHash) getAuthValues() url.Values {
	v := url.Values{}
	v.Set("id", n.ApiId)
	v.Set("key", n.ApiKey)
	return v
}

func Send(v url.Values) *http.Response {
	r, _ := http.Get(fmt.Sprintf("%s?%s", BaseURL, v.Encode()))
	return r
}

// DecreaseBid decrementa la oferta de una orden
func (n *NiceHash) DecreaseBid(o int) {
	v := n.getAuthValues()
	v.Set("method", "orders.set.price.decrease")
	v.Set("location", "1")
	v.Set("algo", "1")
	v.Set("order", strconv.FormatInt(int64(o), 10))
	r := Send(v)
	io.Copy(os.Stdout, r.Body)
}

// IncreaseBid incrementa la oferta de una orden
func (n *NiceHash) IncreaseBid(o int, p float32) {
	price := float64(p)
	v := n.getAuthValues()
	v.Set("method", "orders.set.price")
	v.Set("location", "1")
	v.Set("algo", "1")
	v.Set("order", strconv.FormatInt(int64(o), 10))
	v.Set("price", strconv.FormatFloat(price, 'f', 4, 32))

	// fmt.Println(strconv.FormatFloat(p, 'f', 4, 32))
	r := Send(v)
	io.Copy(os.Stdout, r.Body)
}

// GetAllOrders trae todas las ordenes del market place
func (n *NiceHash) GetAllOrders() Orders {
	v := n.getAuthValues()
	v.Set("method", "orders.get")
	v.Set("location", "1")
	v.Set("algo", "1")

	r := Send(v)

	resp := Response{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&resp)

	return resp.Result.Orders
}

// GetOrders trae las ordenes del usuario
func (n *NiceHash) GetOrders() Orders {
	v := n.getAuthValues()
	v.Set("method", "orders.get")
	v.Set("location", "1")
	v.Set("algo", "1")
	v.Set("my", "")

	r := Send(v)

	resp := Response{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&resp)

	return resp.Result.Orders
}
