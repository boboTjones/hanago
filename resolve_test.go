package resolve

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	fqdn := "assets.98point6.com."
	var T []net.IP

	res, err := Resolve(fqdn, "8.8.8.8:53")
	if err != nil {
		//XXX ToDo(erin): haven't had an error yet so print it for now.
		fmt.Println(err)
	}

	for _, rec := range res.Records {
		fmt.Println(rec)
	}

	assert.Equal(t, fqdn, res.Original)
	if len(res.Records) > 0 && res.Records[0].Type == "A" {
		assert.IsType(t, T, res.Records[0].IPs)
	}
}

func TestIsWildcard(t *testing.T) {
	yes, _ := IsWildCard("sockpuppet.org.", "8.8.8.8:53")
	assert.True(t, yes)

}

func TestBrute(t *testing.T) {
	wl := []string{"bystander", "bystreet", "byth", "bytime", "bytownite", "bytownitite", "bywalk", "bywalker", "byway", "bywoner", "byword", "bywork"}

	_ = Brute("sockpuppet.org", "8.8.8.8:53", wl)
}
