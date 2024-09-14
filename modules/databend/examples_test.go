package databend_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/databend"
)

func ExampleRun() {
	ctx := context.Background()

	databendContainer, err := databend.Run(ctx,
		"datafuselabs/databend:v1.2.615",
		databend.WithUsername("databend"),
		databend.WithPassword("databend"),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(databendContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	state, err := databendContainer.State(ctx)
	if err != nil {
		log.Printf("failed to get container state: %s", err)
		return
	}

	fmt.Println(state.Running)

	// Output:
	// true
}

func ExampleRun_connect() {
	ctx := context.Background()

	databendContainer, err := databend.Run(ctx,
		"datafuselabs/databend:v1.2.615",
		databend.WithUsername("root"),
		databend.WithPassword("password"),
		databend.WithDatabase("test"),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(databendContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	connectionString, err := databendContainer.ConnectionString(ctx)
	if err != nil {
		log.Printf("failed to get connection string: %s", err)
		return
	}

	db, err := sql.Open("databend", connectionString)
	if err != nil {
		log.Printf("failed to connect to Databend: %s", err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Printf("failed to ping Databend: %s", err)
		return
	}
	var i int
	row, err := db.Query("select 1")
	if err != nil {
		log.Printf("failed to query Databend: %s", err)
		return
	}
	err = row.Scan(&i)
	if err != nil {
		log.Printf("failed to scan result: %s", err)
		return
	}

	fmt.Println(i)

	// Output:
	// true
}
