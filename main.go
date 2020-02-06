package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"search/domains/organizations"
	"search/domains/organizations/orgDao"
	"search/domains/tickets"
	ticketsDao "search/domains/tickets/dao"
	"search/domains/users"
	usersDao "search/domains/users/dao"
	"search/web"
)

type Config struct {
	Port int
}

var config Config

func main() {
	app := cli.NewApp()
	app.Name = "Search"
	app.Description = "Used to search fields in organization, users, and tickets"
	app.Commands = []cli.Command{
		{
			Name: "app",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "port",
					Value:       3000,
					EnvVar:      "PORT",
					Destination: &config.Port,
				},
			},
			Action: startApp,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func startApp(_ *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&orgDao.Organization{}).Error
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&ticketsDao.Ticket{}).Error
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&usersDao.User{}).Error
	if err != nil {
		panic(err)
	}
	orgDAO := orgDao.New(db)
	usersDAO := usersDao.New(db)
	ticketsDAO := ticketsDao.New(db)

	err = importOrganizations(ctx, orgDAO)
	if err != nil {
		return err
	}
	logger.Info("successfully imported organizations into database")

	err = importTickets(ctx, ticketsDAO)
	if err != nil {
		return err
	}
	logger.Info("successfully imported tickets into database")

	err = importUsers(ctx, usersDAO)
	if err != nil {
		return err
	}
	logger.Info("successfully imported users into database")

	server := web.HTTPServer{
		OrgFinder:     orgDAO,
		UsersFinder:   usersDAO,
		TicketsFinder: ticketsDAO,
	}

	return runServer(ctx, logger, server)
}

func runServer(ctx context.Context, logger *zap.Logger, server web.HTTPServer) error {
	addr := fmt.Sprintf(":%v", config.Port)
	r := web.MakeHandler(server)
	s := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// launch server in background
	go func() {
		logger.Info("starting web server", zap.Int("port", config.Port))
		if err := s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	// trap signals to allow for graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, os.Kill)

	// block and wait
	<-stopChan
	logger.Info("caught sig int; begin server shutdown")

	return s.Shutdown(ctx)
}

func importOrganizations(ctx context.Context, orgDao *orgDao.DAO) error {
	jsonFile, err := os.Open("organization.json")
	if err != nil {
		return errors.New("unable to open organization json file")
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var orgs []organizations.Organization
	err = json.Unmarshal(byteValue, &orgs)
	if err != nil {
		return errors.New("unable to unmarshal organization")
	}
	err = orgDao.ImportAll(ctx, organizations.ImportAllInput{Orgs: orgs})
	if err != nil {
		return errors.New("unable to import organizations. please try again")
	}
	return nil
}

func importUsers(ctx context.Context, usersDao *usersDao.DAO) error {
	jsonFile, err := os.Open("users.json")
	if err != nil {
		return errors.New("unable to open users json file")
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var u []users.User
	err = json.Unmarshal(byteValue, &u)
	if err != nil {
		return errors.New("unable to unmarshal users")
	}

	err = usersDao.ImportAll(ctx, users.ImportAllInput{Users: u})
	if err != nil {
		return errors.New("unable to import organizations. please try again")
	}
	return nil
}

func importTickets(ctx context.Context, ticketsDao *ticketsDao.DAO) error {
	jsonFile, err := os.Open("tickets.json")
	if err != nil {
		return errors.New("unable to tickets users json file")
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var t []tickets.Ticket
	err = json.Unmarshal(byteValue, &t)
	if err != nil {
		return errors.New("unable to unmarshal users")
	}

	err = ticketsDao.ImportAll(ctx, tickets.ImportAllInput{Tickets: t})
	if err != nil {
		return errors.New("unable to import organizations. please try again")
	}
	return nil
}
