package bidder

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	UserData = `#!/bin/bash
set -e -x
echo "Getting tatanka..."
curl -OL http://tnolet-tatanka.s3-eu-west-1.amazonaws.com/tatanka
chmod +x tatanka
./tatanka
`
	tagKey   = "bidder"
	tagValue = "tatanka"
)

var (
	priceListRegexp = regexp.MustCompile("(callback\\()(.*)\\);$")
)

func New(priceList string, region string) *Bidder {
	log.Println("Initializing bidder for region:", region)
	return &Bidder{priceList, ec2.New(&aws.Config{Region: &region}), region}
}

func (b *Bidder) CreateSpotRequest(price string, size string, amiId string, from time.Time, till time.Time) (requests []*SpotRequest, err error) {

	// var keyName = "aws_dcos"
	var amount int64 = 1
	var reqType = "one-time"
	var IAMRoleARN = "arn:aws:iam::539701811563:instance-profile/tatanka"
	var base64UserData = base64.StdEncoding.EncodeToString([]byte(UserData))

	// create request
	params := &ec2.RequestSpotInstancesInput{
		SpotPrice:     aws.String(price),
		InstanceCount: &amount,
		ValidFrom:     &from,
		ValidUntil:    &till,
		LaunchSpecification: &ec2.RequestSpotLaunchSpecification{
			InstanceType: aws.String(size),
			ImageId:      aws.String(amiId),
			// KeyName:      aws.String(keyName),
			UserData: aws.String(base64UserData),
			IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
				Arn: &IAMRoleARN,
			},
		},
		Type: aws.String(reqType),
	}

	resp, err := b.svc.RequestSpotInstances(params)

	if err != nil {
		return requests, errors.New("Error creating spot request: " + err.Error())
	}

	// compile list
	for _, req := range resp.SpotInstanceRequests {
		requests = append(requests, &SpotRequest{Id: *req.SpotInstanceRequestId})
	}

	// tag them
	for _, req := range requests {
		params := &ec2.CreateTagsInput{
			Resources: []*string{aws.String(req.Id)},
			Tags:      []*ec2.Tag{{Key: aws.String("bidder"), Value: aws.String("tatanka")}},
		}
		_, err := b.svc.CreateTags(params)
		if err != nil {
			return requests, errors.New("Error creating spot request: " + err.Error())
		}
		log.Printf("Created spot request: %v", req.Id)
	}

	return requests, nil
}

func (b *Bidder) SpotInstanceActive(reqId string) bool {

	reqs, err := b.GetSpotInstanceRequests("fulfilled")

	if err != nil {
		return false
	}

	for _, req := range reqs {
		if req.Id == reqId {
			return true
		}
	}
	return false
}

func (b *Bidder) CancelSpotRequests() error {

	log.Println("Cancelling outstanding spot requests")

	reqs, err := b.GetSpotInstanceRequests("open")

	if err != nil {
		return err
	}

	var ids []*string
	for _, req := range reqs {
		ids = append(ids, &req.Id)
		log.Println("Cancelling spot request:" + req.Id)
	}

	params := &ec2.CancelSpotInstanceRequestsInput{
		SpotInstanceRequestIds: ids,
	}
	_, err = b.svc.CancelSpotInstanceRequests(params)

	if err != nil {
		return err
	}
	return nil
}

func (b *Bidder) GetSpotInstanceRequests(reqState string) (requests []*SpotRequest, err error) {

	reqType := "one-time"

	params := &ec2.DescribeSpotInstanceRequestsInput{
		Filters: filters(reqState, tagKey, tagValue, reqType),
	}
	resp, err := b.svc.DescribeSpotInstanceRequests(params)

	if err != nil {
		return requests, err
	}

	for _, req := range resp.SpotInstanceRequests {
		requests = append(requests, &SpotRequest{Id: *req.SpotInstanceRequestId})
	}

	return requests, nil
}

func (b *Bidder) GetSpotPriceHistory(size string) (price string, err error) {

	var now = time.Now()
	var startTime = now.AddDate(0, -1, 0) // yesterday
	var maxResults int64 = 30
	var description = "Linux/UNIX"

	params := &ec2.DescribeSpotPriceHistoryInput{
		StartTime: aws.Time(startTime),
		EndTime:   aws.Time(now),
		InstanceTypes: []*string{
			aws.String(size),
		},
		MaxResults: &maxResults,
		ProductDescriptions: []*string{
			aws.String(description), // Required
			// More values...
		},
	}

	resp, err := b.svc.DescribeSpotPriceHistory(params)

	if err != nil {
		return price, err
	}

	price = *resp.SpotPriceHistory[0].SpotPrice
	return price, nil

}

func (b *Bidder) GetOnDemandPrice(genType string, size string) (price string, err error) {

	list, err := b.getList()
	if err != nil {
		return price, err
	}

	for _, r := range list.Config.Regions {
		if r.Region == b.region {
			for _, g := range r.Instancetypes {
				if g.Type == genType {
					for _, s := range g.Sizes {
						if s.Size == size {
							price = s.Valuecolumns[0].Prices.Usd
						}
					}
				}
			}
		}
	}

	return price, nil
}

func (b *Bidder) getList() (p PriceList, err error) {

	// grab the list page
	res, err := http.Get(b.priceListUrl)
	if err != nil {
		return p, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return p, err
	}

	listJson := parsePageToJson(body)

	if err := json.Unmarshal([]byte(listJson), &p); err != nil {
		return p, err
	}

	return p, nil

}

func parsePageToJson(page []byte) (listJson string) {

	// parse out the javascript object
	listJS := priceListRegexp.FindStringSubmatch(string(page[:]))[2]

	// replace the JS keys with quoted strings
	replaceable := []string{"vers", "rate", "valueColumns", "currencies", "region", "regions",
		"type", "size", "sizes", "vCPU", "ECU", "memoryGiB", "storageGB", "name", "USD", "prices", "instanceTypes", "config"}

	for _, item := range replaceable {
		toReplace := item + ":"
		replaceWith := "\"" + item + "\":"
		listJS = strings.Replace(listJS, toReplace, replaceWith, -1)
	}
	return listJS
}

func filters(reqState string, tagKey string, tagValue string, reqType string) []*ec2.Filter {

	return []*ec2.Filter{
		{
			Name: aws.String("state"),
			Values: []*string{
				aws.String(reqState),
			},
		},
		{
			Name: aws.String("tag-key"),
			Values: []*string{
				aws.String(tagKey),
			},
		},
		{
			Name: aws.String("tag-value"),
			Values: []*string{
				aws.String(tagValue),
			},
		},
		{
			Name: aws.String("type"),
			Values: []*string{
				aws.String(reqType),
			},
		},
	}
}
