package action

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // import file-based migrations driver

	_ "github.com/lib/pq" // import postgres driver
	"github.com/spf13/cobra"
)

type migrateCmdFlags struct {
	host     string
	port     string
	dbName   string
	sslMode  string
	user     string
	password string
	count    int
}

var migrateFlags migrateCmdFlags

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Runs database migrations for gocop",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=%s",
			migrateFlags.host, migrateFlags.port, migrateFlags.user, migrateFlags.password, migrateFlags.dbName, migrateFlags.sslMode)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			log.Fatal(err)
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://sql/migrations",
			"postgres", driver)
		if err != nil {
			log.Fatal(err)
		}

		if migrateFlags.count == 0 {
			err = m.Up()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = m.Steps(migrateFlags.count)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().StringVarP(&migrateFlags.host, "host", "a", "localhost", "database host")
	migrateCmd.Flags().StringVarP(&migrateFlags.port, "port", "t", "5432", "database port")
	migrateCmd.Flags().StringVarP(&migrateFlags.dbName, "database", "x", "postgres", "database name")
	migrateCmd.Flags().StringVarP(&migrateFlags.sslMode, "ssl", "y", "require", "database ssl mode")
	migrateCmd.Flags().StringVarP(&migrateFlags.password, "pass", "p", "", "database password")
	migrateCmd.Flags().StringVarP(&migrateFlags.user, "user", "u", "postgres", "database username")
	migrateCmd.Flags().IntVarP(&migrateFlags.count, "count", "c", 0, "number of migrations to apply")
	err := migrateCmd.MarkFlagRequired("pass")
	if err != nil {
		log.Fatal(err)
	}
}
