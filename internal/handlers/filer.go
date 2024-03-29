package handlers

// import (
// 	"encoding/json"
// 	"os"
// )

// type (
// 	Event struct {
// 		UUID       uint   `json:"uud"`
// 		ShortURL   string `json:"short_url"`
// 		OriginlURL string `json:"original_url"`
// 	}
// 	Producer struct {
// 		file    *os.File
// 		encoder *json.Encoder
// 	}
// 	Consumer struct {
// 		file    *os.File
// 		decoder *json.Decoder
// 	}
// )

// func (p *Producer) Close() error {
// 	return p.file.Close()
// }

// func NewProducer(fileName string) (*Producer, error) {

// 	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Producer{
// 		file:    file,
// 		encoder: json.NewEncoder(file),
// 	}, nil
// }

// func (p *Producer) WriteEvent(event *Event) error {
// 	return p.encoder.Encode(&event)
// }

// func NewConsumer(fileName string) (*Consumer, error) {

// 	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Consumer{
// 		file:    file,
// 		decoder: json.NewDecoder(file),
// 	}, nil
// }

// // func (c *Consumer) ReadEvent() (*Event, error) {
// // 	event := &Event{}
// // 	if err := c.decoder.Decode(&event); err != nil {
// // 		return nil, err
// // 	}

// // 	return event, nil
// // }
