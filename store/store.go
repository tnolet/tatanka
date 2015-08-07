package store

import (
	"bytes"
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
	svc         *s3.S3
}

type State struct {
	LifeTime           int
	LastLifeTimeTarget int
	LastLifeTimeActual int
	BidWindow          int
	BidOffset          int
	BidPrice           float32
	BidRatio           int
	InstanceSize       string
	InstanceType       string
	Regions            []string
	CurrentRegion      string
	CurrentInstanceID  string
	LastRegion         string
	PriceListUrl       string
	TerminationUrl     string
	TwitterHandle      string
}

func New(region string, bucket string, key string, localConfig string) *Store {
	svc := s3.New(&aws.Config{Region: region})
	return &Store{region, bucket, key, localConfig, svc}
}

func (s *Store) GetState() (state State, err error) {

	var configDestination = "state.json"

	// use bucket when there is no override
	if len(s.localConfig) <= 0 {

		log.Println("Calling S3 Bucket for state...")

		result, err := s.svc.GetObject(&s3.GetObjectInput{
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

func (s *Store) PutState(state State) error {

	log.Println("Storing state in S3...")

	stateJson, _ := json.Marshal(state)

	params := &s3.PutObjectInput{
		Bucket:      aws.String(s.homeBucket),
		Key:         aws.String(s.homeKey),
		Body:        bytes.NewReader([]byte(stateJson)),
		ContentType: aws.String("application/json"),
	}

	_, err := s.svc.PutObject(params)

	if err != nil {
		return err
	}

	return nil
}
