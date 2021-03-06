package main

import (
	"SuperQueue/logger"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
)

var (
	DBSession *gocqlx.Session

	// ---------------------------------------------------------------------------
	// DB Models
	// ---------------------------------------------------------------------------

	ItemsMetadata = &table.Metadata{
		Name: "items",
		Columns: []string{
			"namespace",
			"id",
			"payload",
			"bucket",
			"created_at",
			"expire_at",
			"in_flight_timeout",
		},
		PartKey: []string{"namespace", "id"},
	}

	ItemStatesMetadata = &table.Metadata{
		Name: "item_states",
		Columns: []string{
			"namespace",
			"id",
			"version",
			"state",
			"created_at",
			"attempts",
			"delay_to",
			"error",
			"error_message",
		},
		PartKey: []string{"namespace", "id"},
		SortKey: []string{"version"},
	}

	// ---------------------------------------------------------------------------
	// DB Tables
	// ---------------------------------------------------------------------------

	ItemsTable      = table.New(*ItemsMetadata)
	ItemStatesTable = table.New(*ItemStatesMetadata)
)

func DBConnect() {
	var cluster *gocql.ClusterConfig
	if SCYLLA_HOSTS == "" {
		logger.Debug("No SCYLLA_HOSTS env var found, using localhost")
		cluster = gocql.NewCluster("localhost:9042")
	} else {
		cluster = gocql.NewCluster(strings.Split(SCYLLA_HOSTS, ",")...)
	}
	cluster.Keyspace = "sq_ksp"
	cluster.Consistency = gocql.Quorum

	// Increase timeout if testing
	if os.Getenv("TEST_MODE") == "true" {
		cluster.Timeout = 1 * time.Second
	}

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal(err)
	}

	DBSession = &session
}

func DBConnectWithoutKeyspace() {
	var cluster *gocql.ClusterConfig
	if SCYLLA_HOSTS == "" {
		logger.Debug("No SCYLLA_HOSTS env var found, using localhost")
		cluster = gocql.NewCluster("localhost:9042")
	} else {
		cluster = gocql.NewCluster(SCYLLA_HOSTS)
	}
	cluster.Consistency = gocql.All

	// Increase timeout if testing
	if os.Getenv("TEST_MODE") == "true" {
		cluster.Timeout = 1 * time.Second
	}

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal(err)
	}

	DBSession = &session
}

func DBKeyspaceSetup() {
	// Create NS
	err := DBSession.ExecStmt("CREATE KEYSPACE IF NOT EXISTS sq_ksp WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}")
	if err != nil {
		log.Fatal(err)
	}
}

func DBReset() {
	err := DBSession.ExecStmt("DROP KEYSPACE IF EXISTS sq_ksp;")
	if err != nil {
		log.Fatal(err)
	}
}

func DBTableSetup() {
	err := DBSession.ExecStmt(`
		CREATE TABLE IF NOT EXISTS items (
			namespace TEXT,
			id TEXT,
			payload BLOB,
			bucket TEXT,
			created_at TIMESTAMP,
			expire_at TIMESTAMP,
			in_flight_timeout INT,
			PRIMARY KEY((namespace, id))
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	err = DBSession.ExecStmt(`
		CREATE TABLE IF NOT EXISTS item_states (
			namespace TEXT,
			id TEXT,
			version INT,
			state TEXT,
			created_at TIMESTAMP,
			attempts INT,
			delay_to TIMESTAMP,
			error TEXT,
			error_message TEXT,
			PRIMARY KEY((namespace, id), version)
		) WITH CLUSTERING ORDER BY (version DESC);
	`)
	if err != nil {
		log.Fatal(err)
	}
}
