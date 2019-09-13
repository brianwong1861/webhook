package kit

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetParams(regEx, logEntry string) (paramsMap map[string]interface{}) {
	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(logEntry)
	paramsMap = make(map[string]interface{})
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) && name != "" {
			paramsMap[name] = match[i]
		}
	}
	/* Transform "request_processing_time","response_processing_time",
	"target_processing_time" data type from interface -> string -> float64 */
	fieldFloat := []string{"request_processing_time", "response_processing_time", "target_processing_time"}
	for _, field := range fieldFloat {
		v, ok := paramsMap[field].(string)
		if !ok {
			fmt.Println("Assertion Error")
			os.Exit(1)
		}
		vFloat, err := strconv.ParseFloat(v, 64)
		if err != nil {
			fmt.Println(err)
		}
		paramsMap[field] = vFloat
	}
	fieldTrim := []string{"client_ip", "target_ip"}
	for _, field := range fieldTrim {
		str, ok := paramsMap[field].(string)
		if !ok {
			fmt.Println("Assertion Error")
			os.Exit(1)
		}
		strTrim := strings.Split(str, ":")[0]
		paramsMap[field] = strTrim
	}
	// End transform block
	return paramsMap
}
func LogEntryParser(filePath string) error {
	pattern := `(?P<ts>[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}.\d+Z) (?P<elb>\w+) (?P<client_ip>(?:[0-9]{1,3}\.){3}[0-9]{1,3}:\d+) (?P<target_ip>(?:[0-9]{1,3}\.){3}[0-9]{1,3}:\d+) (?P<request_processing_time>[+-]?((\d+\.?\d*)|(\.\d+))) (?P<target_processing_time>[+-]?((\d+\.?\d*)|(\.\d+))) (?P<response_processing_time>[+-]?((\d+\.?\d*)|(\.\d+))) (?P<elb_status_code>[0-9]{3}) (?P<target_status_code>[0-9]{3}) (?P<received_bytes>\d+) (?P<sent_bytes>\d+) \"(?P<req_method>[A-Z]\S+) (?P<req_url>.*) (?P<http_ver>\S+)\" \"(?P<ua>.*)\" (?P<ssl_cipher>[A-Z0-9-]*) (?P<ssl_proto>[A-Za-z0-9.]*)`
	f1, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f1.Close()
	scanner := bufio.NewScanner(f1)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	mapParams := make(map[string]interface{})
	for scanner.Scan() {
		go func() {
			mapParams = GetParams(pattern, scanner.Text())
			err := SinkLog(mapParams)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = SinkLog(mapParams)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}()
	}
	return nil
}
