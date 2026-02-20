// Walkie Talkie Service
package walktalk

import (
	"database/sql"

	"github.com/artrctx/noliteo-core/internal/hub"
)

type WalkTalkService struct {
	DB  *sql.DB
	Hub *hub.Hub
}
