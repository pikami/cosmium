package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/pikami/cosmium/internal/rntbd"
)

func main() {
	input := flag.String("input", "", "Input hex string")
	isResponse := flag.Bool("response", false, "Is response")
	flag.Parse()

	data, err := hex.DecodeString(*input)
	if err != nil {
		fmt.Printf("Error decoding hex string: %v\n", err)
		return
	}

	frame, err := rntbd.ParseFrame(data, *isResponse)
	if err != nil {
		fmt.Printf("Error parsing frame: %v\n", err)
		return
	}

	fmt.Printf("Activity ID: %s\n", hex.EncodeToString(frame.ActivityId))
	fmt.Printf("Resource Type: %s\n", frame.ResourceType.String())
	fmt.Printf("Operation Type: %s\n", frame.OperationType.String())

	if len(frame.RequestHeaders) > 0 {
		fmt.Printf("=== Request Headers ===\n")
		for header, value := range frame.RequestHeaders {
			fmt.Printf("%s: %v\n", header.String(), value)
		}
	}

	if len(frame.ResponseHeaders) > 0 {
		fmt.Printf("=== Response Headers ===\n")
		for header, value := range frame.ResponseHeaders {
			fmt.Printf("%s: %v\n", header.String(), value)
		}
	}

	if len(frame.ContextHeaders) > 0 {
		fmt.Printf("=== Context Headers ===\n")
		for header, value := range frame.ContextHeaders {
			fmt.Printf("%s: %v\n", header.String(), value)
		}
	}

	if len(frame.Payload) > 0 {
		var jsonObj any
		err := json.Unmarshal(frame.Payload, &jsonObj)
		if err != nil {
			fmt.Printf("Payload: %s\n", hex.EncodeToString(frame.Payload))
		} else {
			fmt.Printf("Payload: %+v\n", jsonObj)
		}
	}
}
