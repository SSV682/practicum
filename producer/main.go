package main

import (
	"context"
	_ "context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"practicum/producer/client"
	"practicum/producer/common"
	"practicum/producer/common/serializer"
	"strconv"
	"time"
)

func main() {
	brokers := []string{"127.0.0.1:9092"}
	kafkaClient := client.NewKafkaClient(brokers, "facade_recommendation.logs.request")

	jsonSerializer := serializer.NewJsonSerializer()

	http.HandleFunc("/entryList", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			params := request.URL.Query()

			for j := 0; j < 100; j++ {

				scendrioID, _ := strconv.Atoi(params["scenario_id"][0])
				var entries []common.LogMessage
				ss, _ := params["user_id"]

				for i := 0; i < 100; i++ {

					ip := make(map[string]any)
					ip["user_id"] = strconv.Itoa(i) + ss[0]
					message := generateMessage(ip, scendrioID+j)
					entries = append(entries, message)
				}

				go func() {
					batch := make([]common.LogMessage, 0, 10)
					for _, entry := range entries {
						batch = append(batch, entry)
					}
					log.Infof("batch: %v", batch)
					data, _ := jsonSerializer.Serialize(batch)
					err := kafkaClient.SendMessage(context.Background(), data)
					if err != nil {
						log.Errorf("send message: %v", err)
					}
				}()

				_, err := writer.Write([]byte("Record sent"))
				if err != nil {
					log.Error(err)
				}
			}
		}
	})

	err := http.ListenAndServe(":8091", nil)
	if err != nil {
		log.Error(err)
	}
}

func generateMessage(ip map[string]any, scendrioID int) common.LogMessage {
	tn := time.Now()
	return common.LogMessage{
		Host:            "facade-recommendation-ift-master.ift.amazmetest.ru",
		RequestHost:     "192.168.0.1",
		Endpoint:        "/personal/items/",
		EndpointVersion: "v1",
		ArtifactId:      2533,
		ItemCount:       20,
		OutputCount:     19,
		Mode:            "offline",
		StatusCode:      200,
		Latency:         100,
		Timestamp:       tn,
		InputParams:     ip,
		ScenarioId:      scendrioID,
	}
}
