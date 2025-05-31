package schema_test

import (
	_ "github.com/mattn/go-sqlite3"
)

// func TestSchema(t *testing.T) {
// 	testDb, err := sql.Open("sqlite3", ":memory:")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer testDb.Close()

// 	m := schema.Migrator{
// 		Store:   schema.NewSqlite3SchemaStore(testDb),
// 		Log:     *slog.Default(),
// 		Sources: migrations.All,
// 	}
// 	err = m.ApplyLatest(t.Context())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	var b strings.Builder
// 	var stmts []string
// 	rows, err := testDb.Query("SELECT sql FROM sqlite_schema")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var stmt sql.NullString
// 		if err := rows.Scan(&stmt); err != nil {
// 			t.Fatal(err)
// 		}
// 		if stmt.Valid {
// 			el := stmt.String
// 			if !strings.HasSuffix(b.String(), ";") {
// 				el += ";"
// 			}
// 			stmts = append(stmts, el)
// 		}
// 	}

// 	want := db.Schema
// 	got := strings.Join(stmts, "\n")

// 	if diff := cmp.Diff(want, got); diff != "" {
// 		t.Errorf("schema mismatch (-want +got):\n%s", diff)
// 	}
// }

// func TestSchemaVersion(t *testing.T) {
// 	var ids []uint64
// 	for id, _ := range db.Migrations {
// 		ids = append(ids, id)
// 	}
// 	slices.Sort(ids)
// }
