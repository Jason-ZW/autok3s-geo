package db

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/mmcloughlin/geohash"
	"os"
	"strconv"
	"time"

	"github.com/Jason-ZW/autok3s-geo/pkg/types"

	"github.com/influxdata/influxdb-client-go/v2"
)

const (
	ENDPOINT    = "ENDPOINT"
	TOKEN       = "TOKEN"
	ORG         = "autok3s"
	BUCKET      = "autok3s"
	MEASUREMENT = "geo_locations"
	REPOSITORY  = "autok3s"
	START       = "2022-01-01T00:00:00Z"
)

type DB struct {
	ctx    context.Context
	client influxdb2.Client
}

func NewDB(ctx context.Context) (*DB, error) {
	endpoint, token := GetFromEnv()
	if endpoint == "" || token == "" {
		return nil, errors.New(fmt.Sprintf("environment %s or %s is empty, please check that", ENDPOINT, TOKEN))
	}

	client := influxdb2.NewClientWithOptions(endpoint, token, influxdb2.DefaultOptions().
		SetUseGZip(true).
		SetTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		}))

	isAvailable, err := client.Ping(ctx)
	if err != nil {
		return nil, err
	}

	if !isAvailable {
		return nil, errors.New(fmt.Sprintf("remote server %s is not available", ENDPOINT))
	}

	db := &DB{
		ctx:    ctx,
		client: client,
	}

	return db, nil
}

func (d *DB) Close() {
	if d.client != nil {
		d.client.Close()
	}
}

func (d *DB) QueryByIPIn24Hours(ip string) ([]types.GeoIP, error) {
	geoIPs := make([]types.GeoIP, 0)

	if d.client == nil {
		return geoIPs, errors.New("failed to query metrics, because client is nil")
	}

	query := d.client.QueryAPI(ORG)

	result, err := query.Query(d.ctx,
		fmt.Sprintf(`from(bucket: "%s") |> range(start: %d, stop: now()) |> filter(fn: (r) => r._measurement == "%s" and r.repository == "%s" and r.ip == "%s")`,
			BUCKET, time.Now().Add(-24*time.Hour).UTC().Unix(), MEASUREMENT, REPOSITORY, ip))
	if err != nil {
		return geoIPs, err
	}

	for result.Next() {
		geoIP := types.GeoIP{
			IP:          result.Record().ValueByKey("ip").(string),
			CountryName: result.Record().ValueByKey("country").(string),
			City:        result.Record().ValueByKey("city").(string),
			Active:      result.Record().ValueByKey("_value").(int64),
		}

		if val, ok := result.Record().ValueByKey("city").(string); ok {
			geoIP.City = val
		} else {
			geoIP.City = geohash.Encode(geoIP.Latitude, geoIP.Longitude)
		}

		if val, ok := result.Record().ValueByKey("country").(string); ok {
			geoIP.CountryName = val
		}

		geoIP.Latitude, _ = strconv.ParseFloat(result.Record().ValueByKey("latitude").(string), 32)
		geoIP.Longitude, _ = strconv.ParseFloat(result.Record().ValueByKey("longitude").(string), 32)
		geoIPs = append(geoIPs, geoIP)
	}
	if result.Err() != nil {
		return geoIPs, result.Err()
	}

	return geoIPs, nil
}

func (d *DB) Write(geoIP *types.GeoIP) (*types.GeoIP, error) {
	if d.client == nil {
		return nil, errors.New("failed to write metrics, because client is nil")
	}

	write := d.client.WriteAPIBlocking(ORG, BUCKET)

	if geoIP.City == "" {
		geoIP.City = geohash.Encode(geoIP.Latitude, geoIP.Longitude)
	}

	if geoIP.City == "" {
		geoIP.City = geohash.Encode(geoIP.Latitude, geoIP.Longitude)
	}

	point := influxdb2.NewPointWithMeasurement(MEASUREMENT).
		AddTag("repository", REPOSITORY).
		AddTag("ip", geoIP.IP).
		AddTag("country", geoIP.CountryName).
		AddTag("city", geoIP.City).
		AddTag("latitude", fmt.Sprintf("%f", geoIP.Latitude)).
		AddTag("longitude", fmt.Sprintf("%f", geoIP.Longitude)).
		AddField("active", 1).
		SetTime(time.Now())

	err := write.WritePoint(d.ctx, point)
	if err != nil {
		return nil, err
	}

	return geoIP, nil
}

func GetFromEnv() (string, string) {
	return os.Getenv(ENDPOINT), os.Getenv(TOKEN)
}
