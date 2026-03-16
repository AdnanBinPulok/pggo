package main

import (
	"fmt"
	"pggo"
	"time"
)

// User is the demo model mapped to users table.
type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Age       int       `db:"age"`
	CreatedAt time.Time `db:"created_at"`
}

// RunStrictDemo demonstrates end-to-end strict API usage.
func RunStrictDemo(dbURL string) error {
	conn := pggo.NewDatabaseConnection(dbURL, 20, true)
	defer conn.Close()

	usersTable := pggo.Table[User]{
		Name:       "users_demo",
		Connection: conn,
		Columns: []pggo.Column{
			{Name: "id", DataType: *pggo.DataType.Serial().PrimaryKey()},
			{Name: "name", DataType: *pggo.DataType.Text().NotNull()},
			{Name: "email", DataType: *pggo.DataType.Text().Unique().NotNull()},
			{Name: "age", DataType: *pggo.DataType.Integer()},
			{Name: "created_at", DataType: *pggo.DataType.TimestampTz()},
		},
	}

	usersTable.SetSchemaSync(pggo.SyncOptions{
		SyncColumns:      true,
		DropMissing:      false,
		AllowDestructive: false,
		DryRun:           false,
	})
	usersTable.CacheKey = "id"
	usersTable.EnableCache(5 * time.Second)

	_ = usersTable.DeleteTable()
	if err := usersTable.CreateTable(); err != nil {
		return err
	}

	alice, err := usersTable.InsertOne(User{Name: "Alice", Email: "alice@example.com", Age: 25})
	if err != nil {
		return err
	}

	_, err = usersTable.InsertMany([]User{
		{Name: "Bob", Email: "bob@example.com", Age: 30},
		{Name: "Carol", Email: "carol@example.com", Age: 22},
	})
	if err != nil {
		return err
	}

	_, err = usersTable.FetchOne("id", alice.ID)
	if err != nil {
		return err
	}

	rows, err := usersTable.FetchMany(map[string]any{"age": pggo.Gte(20)})
	if err != nil {
		return err
	}

	count, err := usersTable.Count(map[string]any{"age": pggo.Gte(20)})
	if err != nil {
		return err
	}

	pages, err := usersTable.FetchPages(1, 2, "id", "asc", map[string]any{"age": pggo.Gte(20)})
	if err != nil {
		return err
	}

	searchPages, err := usersTable.PageSearch(1, 10, "id", "desc", map[string]any{}, []string{"name", "email"}, "ali")
	if err != nil {
		return err
	}

	searchCount, err := usersTable.SearchCount(map[string]any{}, []string{"name", "email"}, "example.com")
	if err != nil {
		return err
	}

	_, err = usersTable.Update(map[string]any{"age": 26}, map[string]any{"id": alice.ID})
	if err != nil {
		return err
	}

	_, err = usersTable.Delete(map[string]any{"id": alice.ID})
	if err != nil {
		return err
	}

	fmt.Printf("rows=%d count=%d pages=%d searchPages=%d searchCount=%d\n", len(rows), count, pages.TotalPages, searchPages.TotalPages, searchCount)

	if err := usersTable.DeleteTable(); err != nil {
		return err
	}
	return nil
}
