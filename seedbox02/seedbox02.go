// Package seedbox02 handles all link generation for ddl1.project-gxs.com and tor1.project-gxs.com
package seedbox02

import (
	"fmt"
	"os"
	"path"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"bufio"
	"strings"
	"regexp"
	"net/http"
)

func GenerateDDL(filter string, shorten bool, video bool, torrent bool, html string, secret []byte) {
	fmt.Println()
	fmt.Println("---------- Seedbox-02 ----------")

	video_files_url := "/mnt/data/downloads"
	torrent_files_url := "/mnt/data/torrents"

	if html == "batch" || html == "non-batch" {
		// tor_output := establishSSHConnection(torrent_files_url, secret)
		video_output := establishSSHConnection(video_files_url, secret)
		parseHTMLOutput(filter, shorten, html, video_output)
	} else {
		if video {
			output := establishSSHConnection(video_files_url, secret)
			parseVideoOutput(filter, shorten, output)
		}
		if torrent {
			output := establishSSHConnection(torrent_files_url, secret)
			parseTorrentOutput(filter, shorten, output)
		}
	}
}

func establishSSHConnection(url_path string, secret []byte) string {
	app_path, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// get path of the private key
	ssh_key_path := path.Dir(app_path) + "/ddl-gen-key-01"
	
	key, err := ioutil.ReadFile(ssh_key_path)
	if err != nil {
		fmt.Println(os.Stderr, "Error: Key not found in path. Please provide a valid key.\n", err)
		os.Exit(1)
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, secret)
	if err != nil {
		fmt.Println(os.Stderr, "Error: Invalid Secret.\n", err)
		os.Exit(1)
	}

	config := &ssh.ClientConfig{
		User: "lambchopz",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", "205.185.113.65:22", config)
	if err != nil {
		fmt.Println("ssh.Dial()")
		panic(err)
	}

	session, err := client.NewSession()
	if err != nil {
		fmt.Println("client.NewSession()")
		panic(err)
	}
	defer session.Close()

	fmt.Println()
	fmt.Println("Obtaining files from seedbox...")
	output, err := session.Output("ls -l -h --si -R " + url_path)
	if err != nil {
		fmt.Println("session.Output")
		panic(err)
	}
	fmt.Println()
	return string(output)
}

func getShorteLink(link string) string {
	res, err := http.Get("https://api.shorte.st/s/9960e42c58372da26445d69d348e0baf/" + link)
	if err != nil {
		panic(err)
	}

	shortened_body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}
	shortened_link := strings.Replace(string(shortened_body), `{"status":"ok","shortenedUrl":"`, "", 1)
	shortened_link = strings.Replace(shortened_link[:len(shortened_link) - 2], `\`, "", -1)

	return shortened_link
}


func parseVideoOutput(filter string, shorten bool, output string) {
	var title string

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(strings.ToLower(line), strings.ToLower(filter)) {
			if strings.HasPrefix(line, "/mnt/data/downloads/") {
				fmt.Println()

				title = strings.Replace(line, "/mnt/data/downloads/", "", 1)
				title = title[:len(title) - 1]
			} else if strings.Contains(line, ".mkv") && strings.Contains(line, "[project-gxs]") {
				line = line[strings.Index(line, "[project-gxs]"):]

				if title != "" && string(title[len(title) - 1]) != "/" {
					title = title + "/"
				}

				if shorten {
					shortened_link := getShorteLink("http://ddl2.project-gxs.com/" + title + line)
					
					fmt.Println(line)
					fmt.Println(shortened_link)
				} else {
					fmt.Println("http://ddl2.project-gxs.com/" + title + line)
				}
			}
		}
	}
}

func parseTorrentOutput(filter string, shorten bool, output string) {
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(strings.ToLower(line), strings.ToLower(filter)) {
			if strings.Contains(line, ".torrent") && strings.Contains(line, "[project-gxs]") {
				line = line[strings.Index(line, "[project-gxs]"):]
				if shorten {
					shortened_link := getShorteLink("http://tor2.project-gxs.com/" + line)
					
					fmt.Println(line)
					fmt.Println(shortened_link)
				} else {
					fmt.Println("http://tor2.project-gxs.com/" + line)
				}
			}
		}
	}
}

func parseHTMLOutput(filter string, shorten bool, html string, video_output string) {
	var title string
	var title2 string
	var found bool
	var total_size string
	var file_size string

	scanner := bufio.NewScanner(strings.NewReader(string(video_output)))
	for scanner.Scan() {
		line := scanner.Text()

		if html == "batch" && found && strings.HasPrefix(line, "total ") {
			total_size = strings.Replace(line, "total ", "", 1)
			if total_size[len(total_size) - 1] == 'M' {
				total_size = strings.Replace(total_size, "M", " <em>MB</em>", 1)
			} else if total_size[len(total_size) - 1] == 'G' {
				total_size = strings.Replace(total_size, "G", " <em>GB</em>", 1)
			}
			found = false
		}

		if strings.Contains(strings.ToLower(line), strings.ToLower(filter)) {
			if strings.HasPrefix(line, "/mnt/data/downloads/") {
				found = true
				if html == "batch" && title != "" && total_size != "" {
					fmt.Println(`<p style="text-align: center;"></p>`)
					fmt.Print(`<p style="text-align: center;"><strong>Batch</strong> – (` + total_size + `): <strong><a href="`)

					if title != "" && string(title[len(title) - 1]) == "/" {
						title2 = title[:len(title) - 1]
					} else {
						title2 = title
					}

					if shorten {
						shortened_link := getShorteLink("http://tor2.project-gxs.com/" + title2 + ".torrent")
						
						fmt.Print(shortened_link)
					} else {
						fmt.Print("http://tor2.project-gxs.com/" + title2 + ".torrent")
					}
					fmt.Println(`" target="_blank">Torrent</a></strong></p>`)
					total_size = ""
				}
				fmt.Println()

				title = strings.Replace(line, "/mnt/data/downloads/", "", 1)
				title = title[:len(title) - 1]
				fmt.Println(title)
			} else if strings.Contains(line, ".mkv") && strings.Contains(line, "[project-gxs]") {
				file_size = strings.Fields(line)[4]
				if file_size[len(file_size) - 1] == 'M' {
					file_size = strings.Replace(file_size, "M", " <em>MB</em>", 1)
				} else if file_size[len(file_size) - 1] == 'G' {
					file_size = strings.Replace(file_size, "G", " <em>GB</em>", 1)
				} else if file_size[len(file_size) - 1] == 'K' {
					file_size = strings.Replace(file_size, "K", " <em>KB</em>", 1)
				} else {
					file_size = file_size + " <em>B</em>"
				}
				line = line[strings.Index(line, "[project-gxs]"):]

				re := regexp.MustCompile(`-[^-]+?\[`)
				episode_num := re.FindString(strings.Replace(line, "[project-gxs]", "", 1))

				if episode_num != "" {
					episode_num = episode_num[1:len(episode_num) - 1]
					episode_num = strings.TrimSpace(episode_num)
				}

				
				fmt.Print(`<p style="text-align: center;">Episode – <strong>` + episode_num + `</strong> – (` + file_size + `): <strong><a href="`)

				if title != "" && string(title[len(title) - 1]) != "/" {
					title2 = title + "/"
				}

				if shorten {
					shortened_link := getShorteLink("http://ddl2.project-gxs.com/" + title2 + line)
						
					fmt.Print(shortened_link)
				} else {
					fmt.Print("http://ddl2.project-gxs.com/" + title2 + line)
				}
				fmt.Print(`" target="_blank" rel="noopener">DDL</a></strong>`)

				if html == "non-batch" {
					fmt.Print(` | <strong><a href="`)
					if shorten {
						shortened_link := getShorteLink("http://tor2.project-gxs.com/" + line + ".torrent")
						
						fmt.Print(shortened_link)
					} else {
						fmt.Print(`http://tor2.project-gxs.com/` + line + ".torrent")
					}
					fmt.Println(`" target="_blank" rel="noopener">Torrent</a></strong></p>`)
				} else if html == "batch" {
					fmt.Print(`</p>`)
					fmt.Println()
				} else {
					fmt.Print(`</p>`)
				}
			}
		}
	}
	if html == "batch" && title != "" && total_size != "" {
		fmt.Println(`<p style="text-align: center;"></p>`)
		fmt.Print(`<p style="text-align: center;"><strong>Batch</strong> – (` + total_size + `): <strong><a href="`)

		if title != "" && string(title[len(title) - 1]) == "/" {
			title2 = title[:len(title) - 1]
		} else {
			title2 = title
		}

		if shorten {
			shortened_link := getShorteLink("http://tor2.project-gxs.com/" + title2 + ".torrent")
			
			fmt.Print(shortened_link)
		} else {
			fmt.Print("http://tor2.project-gxs.com/" + title2 + ".torrent")
		}
		fmt.Println(`" target="_blank">Torrent</a></strong></p>`)
		total_size = ""
	}
}