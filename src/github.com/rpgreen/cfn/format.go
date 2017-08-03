package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"time"
	//"text/tabwriter"
	"github.com/mgutz/ansi"
	"os"
	"strings"
)

var red = ansi.ColorCode("red+b")
var green = ansi.ColorCode("green+h:black")
var yellow = ansi.ColorCode("yellow")
var reset = ansi.ColorCode("reset")

//var error = ansi.ColorCode("red+b:white+h")
//var info = ansi.ColorCode("green+b:white+h")
//var writer = tabwriter.NewWriter(os.Stdout, 30, 10, 0, ' ', tabwriter.Debug)
var writer = os.Stdout

func PrintEvents(events []*cloudformation.StackEvent) {
	for _, event := range events {
		printEvent(*event)
	}
}

func printEvent(event cloudformation.StackEvent) {
	status := *event.ResourceStatus
	var color string
	if isInProgress(status) {
		color = yellow
	} else if isFailed(status) {
		color = red
	} else {
		color = green
	}

	fmt.Fprintf(writer, "%s", color)

	// todo: improved formatting
	timestamp := event.Timestamp.Local().Format(time.UnixDate)
	fmt.Fprintf(writer, "%s - %s - %s", *event.LogicalResourceId, timestamp, *event.ResourceStatus)
	if event.ResourceStatusReason != nil {
		fmt.Fprintf(writer, " - %s", *event.ResourceStatusReason)
	}

	fmt.Fprintf(writer, "%s\n", reset)
	//writer.Flush()
}

func isFailed(status string) bool {
	return strings.Contains(status, "ROLLBACK") || strings.Contains(status, "FAIL")
}
