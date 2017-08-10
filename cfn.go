package cfn

import (
    "fmt"
    "strings"
    "os"
    "time"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudformation"
    "flag"
    "sort"
    "github.com/mgutz/ansi"
)

func init() {
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage: cfn [-c command] [-s stackname] [-r region] [-p profile]\n")
        flag.PrintDefaults()
    }
}

func main() {
    command := flag.String("c", "tail", "Command to use (i.e. 'tail')")
    stackname := flag.String("s", "", "Stack to use (optional)")
    profile := flag.String("p", "", "AWS SDK profile name to use (optional)")
    region := flag.String("r", "", "AWS region to use (optional)")

    flag.Parse()

    cfn := createCloudFormationClient(profile, region)

    switch *command {
    case "tail":
        var stack = *stackname

        tail(*cfn, stack)
    default:
        exitWithHelp()
    }
}

func createCloudFormationClient(profile *string, region *string) *cloudformation.CloudFormation {
    options := session.Options{
        SharedConfigState: session.SharedConfigDisable,
    }
    if *region != "" {
        fmt.Println("Using region", *region)
        options.Config = aws.Config{Region: aws.String(*region)}
    }
    if *profile != "" {
        fmt.Println("Using profile", *profile)
        options.Profile = *profile
    }
    if *profile == "" && *region == "" {
        options.SharedConfigState = session.SharedConfigEnable
    }
    sess := session.Must(session.NewSessionWithOptions(options))
    cfn := cloudformation.New(sess)
    return cfn
}

func tail(cfn cloudformation.CloudFormation, stack string) {
    if stack == "" {
        fmt.Println("No stack parameter specified, selecting stack with most recent events.")
        stack = getLastUpdatedStack(cfn)
    }
    fmt.Printf("Tailing stack %s...\n", stack)
    var lasttimestamp = time.Unix(0, 0)

    for {
        res := describeStackEvents(cfn, stack)
        events := filterEvents(res.StackEvents, lasttimestamp)

        sort.Sort(ByTimeStampAscending(events))

        PrintEvents(events)
        time.Sleep(5 * time.Second)

        lasttimestamp = *res.StackEvents[0].Timestamp
    }
}

type ByTimeStampAscending []*cloudformation.StackEvent

func (a ByTimeStampAscending) Len() int           { return len(a) }
func (a ByTimeStampAscending) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimeStampAscending) Less(i, j int) bool { return a[i].Timestamp.Sub(*a[j].Timestamp) < 0 }

func filterEvents(events []*cloudformation.StackEvent, lasttime time.Time) []*cloudformation.StackEvent {
    var filtered []*cloudformation.StackEvent

    for _, event := range events {
        if event.Timestamp.Sub(lasttime) > 0 {
            filtered = append(filtered, event)
        }
    }
    return filtered
}

func isDeleted(status string) bool {
    return strings.Contains(status, "DELETE_COMPLETE")
}

func isInProgress(status string) bool {
    return strings.Contains(status, "IN_PROGRESS")
}

func getStackNames(cfn cloudformation.CloudFormation) []string {
    var stacknames []string
    res, err := cfn.ListStacks(&cloudformation.ListStacksInput{})
    if err != nil {
        fmt.Printf("%+v\n", err)
        os.Exit(1)
    }
    for _, stack := range res.StackSummaries {
        if !isDeleted(*stack.StackStatus) {
            stacknames = append(stacknames, *stack.StackName)
        }
    }
    return stacknames
}

func getLastUpdatedStack(cfn cloudformation.CloudFormation) string {
    res, err := cfn.ListStacks(&cloudformation.ListStacksInput{})
    if err != nil {
        fmt.Printf("%+v\n", err)
        os.Exit(1)
    }
    maxtime := time.Unix(0, 0)
    var maxstack string

    for _, stack := range res.StackSummaries {
        if !isInProgress(*stack.StackStatus) {
            continue
        }
        stackname := *stack.StackName
        events := describeStackEvents(cfn, stackname).StackEvents
        if events[0].Timestamp.Sub(maxtime) > 0 {
            maxtime = *events[0].Timestamp
            maxstack = stackname
        }
    }

    if maxstack == "" {
        errormsg(fmt.Sprintf("Could not find an active stack with recent events. "+
            "Please pass a stack name as an argument. Available stacks: %+v", getStackNames(cfn)))
        os.Exit(1)
    }

    return maxstack
}
func errormsg(msg string) {
    fmt.Println(ansi.Color(msg, "red+b"))
}

func describeStackEvents(cfn cloudformation.CloudFormation,
    stack string) cloudformation.DescribeStackEventsOutput {
    res, err :=cfn.DescribeStackEvents(
        &cloudformation.DescribeStackEventsInput{
            StackName: aws.String(stack),
        })

    // todo: handle deleted
    if err != nil {
        fmt.Printf("%+v\n", err)
        os.Exit(1)
    }

    return *res
}

func exitWithHelp() {
    flag.Usage()
    os.Exit(1)
}
