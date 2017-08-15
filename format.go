package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudformation"
    "github.com/olekukonko/tablewriter"
	tm "github.com/buger/goterm"
	"github.com/mgutz/ansi"
	"os"
	"strings"
	"github.com/dustin/go-humanize"
)

var red = ansi.ColorCode("red+b")
var green = ansi.ColorCode("green+h:black")
var yellow = ansi.ColorCode("yellow")
var reset = ansi.ColorCode("reset")

var writer = os.Stdout
var eventmap map[string]*cloudformation.StackEvent = make(map[string]*cloudformation.StackEvent)
var stackname string
var stackevent *cloudformation.StackEvent
var table = tablewriter.NewWriter(tm.Output)

func PrintEventsAsTable(events []*cloudformation.StackEvent) {
	if len(events) == 0 {
		return
	}
	updateEventMap(events)
	printEventTable()
}

func printEventTable() {
	data := make([][]string, len(eventmap))

	i := 0
	for id, event := range eventmap {
		if id == stackname {
			stackevent = event
		} else {
			data[i] = []string{id, getTimestamp(event), getStatusString(*event.ResourceStatus)}
			i++
		}
		// todo: sort data by resource name or status
	}

	//tm.Clear()
	tm.MoveCursor(1, 1)
	table.SetHeader([]string{"Resource", "Time", "Status"})
	table.SetFooter([]string{stackname, getTimestamp(stackevent), getStatusString(*stackevent.ResourceStatus)})
	table.SetBorder(false)
	table.SetAutoFormatHeaders(false)
	table.AppendBulk(data)
	table.Render()
	tm.Flush()
}

func getStatusString(status string) string {
	return fmt.Sprintf("%s%s%s", getColor(status), status, reset)
}

func getTimestamp(event *cloudformation.StackEvent) string {
	return humanize.Time(*event.Timestamp)
	//return event.Timestamp.Local().Format(time.UnixDate)
}

func updateEventMap(events []*cloudformation.StackEvent) {
	stackname = *events[0].StackName
	for _, event := range events {
		eventmap[*event.LogicalResourceId] = event
	}
}

func PrintEventsAsLog(events []*cloudformation.StackEvent) {
	for _, event := range events {
		printEventLine(*event)
	}
}

func printEventLine(event cloudformation.StackEvent) {
	status := *event.ResourceStatus
	printEventColor(status)
	printEvent(event)

	printReset()
}

func printEvent(event cloudformation.StackEvent) {
	timestamp := getTimestamp(&event)
	fmt.Fprintf(writer, "%s - %s - %s", *event.LogicalResourceId, timestamp, *event.ResourceStatus)
	if event.ResourceStatusReason != nil {
		fmt.Fprintf(writer, " - %s", *event.ResourceStatusReason)
	}
}

func printReset() (int, error) {
	return fmt.Fprintf(writer, "%s\n", reset)
}

func printEventColor(status string) {
	color := getColor(status)
	fmt.Fprintf(writer, "%s", color)
}

func getColor(status string) string {
	var color string
	if isInProgress(status) {
		color = yellow
	} else if isFailed(status) {
		color = red
	} else {
		color = green
	}
	return color
}

func isFailed(status string) bool {
	return strings.Contains(status, "ROLLBACK") || strings.Contains(status, "FAIL")
}
