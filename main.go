package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/urfave/cli.v1"
)

func main() {

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "login",
			Aliases: []string{"l"},
			Usage:   "Submit the token used for authenticating with your GitHub Token.",
			Action:  storeToken,
		},
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Used to create a Github Repository with the command line.",
			Action:  hubinit,
			Flags:[]cli.Flag{
				cli.StringFlag{
					Name: "organisation, o",
					Usage: "To set the organisation of the repository.",
				},
				cli.StringFlag{
					Name: "name, n",
					Usage: "Name of the repository.",
				},
				cli.BoolFlag{
					Name: "public, p",
					Usage: "To set the repo as public",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func storeToken(cli *cli.Context) error {
	//Save the token into the file.

	println(cli.Args().First())
	return nil
}


func authenticate(ctx context.Context) (*github.Client, error){

	//Get Token from file, if not present raise error
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ""},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc),nil
}

func hubinit(cli *cli.Context) error {

	ctx := context.Background()
	client,err := authenticate(ctx)
	rpath, err := os.Getwd()
	if err != nil{
		print("Error, Can't find repo name")
		return err
	}

	reponame := filepath.Base(rpath)
	if(cli.IsSet("name")){
		reponame = cli.String("name")
	}

	var private bool = true
	if(cli.IsSet("public")){
		private = false
	}

	org := ""
	if(cli.IsSet("organisation")) {
		org = cli.String("organisation")
	}

	_ , err = createRepo(ctx,client,reponame,org,private)
	if err != nil{
		return err
	}

	return nil
}

func createRepo(ctx context.Context, client *github.Client,name string,org string, private bool) (bool, error){

	repodetails := &github.Repository{
		Name:    &name,
		Private: &private,
	}

	_,_ , err := client.Repositories.Create(ctx,org,repodetails)
	if err!=nil{
		return false, err
	}

	return true,nil
}
