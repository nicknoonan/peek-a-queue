package main

// import (
// 	"github.com/aws/aws-sdk-go-v2/service/sqs"
// )

type item struct {
	name         string
	lengthString string
	url string
}


func (i item) Title() string       { return i.name }
func (i item) Description() string { 
	if i.lengthString != "" {
		return i.lengthString
	}
	return "-"
}

func (i item) FilterValue() string { return i.name }