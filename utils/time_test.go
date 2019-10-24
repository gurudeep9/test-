// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var format = "2006-01-02 15:04:05.000000000"

func TestMillisFromTime(t *testing.T) {
	input, _ := time.Parse(format, "2015-01-01 12:34:00.000000000")
	actual := MillisFromTime(input)
	expected := int64(1420115640000)

	assert.Equal(t, actual, expected)
}

func TestYesterday(t *testing.T) {
	actual := Yesterday()
	expected := time.Now().AddDate(0, 0, -1)

	assert.Equal(t, actual.Year(), expected.Year())
	assert.Equal(t, actual.Day(), expected.Day())
	assert.Equal(t, actual.Month(), expected.Month())
}

func TestStartOfDay(t *testing.T) {
	input, _ := time.Parse(format, "2015-01-01 12:34:00.000000000")
	actual := StartOfDay(input)
	expected, _ := time.Parse(format, "2015-01-01 00:00:00.000000000")

	assert.Equal(t, actual, expected)
}

func TestEndOfDay(t *testing.T) {
	input, _ := time.Parse(format, "2015-01-01 12:34:00.000000000")
	actual := EndOfDay(input)
	expected, _ := time.Parse(format, "2015-01-01 23:59:59.999999999")

	assert.Equal(t, actual, expected)
}
