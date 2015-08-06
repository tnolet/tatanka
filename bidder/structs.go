package bidder

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Bidder struct {
	priceListUrl string
	svc          *ec2.EC2
	region       string
}

type SpotRequest struct {
	Id string
}

type PriceList struct {
	Vers   float64 `json:"vers"`
	Config struct {
		Rate         string   `json:"rate"`
		Valuecolumns []string `json:"valueColumns"`
		Currencies   []string `json:"currencies"`
		Regions      []struct {
			Region        string `json:"region"`
			Instancetypes []struct {
				Type  string `json:"type"`
				Sizes []struct {
					Size         string `json:"size"`
					Vcpu         string `json:"vCPU"`
					Ecu          string `json:"ECU"`
					Memorygib    string `json:"memoryGiB"`
					Storagegb    string `json:"storageGB"`
					Valuecolumns []struct {
						Name   string `json:"name"`
						Prices struct {
							Usd string `json:"USD"`
						} `json:"prices"`
					} `json:"valueColumns"`
				} `json:"sizes"`
			} `json:"instanceTypes"`
		} `json:"regions"`
	} `json:"config"`
}
