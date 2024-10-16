package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli/v2"
)

func translate_ident_type(itype string) int8 {
	switch itype {
	case "u":
		return IDENT_TYPE_URL
	case "i":
		return IDENT_TYPE_IDENT
	default:
		return IDENT_TYPE_FILENAME
	}
}

func main() {

	var (
		detailed_list bool
		ident_type    string
		download      bool
	)

	cfg := load_config()

	token := login(cfg.Section("").Key("username").String(), cfg.Section("").Key("password").String())

	app := &cli.App{
		Description: "Webshare.cz CLI",
		Commands: []*cli.Command{
			{
				Name:  "ls",
				Usage: "list files in given path",
				Action: func(cCtx *cli.Context) error {
					file_response := files(token, cCtx.Args().First(), true)
					w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
					for _, folder := range file_response.Folders {
						fmt.Fprint(w, "D")
						if detailed_list {
							fmt.Fprintf(w, "\t%s\t", folder.Ident)
						}
						fmt.Fprintf(w, "\t%s/%s\n", cCtx.Args().First(), folder.Name)
					}
					for _, file := range file_response.Files {
						fmt.Fprint(w, "F")
						if detailed_list {
							fmt.Fprintf(w, "\t%s\t%s", file.Ident, file.Size)
						}
						fmt.Fprintf(w, "\t%s/%s\n", cCtx.Args().First(), file.Name)
					}
					w.Flush()
					return nil
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "l",
						Usage:       "long list",
						Destination: &detailed_list,
					},
				},
			},
			{
				Name:  "rm",
				Usage: "remove file in given path",
				Action: func(cCtx *cli.Context) error {
					fmt.Println(remove_file(token, cCtx.Args().First(), translate_ident_type(ident_type)))
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "t",
						Value:       "f",
						Usage:       "type of file identification - f filename (default), u url, i ident",
						Destination: &ident_type,
					},
				},
			},
			{
				Name:  "get",
				Usage: "get download link for given file or download it (-d)",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().First() == "" {
						fmt.Println("No file given")
						return nil
					}
					file_link := file_link(token, cCtx.Args().First(), download, translate_ident_type(ident_type))
					if !download {
						fmt.Println(file_link)
					}
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "t",
						Value:       "f",
						Usage:       "type of file identification - f filename (default), u url, i ident",
						Destination: &ident_type,
					},
					&cli.BoolFlag{
						Name:        "d",
						Usage:       "download the file",
						Destination: &download,
					},
				},
			},
			{
				Name:  "upload",
				Usage: "upload given file",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().First() == "" {
						fmt.Println("No file given")
						return nil
					}
					fmt.Printf("Uploading %s - ", cCtx.Args().First())
					fmt.Println(normal_link(EMPTY_TOKEN, upload(token, cCtx.Args().First(), cCtx.Args().Get(1)), IDENT_TYPE_IDENT))
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "show user data",
				Action: func(cCtx *cli.Context) error {
					user_stats := user_data(token)
					fmt.Printf("%s (%s)\n", user_stats.Username, user_stats.Id)
					fmt.Println("Email: " + user_stats.Email)
					fmt.Println("VIP until " + user_stats.VipUntil)
					fmt.Printf("Usage: %s/%s\n", user_stats.Bytes, user_stats.Space)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
