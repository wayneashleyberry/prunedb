package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()

	cmd := Command(ctx)

	cmd.SilenceUsage = true

	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const prefix = "TESTDB"

type Specification struct {
	Host     string `envconfig:"HOST" default:"127.0.0.1"`
	Port     string `envconfig:"PORT" default:"3306"`
	User     string `envconfig:"USER" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
}

// Config creates a new mysql.Config from a TestDatabase
func (s Specification) Config() *mysql.Config {
	conf := mysql.NewConfig()

	conf.Addr = s.Host + ":" + s.Port
	conf.Net = "tcp"
	conf.User = s.User
	conf.Passwd = s.Password
	conf.ParseTime = true
	conf.Params = map[string]string{
		"charset":         "utf8mb4",
		"multiStatements": "true",
	}

	return conf
}

// Command will create a new "prunedb" cobra command.
func Command(ctx context.Context) *cobra.Command {
	var commit bool

	cmd := &cobra.Command{
		Use:   "prunedb",
		Short: "Drop lingering test databases",
		RunE: func(cmd *cobra.Command, args []string) error {
			var s Specification

			err := envconfig.Process(prefix, &s)
			if err != nil {
				fmt.Println(err)

				envconfig.Usage(prefix, &s)

				os.Exit(1)
			}

			conf := s.Config()

			db, err := sqlx.ConnectContext(ctx, "mysql", conf.FormatDSN())
			if err != nil {
				return fmt.Errorf("connect to database: %w", err)
			}

			rows, err := db.QueryxContext(ctx, "show databases")
			if err != nil {
				return fmt.Errorf("execute database query: %w", err)
			}

			defer func() {
				_ = rows.Close()
			}()

			databases := []string{}

			for rows.Next() {
				var database string

				err := rows.Scan(&database)
				if err != nil {
					return fmt.Errorf("scan row: %w", err)
				}

				if strings.HasPrefix(database, "test_") {
					databases = append(databases, database)
				}
			}

			if err := rows.Err(); err != nil {
				return fmt.Errorf("rows error: %w", err)
			}

			if len(databases) == 0 {
				fmt.Println("no databases to prune")

				return nil
			}

			for _, database := range databases {
				query := fmt.Sprintf("drop database %s", database)

				fmt.Println(query)

				if !commit {
					continue
				}

				_, err := db.ExecContext(ctx, query)
				if err != nil {
					return fmt.Errorf("execute database query: %w", err)
				}
			}

			if !commit {
				fmt.Println("use the --commit flag to run these queries")
			}

			return nil
		},
	}

	cmd.PersistentFlags().BoolVar(&commit, "commit", false, "commit query")

	return cmd
}
