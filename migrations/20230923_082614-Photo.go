package migrations

import (
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20230923_082614-Photo",
		Up: schema.Create("photos", func(table *schema.Blueprint) {
			table.Int("id").Primary().AutoIncrement()
			table.Int("user_id")
			table.String("photo_id").Unique()
		}),
		Down: schema.DropIfExists("photos"),
	})
}
