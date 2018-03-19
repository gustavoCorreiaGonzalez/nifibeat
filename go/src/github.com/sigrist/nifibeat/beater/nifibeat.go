package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/sigrist/nifibeat/config"
)

type Nifibeat struct {
	done          chan struct{}
	config        config.Config
	client        publisher.Client
	lastIndexTime time.Time
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Nifibeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Nifibeat) request(url string, method string, beatname string) bool {
	response := RequestNifi(url, method)
	result := JSONConvert(response)

	if len(result.ProcessGroups) == 0 {
		return false
	}

	for process := range result.ProcessGroups {
		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       beatname,
			"URL":        url,
			"Nifi":       result.ProcessGroups[process],
		}
		bt.client.PublishEvent(event)

		println("Evento Enviado ao Kibana")

		urlRecursiva := fmt.Sprintf("http://nifi.armen.domrock.com.br/nifi-api/process-groups/%s/process-groups", result.ProcessGroups[process].ID)

		bt.request(urlRecursiva, method, beatname)
	}

	return true
}

func (bt *Nifibeat) Run(b *beat.Beat) error {
	logp.Info("lsbeat is running! Hit CTRL-C to stop it.")
	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	for {
		now := time.Now()
		bt.request(bt.config.URL, bt.config.Method, b.Name) // call listDir
		bt.lastIndexTime = now                              // mark Timestamp
		logp.Info("Event sent")
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
	}
}

func (bt *Nifibeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
