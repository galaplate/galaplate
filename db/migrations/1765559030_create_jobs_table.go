package migrations

import (
	"github.com/galaplate/core/database"
)

type Migration1765559030 struct {
	database.BaseMigration
}

func init() {
	migration := &Migration1765559030{
		BaseMigration: database.BaseMigration{
			Name:      "create_jobs_table",
			Timestamp: 1765559030,
		},
	}
	database.Register(migration)
}

func (m *Migration1765559030) Up(schema *database.Schema) error {
	return schema.Create("jobs", func(table *database.Blueprint) {
		table.ID()
		table.String("type").NotNullable()
		table.Text("payload").NotNullable()
		table.Enum("state", []string{"pending", "started", "finished", "failed"}).NotNullable().Default("pending")
		table.Text("error_msg")
		table.Integer("attempts").Default(0)
		table.DateTime("available_at")
		table.DateTime("started_at")
		table.DateTime("finished_at")
		table.Timestamps()
	})
}

func (m *Migration1765559030) Down(schema *database.Schema) error {
	return schema.DropIfExists("jobs")
}
