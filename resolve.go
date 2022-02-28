package resolve

import (
	"crypto/rand"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type Answer struct {
	Name   string
	Type   string
	IPs    []net.IP
	Cnames []string
	Time   time.Duration
}

type Answers struct {
	Host    string
	Answers []Answer
}

func Resolve(name, server string) (Answers, error) {
	answers := []Answer{}

	if !strings.HasSuffix(name, ".") {
		name = name + "."
	}

	msg := new(dns.Msg)
	msg.SetQuestion(name, dns.TypeA)
	client := new(dns.Client)

	reply, t, err := client.Exchange(msg, server)
	if err != nil {
		return Answers{Answers: answers}, err
	}

	for _, r := range reply.Answer {
		if ar, ok := r.(*dns.A); ok {
			if len(answers) > 0 {
				for i, answer := range answers {
					if answer.Name == r.Header().Name {
						answers[i].IPs = append(answers[i].IPs, ar.A)
					}
				}
			} else {
				a := Answer{
					Name: r.Header().Name,
					Type: "A",
					IPs:  make([]net.IP, 0),
					Time: t,
				}
				a.IPs = append(a.IPs, ar.A)
				answers = append(answers, a)
			}
		}
	}

	return Answers{
		Host:    name,
		Answers: answers,
	}, nil
}

func IsWildCard(sld, server string) (bool, error) {
	x := 0
	name := ""
	var rbuf [10]byte

	for i := 1; i <= 4; i++ {
		_, err := rand.Read(rbuf[:])
		if err != nil {
			return false, err
		}

		fqdn := fmt.Sprintf("%x.%s", rbuf[:], sld)
		res, err := Resolve(fqdn, server)

		if len(res.Answers) != 0 {
			answers := res.Answers

			if err != nil {
				return false, err
			}
			if name == answers[0].Name {
				x++
			}
			name = answers[0].Name
		}
	}

	if x != 0 {
		return true, nil
	}

	return false, nil
}
