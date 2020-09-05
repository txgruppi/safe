package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dgraph-io/badger/v2"

	"github.com/txgruppi/safe/buildinfo"
	"github.com/txgruppi/safe/database"
	"github.com/txgruppi/safe/errors"
	"github.com/txgruppi/safe/fs"
	"github.com/urfave/cli"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var db database.Database
	app := cli.App{
		Name:    "safe",
		Usage:   "a safe place to store files",
		Version: buildinfo.Version,
		Commands: []cli.Command{
			cli.Command{
				Name:  "db",
				Usage: "manage a database",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "database",
						Usage:    "the path to your database file",
						EnvVar:   "SAFE_DATABASE",
						Required: true,
					},
					cli.StringFlag{
						Name:   "password",
						Usage:  "the password to unlock your database",
						EnvVar: "SAFE_PASSWORD",
					},
					cli.BoolFlag{
						Name:   "verbose",
						EnvVar: "SAFE_VERBOSE",
					},
				},
				Before: func(c *cli.Context) error {
					var err error
					db, err = database.New(c.String("database"), c.Bool("verbose"))
					if err != nil {
						return err
					}
					if err := db.Unlock([]byte(c.String("password"))); err != nil {
						return err
					}
					return nil
				},
				After: func(c *cli.Context) error {
					if err := db.Tidy(); err != nil {
						return err
					}
					if err := db.Lock(); err != nil {
						return err
					}
					return nil
				},
				Subcommands: []cli.Command{
					cli.Command{
						Name:  "ls",
						Usage: "list database entries",
						Action: func(c *cli.Context) error {
							it, err := db.Iterator("/", false)
							if err != nil {
								return err
							}
							defer it.Close()
							for file := range it.Channel() {
								fmt.Printf("%s %d %s\n", file.Location(), file.Size(), file.MimeType())
							}
							return it.Error()
						},
					},
					cli.Command{
						Name:      "put",
						Usage:     "put one or more files in the database",
						UsageText: "safe db put <database location> <local file>[ <database location> <local file> [...]]",
						Action: func(c *cli.Context) error {
							args := c.Args()
							if len(args) == 0 {
								return errors.ErrMissingFileArgument
							}
							if len(args)%2 != 0 {
								return errors.ErrWrongNumberOfArguments
							}
							for i := 0; i < len(args); i += 2 {
								data, err := ioutil.ReadFile(args[i+1])
								if err != nil {
									return err
								}
								file := fs.NewEmptyFile().
									SetLocation(args[i]).
									SetMimeType(fs.SafeMimeType(args[i+1])).
									SetSize(int64(len(data))).
									SetData(data)
								if err := db.Set(file); err != nil {
									return err
								}
							}
							return nil
						},
					},
					cli.Command{
						Name:      "get",
						Usage:     "get one or more files from the database",
						UsageText: "safe db get <database location> <local file>[ <database location> <local file> [...]]",
						Action: func(c *cli.Context) error {
							args := c.Args()
							if len(args) == 0 {
								return errors.ErrMissingFileArgument
							}
							if len(args)%2 != 0 {
								return errors.ErrWrongNumberOfArguments
							}
							for i := 0; i < len(args); i += 2 {
								file, err := db.Get(args[i])
								if err == badger.ErrKeyNotFound {
									continue
								}
								if err != nil {
									return err
								}
								fp, err := os.OpenFile(args[i+1], os.O_WRONLY|os.O_EXCL|os.O_CREATE, 0600)
								if err != nil {
									return err
								}
								_, err = fp.Write(file.Data())
								if err := fp.Close(); err != nil {
									return err
								}
								if err != nil {
									return err
								}
							}
							return nil
						},
					},
					cli.Command{
						Name:      "rm",
						Usage:     "remove one or more files from the database",
						UsageText: "safe db rm <database location> [ <database location> [...]]",
						Action: func(c *cli.Context) error {
							args := c.Args()
							if len(args) == 0 {
								return errors.ErrMissingFileArgument
							}
							for i := 0; i < len(args); i++ {
								if err := db.Del(args[i]); err != nil {
									return nil
								}
							}
							return nil
						},
					},
				},
			},
		},
	}
	return app.Run(os.Args)
}
