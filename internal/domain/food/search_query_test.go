package food

import "testing"

// TestSanitizeFulltextQuery - operator BOOLEAN MODE harus dibuang.
//
// Regresi: input seperti `ayam "goreng` dulu diteruskan apa adanya ke
// AGAINST(... IN BOOLEAN MODE), MySQL menolaknya dengan error sintaks, dan
// endpoint pencarian makanan balas 500.
func TestSanitizeFulltextQuery(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"nasi", "nasi*"},
		{"nasi goreng", "nasi* goreng*"},
		{`ayam "goreng`, "ayam* goreng*"},
		{"nasi +", "nasi*"},
		{"+++", ""},
		{`@ > < ( ) ~`, ""},
		{"  ayam   bakar  ", "ayam* bakar*"},
	}

	for _, c := range cases {
		if got := sanitizeFulltextQuery(c.in); got != c.want {
			t.Errorf("sanitizeFulltextQuery(%q) = %q, mau %q", c.in, got, c.want)
		}
	}
}

// TestEscapeLikePattern - wildcard LIKE dari user harus jadi teks biasa,
// bukan pola yang mencocokkan semua baris.
func TestEscapeLikePattern(t *testing.T) {
	if got := escapeLikePattern("%"); got != `%\%%` {
		t.Errorf("escapeLikePattern(%%) = %q", got)
	}
	if got := escapeLikePattern("na_si"); got != `%na\_si%` {
		t.Errorf("escapeLikePattern(na_si) = %q", got)
	}
	if got := escapeLikePattern("nasi"); got != "%nasi%" {
		t.Errorf("escapeLikePattern(nasi) = %q", got)
	}
}
