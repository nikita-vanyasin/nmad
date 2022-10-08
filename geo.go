package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type GeoInfo interface {
	LookupCity(ctx context.Context, n string) (*Location, error)
	LookupCountry(ctx context.Context, n string) (string, error)
}

const (
	featureClassCity    = "P"
	featureClassCountry = "A"
)

type geoNamesClient struct {
	loginName string
	// TODO: add cache
}

func NewGeoInfo() (GeoInfo, error) {
	return &geoNamesClient{
		loginName: CONFIG.GeoNamesAPILogin,
	}, nil
}

func (c *geoNamesClient) LookupCity(ctx context.Context, name string) (*Location, error) {
	return c.lookupName(ctx, name, featureClassCity)
}

func (c *geoNamesClient) LookupCountry(ctx context.Context, name string) (string, error) {
	l, err := c.lookupName(ctx, name, featureClassCountry)
	if err != nil {
		return "", err
	}

	if l == nil {
		return "", nil
	}

	return l.Country, nil
}

func (c *geoNamesClient) lookupName(ctx context.Context, name string, featureClass string) (*Location, error) {
	var u = url.URL{
		Scheme: "http",
		Host:   "api.geonames.org",
		Path:   "/searchJSON",
		RawQuery: url.Values{
			"username":     []string{c.loginName},
			"maxRows":      []string{"1"},
			"q":            []string{name},
			"featureClass": []string{featureClass},
		}.Encode(),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, errors.WithMessagef(err, "NewRequest")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithMessagef(err, "Do request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessagef(err, "ReadAll")
	}

	var response struct {
		GeoNames []struct {
			Name        string `json:"name,omitempty"`
			CountryName string `json:"countryName,omitempty"`
			Lat         string `json:"lat"`
			Lng         string `json:"lng"`
		} `json:"geonames"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.WithMessagef(err, "Unmarshall")
	}

	if len(response.GeoNames) == 0 {
		return nil, nil
	}

	r := response.GeoNames[0]
	lat, err := strconv.ParseFloat(r.Lat, 32)
	if err != nil {
		return nil, errors.WithMessagef(err, "%s", r.Lat)
	}
	lng, err := strconv.ParseFloat(r.Lng, 32)
	if err != nil {
		return nil, errors.WithMessagef(err, "%s", r.Lng)
	}

	return &Location{
		Country: r.CountryName,
		City:    r.Name,
		Lat:     lat,
		Lng:     lng,
	}, nil
}
