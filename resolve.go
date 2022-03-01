package resolve

import (
	"crypto/rand"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type Record struct {
	Name  string
	Type  string
	IPs   []net.IP
	Cname string
	Time  time.Duration
}

type Answers struct {
	Original string
	Records  []Record
}

func Resolve(name, server string) (Answers, error) {
	recs := []Record{}

	if !strings.HasSuffix(name, ".") {
		name = name + "."
	}

	msg := new(dns.Msg)
	msg.SetQuestion(name, dns.TypeA)
	client := new(dns.Client)

	reply, t, err := client.Exchange(msg, server)
	if err != nil {
		return Answers{
			Original: name,
			Records:  recs,
		}, err
	}

	for i, r := range reply.Answer {
		switch r.Header().Rrtype {
		case 1:
			wtf := 0
			if ar, ok := r.(*dns.A); ok {
				for i := range recs {
					if recs[i].Type == "A" && recs[i].Name == ar.Header().Name {
						recs[i].IPs = append(recs[i].IPs, ar.A)
						wtf++
						break
					}
				}
				if wtf == 0 {
					record := Record{
						Name: ar.Header().Name,
						Type: "A",
						IPs:  make([]net.IP, 0),
						Time: t,
					}
					record.IPs = append(record.IPs, ar.A)
					recs = append(recs, record)
				}
			} else {
				fmt.Println(ar)
			}
		case 5:
			if ar, ok := r.(*dns.CNAME); ok {
				rec := Record{
					Name:  r.Header().Name,
					Type:  "CNAME",
					Cname: ar.Target,
					Time:  t,
				}
				recs = append(recs, rec)
			} else {
				fmt.Println(ar)
			}
		default:
			fmt.Println(reply.Answer[i].Header().String())
		}

	}

	return Answers{
		Original: name,
		Records:  recs,
	}, nil
}

func IsWildCard(sld, server string) (bool, error) {
	if !strings.HasSuffix(sld, ".") {
		sld = sld + "."
	}

	x := 0
	var rbuf [10]byte
	//nip := make([]net.IP, 0)

	// Get a baseline.
	base, err := Resolve("www."+sld, server)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	if len(base.Records) > 0 {
		for _, rec := range base.Records {
			fmt.Println(rec)
		}
	} else {
		err = fmt.Errorf("No answers for www.%s\n", sld)
		return false, err
	}

	for i := 1; i <= 4; i++ {
		_, err := rand.Read(rbuf[:])
		if err != nil {
			return false, err
		}

		fqdn := fmt.Sprintf("%x.%s", rbuf[:], sld)
		res, err := Resolve(fqdn, server)
		if err != nil {
			return false, err
		}

		if len(res.Records) != 0 {
			for _, answer := range res.Records {
				for _, ip := range answer.IPs {

					fmt.Println(ip)
				}

			}
		}
	}

	if x != 0 {
		return true, nil
	}

	return false, nil
}

func Brute(sld, server string, wl []string) error {
	if !strings.HasSuffix(sld, ".") {
		sld = sld + "."
	}

	for _, w := range wl {
		fqdn := fmt.Sprintf("%s.%s", w, sld)
		res, err := Resolve(fqdn, server)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}

	return nil
}
