package job

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://app:app@localhost:5432/app?sslmode=disable"
	}

	var err error
	testDB, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		panic("failed to connect to test DB: " + err.Error())
	}
	defer testDB.Close()

	os.Exit(m.Run())
}

func cleanJobs(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec(context.Background(), "DELETE FROM jobs")
	if err != nil {
		t.Fatalf("failed to clean jobs: %v", err)
	}
}

func TestStore_Create(t *testing.T) {
	cleanJobs(t)
	// TODO
}

func TestStore_Claim(t *testing.T) {
	cleanJobs(t)
	// TODO: insert a pending job, call Claim, assert status is 'processing'
}

func TestStore_Claim_SkipsProcessing(t *testing.T) {
	cleanJobs(t)
	// TODO: insert a job with status='processing', call Claim, assert returns nil (no job)
}

func TestStore_Complete(t *testing.T) {
	cleanJobs(t)
	// TODO: insert a job, call Complete, assert status is 'completed' and result is set
}

func TestStore_Fail_Retries(t *testing.T) {
	cleanJobs(t)
	// TODO: insert a job with retry_count=0, call Fail, assert status back to 'pending' and retry_after is set
}

func TestStore_Fail_MaxRetries(t *testing.T) {
	cleanJobs(t)
	// TODO: insert a job with retry_count=MaxRetry-1, call Fail, assert status is 'failed'
}
