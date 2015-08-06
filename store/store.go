package store

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Store struct {
	homeRegion  string
	homeBucket  string
	homeKey     string
	localConfig string
}

type State struct {
	LifeTime           int
	LastLifeTimeTarget int
	LastLifeTimeActual int
	BidWindow          int
	InstanceSize       string
	InstanceType       string
	Regions            []string
	PriceListUrl       string
	TerminationUrl     string
	BidPrice           float32
	BidRatio           int
	TwitterHandle      string
}

func New(region string, bucket string, key string, localConfig string) *Store {
	return &Store{region, bucket, key, localConfig}
}

func (s *Store) GetState() (state State, err error) {

	var configDestination = "state.json"

	// use bucket when there is no override
	if len(s.localConfig) <= 0 {

		log.Println("Calling S3 Bucket for state...")

		aws.DefaultConfig.Region = s.homeRegion
		svc := s3.New(nil)
		result, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(s.homeBucket),
			Key:    aws.String(s.homeKey),
		})

		if err != nil {
			return state, err
		}

		file, err := os.Create(configDestination)
		if err != nil {
			return state, err
		}

		if _, err := io.Copy(file, result.Body); err != nil {
			return state, err
		}

		// download the full file
		result.Body.Close()
		file.Close()

	} else {
		log.Printf("Using local file %v for state...", s.localConfig)

		// use the local file directly
		configDestination = s.localConfig
	}

	if s, err := ioutil.ReadFile(configDestination); err != nil {
		return state, err
	} else {
		if err := json.Unmarshal(s, &state); err != nil {
			return state, err
		}
	}

	return state, nil
}
