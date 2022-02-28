package resolve

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	fqdn := "www.sockpuppet.org."
	res, _ := Resolve(fqdn, "8.8.8.8:53")
	var T []net.IP

	assert.Equal(t, fqdn, res.Host)
	if len(res.Answers) > 0 && res.Answers[0].Type == "A" {
		assert.IsType(t, T, res.Answers[0].IPs)
	}
}

func TestIsWildcard(t *testing.T) {
	yes, _ := IsWildCard("sockpuppet.org", "8.8.8.8:53")
	assert.False(t, yes)

	// XXX ToDo(erin): find a site with wildcard DNS for testing.
	//yes, _ = IsWildCard("sockpuppet.org", "8.8.8.8:53")
	//assert.True(t, yes)

}
