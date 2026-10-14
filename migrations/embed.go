package migrations

// embedded SQL migration files, so the API binary is self-contained on
// platforms (e.g. Vercel) that do not ship the source tree alongside it.
import "embed"

//go:embed *.sql
var FS embed.FS
