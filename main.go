package main

import (
	"asset/sync/consul"
	"asset/sync/config"
	"flag"
	"fmt"

	Aws "asset/sync/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	var cloud string
	flag.StringVar(&cloud, "c","aws", "set cloud,such aws、ali")
	flag.Parse()

	config.Setup()
	var ec2Client *ec2.EC2
	if cloud == "aws" {
		ec2Client = Aws.NewClient(config.AwsSetting.AccessKey, config.AwsSetting.SecretKey)
	}
	if cloud == "ali" {
		return
	}
	c, _ := consul.NewConsulServiceRegistry(config.ConsulSetting.Address, config.ConsulSetting.Port)

	instances, err := Aws.GetRunningInstances(ec2Client)
	if err != nil {
		fmt.Printf("Couldn't retrieve running instances: %v", err)
		return
	}

	for _, i := range instances {
		c.Register(i)
	}
	currmaps, _ := c.Lister()
	newInstances := consul.DiffRegister(instances, currmaps)

	for _, instance := range newInstances {
		c.Register(instance)
	}

}
