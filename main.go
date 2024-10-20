package main

import (
	"fmt"
	"os"
	"strconv"
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
		public_folder bool
	)

	cfg := load_config()

	token := login(cfg.Section("").Key("username").String(), cfg.Section("").Key("password").String())

	app := &cli.App{
		Usage:                  "Webshare.cz CLI",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
				Name:      "ls",
				Usage:     "list files in given path",
				UsageText: "wscli ls [command options] PATH",
				Action: func(cCtx *cli.Context) error {
					file_response := files(token, cCtx.Args().First(), !public_folder)
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
					&cli.BoolFlag{
						Name:        "p",
						Usage:       "list files in public folder",
						Destination: &public_folder,
					},
				},
			},
			{
				Name:      "rm",
				Usage:     "remove file in given path",
				UsageText: "wscli rm [command options] PATH/IDENT",
				Action: func(cCtx *cli.Context) error {
					fmt.Println(remove_file(token, cCtx.Args().First(), translate_ident_type(ident_type), !public_folder))
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "t",
						Value:       "f",
						Usage:       "type of file identification - f filename, i ident",
						Destination: &ident_type,
					},
					&cli.BoolFlag{
						Name:        "p",
						Usage:       "remove from public folder",
						Destination: &public_folder,
					},
				},
			},
			{
				Name:      "get",
				Usage:     "get download link for given file or download it (-d)",
				UsageText: "wscli get [command options] PATH/IDENT/URL",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().First() == "" {
						fmt.Println("No file given")
						return nil
					}
					file_link := file_link(token, cCtx.Args().First(), download, translate_ident_type(ident_type), !public_folder)
					if !download {
						fmt.Println(file_link)
					}
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "t",
						Value:       "f",
						Usage:       "type of file identification - f filename, u url, i ident",
						Destination: &ident_type,
					},
					&cli.BoolFlag{
						Name:        "d",
						Usage:       "download the file",
						Destination: &download,
					},
					&cli.BoolFlag{
						Name:        "p",
						Usage:       "download from public folder",
						Destination: &public_folder,
					},
				},
			},
			{
				Name:      "upload",
				Usage:     "upload given local file to remote path",
				UsageText: "wscli upload [command options] LOCAL_PATH PATH",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().First() == "" {
						fmt.Println("No file given")
						return nil
					}
					fmt.Printf("Uploading %s - ", cCtx.Args().First())
					fmt.Println(normal_link(EMPTY_TOKEN, upload(token, cCtx.Args().First(), cCtx.Args().Get(1), !public_folder), IDENT_TYPE_IDENT, !public_folder))
					return nil
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "p",
						Usage:       "upload as public",
						Destination: &public_folder,
					},
				},
			},
			{
				Name:      "status",
				Usage:     "show user data",
				UsageText: "wscli status",
				Action: func(cCtx *cli.Context) error {
					user_stats := user_data(token)
					fmt.Printf("%s (%s)\n", user_stats.Username, user_stats.Id)
					fmt.Println("Email: " + user_stats.Email)
					fmt.Println("VIP until " + user_stats.VipUntil)
					used, _ := strconv.Atoi(user_stats.Bytes)
					space, _ := strconv.Atoi(user_stats.Space)
					fmt.Printf("Usage: %s/%s (%.2f%%)\n", user_stats.Bytes, user_stats.Space, 100*float64(used)/float64(space))
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
