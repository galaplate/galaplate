package migrations

import (
	"github.com/galaplate/core/database"
)

type Migration1765558170 struct {
	database.BaseMigration
}

func init() {
	migration := &Migration1765558170{
		BaseMigration: database.BaseMigration{
			Name:      "create_users_table",
			Timestamp: 1765558170,
		},
	}
	database.Register(migration)
}

func (m *Migration1765558170) Up(schema *database.Schema) error {
	return schema.Create("users", func(table *database.Blueprint) {
		table.ID()
		table.String("username", 50).NotNullable()
		table.String("email", 100).Unique().NotNullable()
		table.String("password").NotNullable()
		table.String("description")
		table.Boolean("status").Default(false)
		table.Timestamps()
		table.DateTime("deleted_at").Nullable()
	})
}

func (m *Migration1765558170) Down(schema *database.Schema) error {
	return schema.DropIfExists("users")
}
