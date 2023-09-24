package migrations

import (
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20230923_082642-User",
		Up: schema.Create("users", func(table *schema.Blueprint) {
			table.Int("id").Primary().AutoIncrement()
			table.String("email")
			table.String("access_token")
			table.DateTime("expires_in")
			table.String("refresh_token")
			table.Bool("finished_initial_fetch")
		}),
		Down: schema.DropIfExists("users"),
	})
}
