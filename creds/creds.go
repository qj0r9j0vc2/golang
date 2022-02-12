package main

import (
	"flag"
	"fmt"
	"github.com/awoodbeck/gnp/ch07/creds/auth"
	"log"
	"net"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
)

func init() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(),
			"Usage:\n\t%s <group names>\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func parseGroupNames(args []string) map[string]struct{} {
	groups := make(map[string]struct{})

	for _, args := range args {
		grp, err := user.LookupGroup(args)
		if err != nil {
			log.Print(err)
			continue
		}

		groups[grp.Gid] = struct{}{}
	}

	return groups
}

func main() {
	flag.Parse()

	groups := parseGroupNames(flag.Args())
	socket := filepath.Join(os.TempDir(), "creds.sock")
	addr, err := net.ResolveUnixAddr("unix", socket)
	if err != nil {
		log.Fatal(err)
	}

	s, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_ = s.Close()
	}()

	fmt.Printf("Listening on %s ...\n", socket)

	for {
		conn, err := s.AcceptUnix()
		if err != nil {
			break
		}
		if auth.Allowed(conn, groups) {
			_, err = conn.Write([]byte("Welcome\n"))
			if err == nil {

				continue
			}
		} else {
			_, err = conn.Write([]byte("Access denied\n"))
		}
		if err != nil {
			log.Println(err)
		}
		_ = conn.Close()
	}

}
