package processor

import (
	log "github.com/sirupsen/logrus"
	"practicum/consumer/common"
	"time"
)

func ProcessEntry(entry *common.Entry) {
	time.Sleep(400 * time.Millisecond)
	log.Info("value ", entry.Value)
}
