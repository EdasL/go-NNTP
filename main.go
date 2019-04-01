package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/EdasL/NNTP/nntpclient"
)

func main() {

	serverAddr := flag.String("addr", "-", "nntp server address")
	flag.Parse()

	client, err := nntpclient.New("tcp", *serverAddr)
	if err != nil {
		log.Fatal("Failed to create a new conncetion: ", err)
	}
	defer client.Close()
	log.Printf("Got banner:  %v", client.Banner)

	fmt.Println("Possible commands:( ARTICLE, GROUP, HELP, IHAVE, LAST, LIST, NEWGROUPS, NEWNEWS, NEXT, POST, STAT, QUIT)")

	for {

		var command string

		if _, err := fmt.Scanln(&command); err != nil {
			log.Fatal(err)
		}
		switch command {
		case "ARTICLE":
			fmt.Println("Enter artical name or id-number:")
			var id string
			if _, err := fmt.Scanln(&id); err != nil {
				fmt.Println("Failed to scan artical id: ", err)
				break
			}

			r, err := client.Article(id)
			if err != nil {
				log.Println("Failed to select and article: ", err)
				break
			}
			if _, err = io.Copy(os.Stdout, r); err != nil {
				log.Println("Failed to print selected article: ", err)
				break
			}
			break
		case "BODY":
			fmt.Println("Enter artical name or id-number:")
			var id string
			if _, err := fmt.Scanln(&id); err != nil {
				fmt.Println("Failed to scan artical id: ", err)
				break
			}

			r, err := client.Body(id)
			if err != nil {
				log.Println("Failed to select and articles body: ", err)
				break
			}
			if _, err = io.Copy(os.Stdout, r); err != nil {
				log.Println("Failed to print selected articles body: ", err)
				break
			}

			break
		case "GROUP":
			fmt.Println("Enter the name of the group:")
			var name string
			if _, err := fmt.Scanln(&name); err != nil {
				fmt.Println("Failed to scan group name: ", err)
				break
			}
			// Select a group
			g, err := client.Group(name)
			if err != nil {
				fmt.Println("Failed to select a group: ", err)
				break
			}
			fmt.Println("Group selected: ", g)
			break
		case "HEAD":
			fmt.Println("Enter artical name or id-number:")
			var id string
			if _, err := fmt.Scanln(&id); err != nil {
				fmt.Println("Failed to scan artical id: ", err)
				break
			}

			r, err := client.Head(id)
			if err != nil {
				log.Println("Failed to select and articles head: ", err)
				break
			}
			if _, err = io.Copy(os.Stdout, r); err != nil {
				log.Println("Failed to print selected articles head: ", err)
				break
			}

			break
		case "HELP":
			r, err := client.Help()
			if err != nil {
				log.Println("Help command failed: ", err)
				break
			}
			if _, err = io.Copy(os.Stdout, r); err != nil {
				log.Println("Failed to print commands: ", err)
				break
			}
			break
		case "IHAVE":
			fmt.Println("Enter artical name or id-number(example 1, 2 ... <4106@ucbvax.ARB>):")
			var id string
			if _, err := fmt.Scanln(&id); err != nil {
				fmt.Println("Failed to scan artical id: ", err)
				break
			}
			const examplepost = `From: <test@example.com>
Newsgroups: misc.test
Subject: Code test
Organization: testers

This is a test post.
`
			// Post an article
			resp, err := client.Ihave(strings.NewReader(examplepost), id)
			if err != nil {
				log.Println("Failed to post article or server doesn't want the article: ", err)
				break
			}

			fmt.Println(*resp)
			break
		case "LAST":
			if _, err := client.Last(); err != nil {
				log.Println("Failed LAST command: ", err)
				break
			}
			break
		case "LIST":
			r, err := client.List()
			if err != nil {
				log.Println("Failed to list groups: ", err)
				break
			}

			if _, err = io.Copy(os.Stdout, r); err != nil {
				log.Println("Failed to print the list: ", err)
				break
			}
			break
		case "NEWGROUPS":
			fmt.Println("Enter date [YY]YYMMDD:")
			var date string
			if _, err := fmt.Scanln(&date); err != nil {
				fmt.Println("Failed to scan date: ", err)
				break
			}

			fmt.Println("Enter time HHMMSS:")
			var time string
			if _, err := fmt.Scanln(&time); err != nil {
				fmt.Println("Failed to scan time: ", err)
				break
			}

			r, err := client.Newgroups(date, time)
			if err != nil {
				log.Println("Failed to list: ", err)
				break
			}

			if _, err = io.Copy(os.Stdout, r); err != nil {
				log.Println("Failed to print the list: ", err)
				break
			}
			break
		case "NEWNEWS":
			fmt.Println("Enter newsgroup(* for all):")
			var group string
			if _, err := fmt.Scanln(&group); err != nil {
				fmt.Println("Failed to scan newsgroup: ", err)
				break
			}

			fmt.Println("Enter date [YY]YYMMDD:")
			var date string
			if _, err := fmt.Scanln(&date); err != nil {
				fmt.Println("Failed to scan date: ", err)
				break
			}

			fmt.Println("Enter time HHMMSS:")
			var time string
			if _, err := fmt.Scanln(&time); err != nil {
				fmt.Println("Failed to scan time: ", err)
				break
			}

			r, err := client.Newnews(group, date, time)
			if err != nil {
				log.Println("Failed to list new news: ", err)
				break
			}

			if _, err = io.Copy(os.Stdout, r); err != nil {
				log.Println("Failed to print the list: ", err)
				break
			}
			break
		case "NEXT":
			if _, err := client.Next(); err != nil {
				log.Println("Failed NEXT command: ", err)
				break
			}
			break
		case "POST":
			const examplepost = `From: <test@example.com>
Newsgroups: misc.test
Subject: Code test
Organization: testers

This is a test post.
`
			// Post an article
			resp, err := client.Post(strings.NewReader(examplepost))
			if err != nil {
				log.Println("Failed to post article: ", err)
				break
			}

			fmt.Println(*resp)
			break
		case "STAT":
			fmt.Println("Enter artical name or id-number:")
			var id string
			if _, err := fmt.Scanln(&id); err != nil {
				fmt.Println("Failed to scan artical id: ", err)
				break
			}

			msg, err := client.Stat(id)
			if err != nil {
				log.Println("Failed to select and articles statistics: ", err)
				break
			}

			fmt.Println(*msg)
			break
		case "QUIT":
			msg, err := client.Quit()
			if err != nil {
				log.Println("Failed to quit: ", err)
				break
			}
			fmt.Println(*msg)
			break
		default:
			fmt.Printf("Command: %v doesn't exist \n", command)
			break
		}
	}
}
