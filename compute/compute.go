package compute

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	LocalTerminationAddress = "http://169.254.169.254/latest/meta-data/spot/termination-time"
	LocalInstanceIDAddress  = "http://169.254.169.254/latest/meta-data/instance-id"
	LocalRegionAddress      = "http://169.254.169.254/latest/meta-data/placement/availability-zone"
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

type Compute struct {
	svc *ec2.EC2
}

func New(region string) *Compute {
	return &Compute{ec2.New(&aws.Config{Region: region})}
}

func GetRandomRegion(regions []string) string {
	rand.Seed(time.Now().UTC().UnixNano())
	return regions[rand.Intn(len(regions))]
}

func GetLinuxAMIforRegion(region string) string {
	return LinuxAmis[region]
}

func GetCurrentInstanceID() (id string) {
	res, err := http.Get(LocalInstanceIDAddress)
	if err != nil {
		return id
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return id
	}
	id = string(body)

	log.Println("Local instance ID is:" + id)

	return id
}

func GetCurrentRegion() (region string) {
	res, err := http.Get(LocalRegionAddress)
	if err != nil {
		return region
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return region
	}
	//trim the last character
	ln := len(string(body))
	region = string(body[:ln-1])

	return region
}

// return true if AWS has scheduled the instance to be terminated
func InstanceToBeTerminated(terminationUrl string) bool {

	if len(terminationUrl) <= 0 {
		terminationUrl = LocalTerminationAddress
	}

	res, err := http.Get(terminationUrl)
	if err != nil {
		return false
	}
	if res.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

func (c *Compute) TerminateInstance(instanceID string) (status string, err error) {

	log.Println("Terminating instance:" + instanceID)

	// status Codes and Names
	//  0 : pending
	//
	// 16 : running
	//
	// 32 : shutting-down
	//
	// 48 : terminated
	//
	// 64 : stopping
	//
	// 80 : stopped

	status = ""
	params := &ec2.TerminateInstancesInput{
		InstanceIDs: []*string{aws.String(instanceID)},
	}

	resp, err := c.svc.TerminateInstances(params)

	if err != nil {
		return status, err
	}

	status = *resp.TerminatingInstances[0].CurrentState.Name

	return status, nil
}
