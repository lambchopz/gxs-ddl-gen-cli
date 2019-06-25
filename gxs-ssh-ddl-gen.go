package main

import (
	"os"
	"flag"
	"fmt"
	"strconv"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"github.com/lambchopz/gxs-ddl-gen-cli/seedbox01"
	"github.com/lambchopz/gxs-ddl-gen-cli/seedbox02"
)

func main() {
	const app_version = "1.2.6"

	version := flag.Bool("version", false, "print version of the app")
	shorten := flag.Bool("shorten", true, "shortens generated DDL links using shorte. To disable (./gxs-ddl-gen -shorten=false")
	filter := flag.String("filter", "", "filter results in seedbox, wrap multi-word filters with quotes. ex: (./gxs-ddl-gen -filter=\"One Piece\")")
	video := flag.Bool("video", true, "toggles generation of DDL links for media files. To disable (./gxs-ddl-gen -video=false)")
	torrent := flag.Bool("torrent", true, "toggles generation of DDL links for torrent files. To disable (./gxs-ddl-gen -torrent=false)")
	html := flag.String("html", "", "generates html output of DDL links. This ignores the video and torrent flags. \nTwo options are available (-html=non-batch or -html=batch)")
	flag.Parse()
	
	if *version {
		fmt.Println("gxs_ddl_gen version " + app_version)
		os.Exit(0)
	}
	if flag.Lookup("filter") != nil {
		fmt.Println("filter=\"" + *filter + "\"")
	}
	if flag.Lookup("shorten") != nil {
		fmt.Println("shorten=" + strconv.FormatBool(*shorten))
	}
	if flag.Lookup("video") != nil {
		fmt.Println("video=" + strconv.FormatBool(*video))
	}
	if flag.Lookup("torrent") != nil {
		fmt.Println("torrent=" + strconv.FormatBool(*torrent))
	}
	if flag.Lookup("html") != nil {
		fmt.Println("html=\"" + *html + "\"")
		if *html != "" && *html != "non-batch" && *html != "batch" {
			panic("Invalid html flag value provided. use the -help or -h flag for a list of valid commands")
		}
	}

	fmt.Printf("Enter secret: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	fmt.Println()

	seedbox01.GenerateDDL(*filter, *shorten, *video, *torrent, *html, bytePassword)
	seedbox02.GenerateDDL(*filter, *shorten, *video, *torrent, *html, bytePassword)
	fmt.Println()
	fmt.Println("DDL Generation Complete")
}
