// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/marketplacemetering"
	"github.com/aws/aws-sdk-go/service/marketplacemetering/marketplacemeteringiface"
)

type AWSMeterReport struct {
	Dimension string    `json:"dimension"`
	Value     int64     `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func (o *AWSMeterReport) ToJSON() string {
	b, _ := json.Marshal(o)
	return string(b)
}

type AWSMeterService struct {
	AwsDryRun      bool
	AwsProductCode string
	AwsMeteringSvc marketplacemeteringiface.MarketplaceMeteringAPI
}

func (a *App) newAWSMeterService() (*AWSMeterService, error) {
	svc := &AWSMeterService{
		AwsDryRun:      false,
		AwsProductCode: "12345", //TODO
	}

	service, err := newAWSMarketplaceMeteringService()

	if err != nil {
		mlog.Error("newAWSMeterService", mlog.String("error", err.Error()))
		return nil, err
	}

	svc.AwsMeteringSvc = service
	return svc, nil
}

func newAWSMarketplaceMeteringService() (*marketplacemetering.MarketplaceMetering, error) {
	region := os.Getenv("AWS_REGION")
	s := session.Must(session.NewSession(&aws.Config{Region: &region}))

	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(s),
			},
		})

	_, err := creds.Get()
	if err != nil {
		mlog.Error("session is invalid", mlog.String("error", err.Error()))
		return nil, errors.New("cannot obtain credentials")
	}

	return marketplacemetering.New(session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
	}))), nil
}

// a report entry is for all metrics
func (a *App) getUserCategoryUsage(dimensions []string, startTime time.Time, endTime time.Time) []*AWSMeterReport {
	reports := make([]*AWSMeterReport, 0)

	for _, dimension := range dimensions {
		var userCount int64
		var appErr *model.AppError

		switch dimension {
		case model.AWS_METERING_DIMENSION_USAGE_HRS:
			userCount, appErr = a.Srv().Store.User().AnalyticsActiveCountForPeriod(model.GetMillisForTime(startTime), model.GetMillisForTime(endTime), model.UserCountOptions{})

			if appErr != nil {
				mlog.Error("Failed to obtain usage data", mlog.String("dimension", dimension), mlog.String("start", startTime.String()), mlog.Int64("count", userCount))
			}
		default:
			mlog.Error("Dimension does not exist!", mlog.String("dimension", dimension))
		}

		if appErr != nil {
			mlog.Error("Failed to obtain usage.", mlog.String("dimension", dimension))
			return reports
		}

		report := &AWSMeterReport{
			Dimension: dimension,
			Value:     userCount,
			Timestamp: startTime,
		}

		reports = append(reports, report)
	}

	return reports
}

func (a *App) reportUserCategoryUsage(ams *AWSMeterService, reports []*AWSMeterReport) *model.AppError {
	for _, report := range reports {
		err := sendReportToMeteringService(ams, report)
		if err != nil {
			return err
		}
	}
	return nil
}

func sendReportToMeteringService(ams *AWSMeterService, report *AWSMeterReport) *model.AppError {
	params := &marketplacemetering.MeterUsageInput{
		DryRun:         aws.Bool(ams.AwsDryRun),
		ProductCode:    aws.String(ams.AwsProductCode),
		UsageDimension: aws.String(report.Dimension),
		UsageQuantity:  aws.Int64(report.Value),
		Timestamp:      aws.Time(report.Timestamp),
	}

	resp, err := ams.AwsMeteringSvc.MeterUsage(params)
	if err != nil {
		return model.NewAppError("sendReportToMeteringService", "app.system.aws_metering_service.error", nil, err.Error(), http.StatusNotFound)
	}
	if resp.MeteringRecordId == nil {
		return model.NewAppError("sendReportToMeteringService", "app.system.aws_metering_service.error", nil, "", http.StatusNotFound)
	}

	mlog.Debug("Sent record to AWS metering service", mlog.String("dimension", report.Dimension), mlog.Int64("value", report.Value), mlog.String("timestamp", report.Timestamp.String()))

	return nil
}
