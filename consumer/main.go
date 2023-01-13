package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"github.com/sirupsen/logrus"
	"net/http"
	clientConsumer "practicum/consumer/client"
	"practicum/consumer/common"
	"practicum/consumer/provider"
	"practicum/consumer/repo"
	"time"
)

var (
	logger *logrus.Logger
	config repo.Config
)

func wrapJson(values []string) []byte {
	res, err := json.Marshal(values)
	if err != nil {
		return nil
	}
	return res
}

func main() {
	logger = logrus.New()
	err := confita.NewLoader(file.NewBackend("./consumer/config.json")).Load(context.Background(), &config)
	if err != nil {
		logger.Error(err)
	}

	redis, err := provider.NewRedisProvider()
	if err != nil {
		logger.Fatal(err)
	}

	repository := repo.NewRepo(config)

	kafkaClient := clientConsumer.NewKafkaClient([]string{"localhost:9092"}, "tests")
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	batch, err := kafkaClient.Read(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	//var wg sync.WaitGroup
	logger.Info(len(batch))

	//
	//for idx, rawEntry := range batch {
	//	var entry common.Entry
	//	json.Unmarshal(rawEntry, &entry)
	//	result = append(result, entry)
	//
	//	wg.Add(1)
	//	index := idx
	//	go func(index int, input *common.Entry) {
	//		defer wg.Done()
	//		processor.ProcessEntry(input)
	//	}(index, &result[index])
	//}

	go func(ctx context.Context) {
		for {
			select {
			case <-time.After(500 * time.Millisecond):
				msg, err := kafkaClient.Read(ctx)
				if err != nil {
					fmt.Println("Couldnt read from kafka" + err.Error())
				}

				var blm []common.LogMessage
				json.Unmarshal(msg, &blm)
				//for _, event := range blm {
				//	fmt.Printf("event: %v \n", event)
				//}

				for _, event := range blm {
					fmt.Printf("event: %v \n", event)
					date := event.Timestamp.Format("2006-01-01") // to yyyy-MM-dd
					fmt.Printf("date: %v \n", date)
					if x, found := event.InputParams["user_id"]; found {
						id, _ := x.(string)
						fmt.Printf("id : %v \n", id)
						err := redis.Write(ctx, time.Hour, fmt.Sprintf("%s_%d", date, 1), id)
						fmt.Errorf(" error write %v \n", err)
					}

				}
				//err = repository.AddValue(ctx, string(msg))
				//if err != nil {
				//	fmt.Println("Couldnt write kafka msg to DB" + err.Error())
				//	// TODO: dead letter queue

			case <-ctx.Done():
				fmt.Println("Done listening")
			}
		}
	}(ctx)
	http.HandleFunc("/list", func(writer http.ResponseWriter, request *http.Request) {
		values, err := repository.ReadValues(request.Context())
		if err != nil {
			fmt.Println("isnt ok")
			writer.WriteHeader(http.StatusInternalServerError)
			_, err := writer.Write([]byte("Internal server error"))
			if err != nil {
				return
			}
		}
		//writer.WriteHeader(http.StatusOK)
		_, err = writer.Write(wrapJson(values))
		if err != nil {
			return
		}

	})
	err = http.ListenAndServe(":8091", nil)
	if err != nil {
		logger.Error(err)
	}
}
