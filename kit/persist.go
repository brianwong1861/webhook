package kit

import (
	"encoding/json"
	"fmt"
	"os"
)

func PersistLog(data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	jFile, err := os.Create("/tmp/logdata.json")
	if err != nil {
		panic(err)
	}
	defer jFile.Close()
	_, err = jFile.Write(jsonData)
	if err != nil {
		panic(err)
	}
	fmt.Println("Data written to ", jFile.Name())
	return nil
}
