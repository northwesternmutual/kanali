package monitor

import (
	"errors"
	"regexp"
	"testing"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	db string
}

func (c *mockClient) Ping(timeout time.Duration) (time.Duration, string, error) {
	return 123456789, "", nil
}

func (c *mockClient) Write(bp influx.BatchPoints) error {
	if bp.Database() == "" || c.db != bp.Database() {
		return errors.New("database does not exist")
	}
	return nil
}

func (c *mockClient) Query(q influx.Query) (*influx.Response, error) {
	re := regexp.MustCompile("^CREATE DATABASE (.*)")
	match := re.FindStringSubmatch(q.Command)
	if len(match) != 2 {
		return nil, errors.New("query incorrect")
	}
	if match[1] == "" {
		return nil, errors.New("no database name")
	}
	c.db = match[1]
	return nil, nil
}

func (c *mockClient) Close() error {
	return nil
}

func TestWriteRequestData(t *testing.T) {
	ctlr := &InfluxController{Client: &mockClient{}}
	m := &metrics.Metrics{
		metrics.Metric{"metric-one", "value-one", true},
		metrics.Metric{"metric-two", "value-two", false},
	}
	assert.Equal(t, ctlr.WriteRequestData(m).Error(), "no database name")
	viper.SetDefault(config.FlagInfluxdbDatabase.GetLong(), "mydb")
	assert.Nil(t, ctlr.WriteRequestData(m))
	assert.Nil(t, ctlr.WriteRequestData(m))
	viper.SetDefault(config.FlagInfluxdbDatabase.GetLong(), "")
	ctlr = &InfluxController{Client: nil}
	assert.Equal(t, ctlr.WriteRequestData(m).Error(), "influxdb paniced while attempting to write")
}

func TestNewInfluxdbController(t *testing.T) {
	_, err := NewInfluxdbController()
	assert.NotNil(t, err)
	viper.SetDefault(config.FlagInfluxdbAddr.GetLong(), "http://foo.bar.com")
	_, err = NewInfluxdbController()
	assert.Nil(t, err)
}

func TestCreateDatabase(t *testing.T) {
	err := createDatabase(&mockClient{})
	assert.Equal(t, err.Error(), "no database name")
	viper.SetDefault(config.FlagInfluxdbDatabase.GetLong(), "mydb")
	assert.Nil(t, createDatabase(&mockClient{}))
	viper.SetDefault(config.FlagInfluxdbDatabase.GetLong(), "")
}

func TestGetFields(t *testing.T) {
	assert.Equal(t, getFields(&metrics.Metrics{
		metrics.Metric{"metric-one", "value-one", true},
		metrics.Metric{"metric-two", "value-two", false},
	}), map[string]interface{}{
		"metric-one": "value-one",
		"metric-two": "value-two",
	})
}

func TestGetTags(t *testing.T) {
	tags, err := getTags(&metrics.Metrics{
		metrics.Metric{"metric-one", "value-one", true},
		metrics.Metric{"metric-two", "value-two", false},
	})
	assert.Nil(t, err)
	assert.Equal(t, tags, map[string]string{
		"metric-one": "value-one",
	})
	tags, err = getTags(&metrics.Metrics{
		metrics.Metric{"metric-one", 5, true},
		metrics.Metric{"metric-two", "value-two", false},
	})
	assert.Nil(t, tags)
	assert.Equal(t, err.Error(), "InfluxDB requires that the indexed field metric-one be of type string")
}
