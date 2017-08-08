# cfn
A simple command-line tool to tail human-readable [AWS CloudFormation](https://aws.amazon.com/cloudformation/) events suitable for observing stack creation and updates in real-time. This may be useful when developing and testing CloudFormation templates, or when monitoring infrastructure deployments.
## Install
```
go get github.com/rpgreen/cfn
go build
```
(Optional) Configure aws cli
```
aws configure
```
## Use
```
Usage: cfn [-c command] [-s stackname] [-r region] [-p profile]
  -c string
    	Command to use (i.e. 'tail') (default "tail")
  -p string
    	AWS SDK profile name to use (optional)
  -r string
    	AWS region to use (optional)
  -s string
    	Stack to use (optional)
```
## Examples
### Tail events from the most recently updated stack using default AWS credentials and region
```
./cfn
```
### Tail events using specified stack and region, and using specified AWS credentials profile
```
./cfn -s mystack -p myprofile -r us-east-1
```
<img src="https://github.com/rpgreen/cfn/blob/master/ss.png" width="769" height="176"/>

By default, cfn will look for credentials using the default credential provider chain used by the [AWS CLI](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#config-settings-and-precedence). Similarly, the default region is based on the AWS CLI configuration.

To provide a custom credential profile, use the -p option. To override the default region, use the -r option.

If no stack is specified with -s option, cfn will attempt to find the most recently updated stack.
