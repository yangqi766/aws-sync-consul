package aws

import (
	"asset/sync/model"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func NewClient(ak, sk string) *ec2.EC2 {
	sess, _ := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentials(ak, sk, ""),
			Region: aws.String("ap-southeast-1"),
		},
	})

	return ec2.New(sess)
}

func GetRunningInstances(client *ec2.EC2) ([]model.Instance, error) {
	var Instances []model.Instance
	result, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			i := model.Instance{
				InstanceId:      *instance.InstanceId,
				InstanceAddress: *instance.PrivateIpAddress,
				InstanceName:    *instance.Tags[0].Value,
			}
			Instances = append(Instances, i)
		}
	}
	return Instances, err
}
