package health

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	cfg "github.com/artrctx/noliteo-core/internal/config"
	"github.com/artrctx/noliteo-core/internal/database"
	"github.com/artrctx/noliteo-core/tests/dbtest"
)

var testDBConfig *cfg.DatabaseConfig = &cfg.DatabaseConfig{}

func TestMain(m *testing.M) {
	dbtest.RunWithPostgres(m, testDBConfig)
}

func TestHealthHandler(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(HealthHandlerFunc(database.New(testDBConfig.ConnStr()))))
	defer srv.Close()
	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("error making request to server. err: %v", err)
	}
	defer resp.Body.Close()

	// assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected Status OK; got %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading reasponse body. Err: %v", err)
	}

	var res map[string]interface{}
	json.Unmarshal(body, &res)

	if res["status"] != "healthy" {
		t.Errorf("expected response body status to be healthy; got %v", res["status"])
	}
}
