package Application

import (
	"gopkg.in/urfave/cli.v1"
	"fmt"
	"os"
	"gomypartition/Application/Action/Factory"
	"gomypartition/Application/Database"
	"gomypartition/Application/Service"
)

type Application struct {
	Console   *cli.App
	DefaultDb *Database.DbDriver
}

var databaseFlags = []cli.Flag {
	cli.StringFlag {
		Name: "host",
		Usage: "REQUIRED database server hostname/ip address",
	},

	cli.IntFlag {
		Name: "port",
		Usage: "database server port",
		Value: 3306,
	},

	cli.StringFlag {
		Name: "user",
		Usage: "REQUIRED database username",
	},

	cli.StringFlag {
		Name: "password",
		Usage: "REQUIRED database password",
	},

	cli.StringFlag {
		Name: "database",
		Usage: "REQUIRED database schema name",
	},

	cli.StringFlag {
		Name: "table",
		Usage: "REQUIRED database table name",
	},
}

var infoSpecificFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "orderfld, o",
		Usage: "column for sorting partitions info records",
		Value: "PARTITION_ORDINAL_POSITION",
	},
}

var maintenanceSpecificFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "column",
		Usage: "REQUIRED column for partitioning",
	},
	cli.IntFlag{
		Name:  "max-partitions",
		Usage: "max. partition count",
		Value: 50,
	},
	cli.IntFlag{
		Name: "range",
		Usage: "partition range in days, e.g. 1, 7, 30 ...",
		Value: 30,
	},
	cli.IntFlag{
		Name: "retention",
		Usage: "partition retention in days. Partitions older than NOW() + <retention> will be removed",
		Value: 90,
	},
	cli.StringFlag{
		Name: "prefix",
		Usage: "partition name prefix, e.g. prefix `to` will be used in partition name as to_20180801",
		Value: "to",
	},
	cli.BoolFlag{
		Name: "dry-run",
		Usage: "No sql will be executed. It will be printed to output only.",
	},
}

func (application *Application) setConnectionFromCli(c *cli.Context) error {

	application.DefaultDb.Host = c.String("host")
	application.DefaultDb.Port = c.Int("port")
	application.DefaultDb.User = c.String("user")
	application.DefaultDb.Password = c.String("password")
	application.DefaultDb.Database = c.String("database")

	requiredFields := []string{"Host","User","Password","Database"}
	return Service.CheckRequiredFields(application.DefaultDb, requiredFields)
}

func (application *Application) init() {

	app := cli.NewApp()
	app.Name = "gomypartition"
	app.Usage = "primitive partioning tool for mysql databases"
	app.Version = "0.9.0"
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "Unknown command %q\n", command)
	}

	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelp(c)
		return nil
	}

	application.Console = app

	application.DefaultDb = &Database.DbDriver{}

	application.addActions()
}

func (application *Application) addActions() {

	application.Console.Commands = []cli.Command{
		{

			Name:    "docker-test",
			Aliases: []string{"dt"},
			Usage:   "Create partitioned table in docker (see .env settings) and insert random records there",
			Action:  func(c *cli.Context) error {

				action := Factory.NewDockerTestAction()
				err := action.Process()

				return handleErr(err)
			},
		},
		{
			Name:    "info",
			Usage:   "Read and display partition info for specified table",
			Flags: append(databaseFlags, infoSpecificFlags...),
			Action:  func(c *cli.Context) error {
				err := application.setConnectionFromCli(c)

				if err != nil {
					return handleErr(err)
				}

				action, err := Factory.NewInfoAction(application.DefaultDb, c)

				if err != nil {
					return handleErr(err)
				}

				return handleErr(action.Process())
			},
		},
		{
			Name:  "maintenance",
			Usage: "Process time-series based partition maintenance for specified table",
			Flags: append(databaseFlags, maintenanceSpecificFlags...),
			Action: func(c *cli.Context) error {

				err := application.setConnectionFromCli(c)

				if err != nil {
					return handleErr(err)
				}
				action, err := Factory.NewMaintenanceAction(application.DefaultDb, c)

				if err != nil {
					return handleErr(err)
				}

				return handleErr(action.Process())
			},
		},
	}
}

func (application *Application) Run() {
	application.init()
	application.Console.Run(os.Args)
}

func handleErr(err error) error {
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("\nError: %s\n", err.Error()), 1)
	}

	return nil
}
