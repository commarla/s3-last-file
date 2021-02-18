package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()

	session := session.Must(session.NewSession())

	svc := s3.New(session)

	bucket := viper.GetString("bucket")

	//nolint: exhaustivestruct
	objectsInput := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(1000), //nolint: gomnd
	}

	lastTime := time.Time{}

	err := svc.ListObjectsV2Pages(objectsInput, func(page *s3.ListObjectsV2Output, last bool) bool {
		time := processPages(page)

		diff := time.Sub(lastTime)
		if diff > 0 {
			lastTime = time
		}

		return !last
	})
	if err != nil {
		log.Fatal().Err(err).Msgf("unable to list objects")
	}

	log.Info().Msgf("lastest file is %v on %s", lastTime, bucket)

	loc, _ := time.LoadLocation("UTC")

	oneHourAgo := time.Now().In(loc).Add(-1 * time.Hour)

	diff := lastTime.Sub(oneHourAgo)
	if diff < 0 {
		log.Fatal().Msgf("lastest file is older than 1 hour %v on bucket %s", diff, bucket)
	}
}

func processPages(page *s3.ListObjectsV2Output) time.Time {
	lastTime := time.Time{}

	for _, object := range page.Contents {
		diff := object.LastModified.Sub(lastTime)
		if diff > 0 {
			lastTime = *object.LastModified
		}
	}

	return lastTime
}
