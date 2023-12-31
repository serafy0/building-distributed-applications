package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/streadway/amqp"
)

func main() {
	amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatal(err)

	}
	defer amqpConnection.Close()
	channelAmqp, _ := amqpConnection.Channel()
	defer channelAmqp.Close()
	forever := make(chan bool)
	msgs, err := channelAmqp.Consume(os.Getenv("RABBITMQ_QUEUE"), "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for messages. to exit press CTRL+C")
	<-forever
}

type Request struct {
	URL string `json:"url"`
}

type Feed struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Link struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Thumbnail struct {
		URL string `xml:"url,attr"`
	} `xml:"thumbnail"`
	Title string `xml:"title"`
}

func GetFeedEntries(url string) ([]Entry, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	byteValue, _ := ioutil.ReadAll(resp.Body)
	var feed Feed
	xml.Unmarshal(byteValue, &feed)

	return feed.Entries, nil
}
