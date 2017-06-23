package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"reflect"
	"strings"
	"tc/stat"
	"time"
)

/*
(let ((class "statsByCampaignsHourly")
      (file "stats_by_campaigns_hourly.go"))
  (copy-file "stats_by_campaigns.go" file 1)
  (pop-to-buffer (find-file file))
  (insert "// generated from tables.go\n")
  (while (search-forward "statsByCampaigns" nil t)
    (replace-match class nil t)))
*/

type (
	Table interface {
		Add(stat.Slice)
		Rows() Rows
		UpdateSql() string
		InsertSql() string
		Name() string
	}
	ComplexTable interface {
		Table
		BeforeSql() string
	}
	Rows       []interface{}
	ComplexRow interface {
		Update()
	}
	table struct{}
)

func truncateDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func makeInsertSql(tbl interface{}, tblName string) string {
	t := reflect.TypeOf(tbl)
	fields, conditions := parseTypeSql(t, false)

	fields = append(fields, conditions...)
	placeholders := make([]string, len(fields))

	for i, f := range fields {
		fields[i] = fmt.Sprintf(`"%s"`, f)
		placeholders[i] = ":" + strings.ToLower(f)
	}
	fieldsSql := strings.Join(fields, ", ")
	placeholdersSql := strings.Join(placeholders, ", ")
	conditionsSql := makeConditionsSql(conditions)

	return fmt.Sprintf(`INSERT INTO tc."%s" (%s) `+
		`SELECT %s WHERE NOT EXISTS (SELECT 1 FROM tc."%s" WHERE %s)`,
		tblName, fieldsSql, placeholdersSql, tblName, conditionsSql)
}

func makeUpdateSql(tbl interface{}, tblName string) string {
	t := reflect.TypeOf(tbl)
	fields, conditions := parseTypeSql(t, true)

	for i, f := range fields {
		fields[i] = fmt.Sprintf(`"%s" = "%s" + :%s`, f, f, strings.ToLower(f))
	}
	fieldsSql := strings.Join(fields, ", ")
	conditionsSql := makeConditionsSql(conditions)

	return fmt.Sprintf(`UPDATE tc."%s" SET %s WHERE %s`, tblName, fieldsSql, conditionsSql)
}

func parseTypeSql(t reflect.Type, noupdate bool) (fields, conditions []string) {
	addConditions := func(t reflect.Type) {
		for i := 0; i < t.NumField(); i++ {
			conditions = append(conditions, t.Field(i).Name)
		}
	}

	var iterate func(t reflect.Type)
	iterate = func(t reflect.Type) {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			switch field.Tag.Get("table") {
			case "-":
				continue
			case "key":
				addConditions(field.Type)
			case "noupdate":
				if noupdate {
					continue
				}
				fallthrough
			default:
				if field.Type.Kind() == reflect.Struct && field.Type.Name() != "Time" {
					iterate(field.Type)
					continue
				}
				fields = append(fields, field.Name)
			}
		}
	}
	iterate(t)
	return
}

func makeConditionsSql(conditions []string) string {
	for i, c := range conditions {
		conditions[i] = fmt.Sprintf(`"%s" = :%s`, c, strings.ToLower(c))
	}
	return strings.Join(conditions, " AND ")
}

func makeTables() []Table {
	return []Table{
		&statsByAdsTable{},

		&statsByCampaignsTable{},
		&statsByCampaignsHourlyTable{},
		&statsByCountryForCampaignsTable{},

		&statsBySitesTable{},
		&statsBySitesHourlyTable{},
		&statsByCountryForSitesTable{},

		&statsByUsersTable{},

		&campaignsStatsTable{},
		&taxsTable{},
		&referralsTable{},
	}
}

func saveTables(tables []Table) {
	tx, err := db.Beginx()
	if err != nil {
		log.Println(err)
		return
	}
	defer tx.Rollback()

	for _, table := range tables {
		updateSql := table.UpdateSql()
		update, err := tx.PrepareNamed(updateSql)
		if err != nil {
			log.Println(err, "\n", updateSql)
			return
		}

		insertSql := table.InsertSql()
		insert, err := tx.PrepareNamed(insertSql)
		if err != nil {
			log.Println(err, "\n", insertSql)
			return
		}

		var before *sqlx.NamedStmt
		if complex, ok := table.(ComplexTable); ok {
			var err error
			beforeSql := complex.BeforeSql()
			before, err = tx.PrepareNamed(beforeSql)
			if err != nil {
				log.Println(err, "\n", beforeSql)
				return
			}
		}

		rows := table.Rows()
		if config.verbosity > VerbosityNone {
			log.Printf("Saving %d rows into the %s table\n", len(rows), table.Name())
		}
		for _, row := range rows {
			if complexRow, ok := row.(ComplexRow); ok {
				complexRow.Update()
			}

			// before
			if before != nil {
				_, err := before.Exec(row)
				if err != nil {
					log.Println(err)
					return
				}
			}
			// try update
			res, err := update.Exec(row)
			if err != nil {
				log.Println(err)
				return
			}
			n, err := res.RowsAffected()
			if err != nil {
				log.Printf("%v\n%+v\n%+v\n", err, row, updateSql)
				return
			}

			if n != 0 {
				continue
			}
			// insert elsewhere
			_, err = insert.Exec(row)
			if err != nil {
				log.Printf("%v\n%+v\n%+v\n", err, row, insertSql)
				return
			}
		}
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
	}
}
