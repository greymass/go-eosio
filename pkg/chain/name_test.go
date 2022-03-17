package chain_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/greymass/go-eosio/internal/assert"
	"github.com/greymass/go-eosio/pkg/chain"
)

func TestName(t *testing.T) {
	assert.Equal(t, chain.NewName(""), chain.Name(0))
	assert.Equal(t, chain.NewName("").String(), "")
	assert.Equal(t, chain.Name(14595364149838066048).String(), "teamgreymass")
	assert.Equal(t, chain.NewName("teamgreymass"), chain.Name(14595364149838066048))
	assert.Equal(t, chain.NewName("invålid").String(), "inv..lid")
	assert.Equal(t, chain.NewName("overflowmyname").String(), "overflowmyna2")

	assert.JSONCoding(t, chain.Name(14595364149838066048), `"teamgreymass"`)
	assert.ABICoding(t, chain.Name(14595364149838066048), []byte{0x80, 0xb1, 0x91, 0x5e, 0x5d, 0x26, 0x8d, 0xca})
}

func FuzzName(f *testing.F) {
	f.Add("")
	f.Add(" ")
	f.Add("teamgreymass")
	f.Add("invålid")
	f.Add("overflowmyname")
	f.Add("f.u.z.z.")
	re := regexp.MustCompile("^[a-z1-5.]{0,13}$")
	f.Fuzz(func(t *testing.T, orig string) {
		name := chain.NewName(orig)
		actual := name.String()
		t.Logf("%q -> %q (%d)", orig, actual, name)
		if !re.MatchString(actual) || strings.HasSuffix(actual, ".") {
			t.Errorf("%q is not a valid name", orig)
		}
	})
}

func BenchmarkNameFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		chain.NewName("teamgreymass")
	}
}

func BenchmarkStringFromName(b *testing.B) {
	n1 := chain.NewName("teamgreymass")
	for i := 0; i < b.N; i++ {
		_ = n1.String()
	}
}
