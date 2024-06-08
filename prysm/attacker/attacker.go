package attacker

import (
	"github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/attackclient"
	"os"
)

var (
	serviceUrl string
	client     *attackclient.Client
)

func InitAttacker(url string) error {
	env := os.Getenv("ATTACKER_SERVICE_URL")
	if url != "" {
		env = url
	}
	serviceUrl = env
	logrus.WithField("url", serviceUrl).Info("Attacker service init")
	return nil
}

func GetAttacker() *attackclient.Client {
	if client != nil {
		return client
	}
	if serviceUrl == "" {
		return nil
	}

	c, err := attackclient.Dial(serviceUrl)
	if err != nil {
		return nil
	}
	client = c
	return client
}
