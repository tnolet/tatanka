package instance

import (
	"testing"
)

var (
	Regions = []string{
		"us-east-1",
		"us-west-2",
		"us-west-1",
		"eu-west-1",
		"eu-central-1",
		"ap-northeast-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"sa-east-1",
		"cn-north-1",
	}
)

func TestGetLinuxAMIforRegion(t *testing.T) {
	if image := GetLinuxAMIforRegion("eu-west-1"); image != "ami-a10897d6" {
		t.Errorf("Failed to get correct AMZ Linux AMI for region eu-west. Got %v", image)
	}
}

func TestGetRandomRegion(t *testing.T) {
	region := GetRandomRegion(Regions)
	found := false
	for _, reg := range Regions {
		if reg == region {
			found = true
		}
	}
	if found == false {
		t.Errorf("Failed to get a random region, got %v", region)
	}

}
