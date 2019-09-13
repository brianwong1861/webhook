package kit

import (
	"fmt"
	"github.com/fluent/fluent-logger-golang/fluent"
	"os"
	"strconv"
)

func SinkLog(data map[string]interface{}) error {
	fluentPort, err := strconv.Atoi(os.Getenv("FLUENTD_PORT"))
	if err != nil {
		fmt.Printf("fluentd port parsing error %v", err)
	}
	logger, err := fluent.New(fluent.Config{
		FluentHost: os.Getenv("FLUENTD_HOST"),
		FluentPort: fluentPort,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Close()
	tag := os.Getenv("TAG")
	err = logger.Post(tag, data)
	// error := logger.PostWithTime(tag, time.Now(), data)
	if err != nil {
		panic(err)
	}
	return nil
}
