package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"

	"github.com/pkg/errors"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type Storage interface {
	Save(ctx context.Context, nl NomadLocation) error
	List(ctx context.Context, chatID string) ([]NomadLocation, error)
	ListByCountry(ctx context.Context, chatID, country string) ([]NomadLocation, error)
	ListByCity(ctx context.Context, chatID, city string) ([]NomadLocation, error)
}

const dbName = "nmad"
const colName = "nomad_locations"

type arangodbStorage struct {
	c driver.Client
}

func NewArangodbDBStorage() (Storage, error) {
	// Decode CA certificate
	caCertificate, err := base64.StdEncoding.DecodeString(CONFIG.ArangoDBCA)
	if err != nil {
		return nil, errors.WithMessagef(err, "DecodeString")
	}

	tlsConfig := &tls.Config{}
	certpool := x509.NewCertPool()
	if success := certpool.AppendCertsFromPEM(caCertificate); !success {
		return nil, errors.WithMessagef(err, "Invalid certificate")
	}
	tlsConfig.RootCAs = certpool

	endpoint := CONFIG.ArangoDBEndpoint
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{endpoint},
		TLSConfig: tlsConfig,
	})
	if err != nil {
		return nil, errors.WithMessagef(err, "NewConnection %s", endpoint)
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(CONFIG.ArangoDBUser, CONFIG.ArangoDBPassword),
	})
	if err != nil {
		return nil, errors.WithMessagef(err, "NewClient")
	}
	return &arangodbStorage{client}, nil
}

func (s *arangodbStorage) Save(ctx context.Context, nl NomadLocation) error {
	col, err := s.getCollection(ctx)
	if err != nil {
		return err
	}

	_, err = col.CreateDocument(ctx, nl)
	if err != nil {
		return errors.WithMessagef(err, "CreateDocument")
	}

	return nil
}

func (s *arangodbStorage) List(ctx context.Context, chatID string) ([]NomadLocation, error) {
	q := `
FOR doc IN @@col_name
  FILTER doc.chat_id == @chat_id
  SORT doc.at DESC, doc.city ASC
  COLLECT username = doc.username INTO users
  LET nl = FIRST(FOR u IN users[*].doc SORT u.at DESC LIMIT 1 RETURN u) 
RETURN nl
`
	vals := map[string]interface{}{
		"@col_name": colName,
		"chat_id":   chatID,
	}
	return s.queryNomads(ctx, q, vals)
}

func (s *arangodbStorage) ListByCity(ctx context.Context, chatID, city string) ([]NomadLocation, error) {
	q := `
FOR doc IN @@col_name
  FILTER doc.chat_id == @chat_id
  SORT doc.at DESC, doc.city ASC
  COLLECT username = doc.username INTO users
  LET nl = FIRST(FOR u IN users[*].doc SORT u.at DESC LIMIT 1 RETURN u) 
  FILTER nl.city == @city
RETURN nl
`
	vals := map[string]interface{}{
		"@col_name": colName,
		"chat_id":   chatID,
		"city":      city,
	}
	return s.queryNomads(ctx, q, vals)
}

func (s *arangodbStorage) ListByCountry(ctx context.Context, chatID, country string) ([]NomadLocation, error) {
	q := `
FOR doc IN @@col_name
  FILTER doc.chat_id == @chat_id
  SORT doc.at DESC, doc.city ASC
  COLLECT username = doc.username INTO users
  LET nl = FIRST(FOR u IN users[*].doc SORT u.at DESC LIMIT 1 RETURN u) 
  FILTER nl.country == @country
RETURN nl
`
	vals := map[string]interface{}{
		"@col_name": colName,
		"chat_id":   chatID,
		"country":   country,
	}
	return s.queryNomads(ctx, q, vals)
}

func (s *arangodbStorage) queryNomads(ctx context.Context, q string, vals map[string]interface{}) ([]NomadLocation, error) {
	db, err := s.getDB(ctx)
	if err != nil {
		return nil, err
	}
	c, err := db.Query(ctx, q, vals)
	if err != nil {
		return nil, errors.WithMessagef(err, "Query %s", q)
	}
	defer c.Close()

	result := make([]NomadLocation, 0)
	for {
		var doc NomadLocation
		_, err = c.ReadDocument(nil, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, errors.WithMessagef(err, "While reading document")
		}

		result = append(result, doc)
	}

	return result, nil
}

func (s *arangodbStorage) getDB(ctx context.Context) (driver.Database, error) {
	db, err := s.c.Database(ctx, dbName)
	if err != nil {
		return nil, errors.WithMessagef(err, "Database failed")
	}
	return db, nil
}

func (s *arangodbStorage) getCollection(ctx context.Context) (driver.Collection, error) {
	db, err := s.getDB(ctx)

	col, err := db.Collection(ctx, colName)
	if err != nil {
		return nil, errors.WithMessagef(err, "Collection failed")
	}
	return col, nil
}
