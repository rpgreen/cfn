# cfn
A simple command-line tool to tail human-readable [AWS CloudFormation](https://aws.amazon.com/cloudformation/) events. This may be useful when developing and testing CloudFormation templates, or when monitoring infrastructure deployments.
##Install
```
go get github.com/rpgreen/cfn
go build
```
(Optional) Configure aws cli
```
aws configure
```

##Use
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
i.e.
```
./cfn -s mystack -p rpgreen -r us-east-1
```
![screenshot](https://github.com/rpgreen/cfn/blob/master/ss.png)

By default, cfn will use the default credential provider chain and shared region config used by the [AWS CLI](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#config-settings-and-precedence).

To provide a custom credential profile, use the -p option. To override the default region, use the -r option.

If no stack is specified with -s option, cfn will attempt to find the most recently updated stack.