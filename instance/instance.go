package instance

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"math/rand"
	"time"
)

var (
	LinuxAmis = map[string]string{
		"us-east-1":      "ami-1ecae776",
		"us-west-2":      "ami-e7527ed7",
		"us-west-1":      "ami-d114f295",
		"eu-west-1":      "ami-a10897d6",
		"eu-central-1":   "ami-a8221fb5",
		"ap-northeast-1": "ami-cbf90ecb",
		"ap-southeast-1": "ami-68d8e93a",
		"ap-southeast-2": "ami-fd9cecc7",
		"sa-east-1":      "ami-b52890a8",
		"cn-north-1":     "ami-f239abcb",
	}
)

type Instance struct {
	svc *ec2.EC2
}

func New(region string) *Instance {
	return &Instance{ec2.New(&aws.Config{Region: region})}
}

func GetRandomRegion(regions []string) string {
	rand.Seed(time.Now().UTC().UnixNano())
	return regions[rand.Intn(len(regions))]
}

func GetLinuxAMIforRegion(region string) string {
	return LinuxAmis[region]
}
