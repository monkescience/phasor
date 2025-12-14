// Package fixtures provides shared test utilities for integration tests.
package fixtures

// Testcontainers support for external dependencies like databases.
//
// When you need to test with a real database, add the testcontainers dependency:
//   go get github.com/testcontainers/testcontainers-go
//   go get github.com/testcontainers/testcontainers-go/modules/postgres
//
// Example usage:
//
//	func StartPostgres(ctx context.Context, t *testing.T) *PostgresContainer {
//	    container, err := postgres.Run(ctx, "postgres:16-alpine",
//	        postgres.WithDatabase("testdb"),
//	        postgres.WithUsername("test"),
//	        postgres.WithPassword("test"),
//	    )
//	    if err != nil {
//	        t.Fatalf("failed to start postgres: %v", err)
//	    }
//
//	    dsn, err := container.ConnectionString(ctx, "sslmode=disable")
//	    if err != nil {
//	        t.Fatalf("failed to get connection string: %v", err)
//	    }
//
//	    t.Cleanup(func() {
//	        if err := container.Terminate(ctx); err != nil {
//	            t.Logf("failed to terminate postgres: %v", err)
//	        }
//	    })
//
//	    return &PostgresContainer{DSN: dsn}
//	}
//
// Then in your test:
//
//	func TestWithDatabase(t *testing.T) {
//	    pg := fixtures.StartPostgres(ctx, t)
//	    backend := backendserver.NewTestServer(
//	        backendserver.WithDatabaseDSN(pg.DSN),
//	    )
//	    // Test via HTTP...
//	}
