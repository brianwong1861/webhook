package main

import (
	"fmt"
	"os"
	"strings"
	"webhook/kit"
)

func (c *BucketMessage) LogNameExtract() string {
	logArr := strings.Split(c.Key, "/")
	logName := logArr[len(logArr)-1]
	return logName
}

type LogProperty struct {
	EphemeralPath string
	Name          string
}

func (l *LogProperty) CreateFile() *os.File {
	file, err := os.Create(l.EphemeralPath + l.Name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return file
}
func LogParser(s3 *BucketMessage) error {
	fName := s3.LogNameExtract()
	p := &LogProperty{
		EphemeralPath: "/tmp/",
		Name:          fName,
	}
	fObj := p.CreateFile()
	err := kit.FetchFromS3(fObj, s3.Bucket, s3.Key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = kit.LogEntryParser(p.EphemeralPath + p.Name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fObj.Close()
	// Housekeeping for log under /tmp directory
	postProcess(p.EphemeralPath + p.Name)
	return nil
}

func postProcess(filepath string) {
	err := os.Remove(filepath)
	if err != nil {
		fmt.Println(err)
	}
}
