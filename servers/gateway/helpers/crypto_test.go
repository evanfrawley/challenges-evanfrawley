package helpers

import "testing"

func TestGetMD5Hash(t *testing.T) {
    cases := []struct{
        name string
        str string
        hash string
        expectedSame bool
    }{
        {
            name: "MD5 hash should work",
            str: "test@test.com",
            hash: "b642b4217b34b1e8d3bd915fc65c4452",
            expectedSame: true,
        },
        {
            name: "MD5 hash should work",
            str: "test@test.com",
            hash: "hashdoesntequal",
            expectedSame: false,
        },
    }

    for _, c := range cases {
        hash := GetMD5Hash(c.str)
        if hash != c.hash && c.expectedSame {
            t.Errorf("Name: %s\nExpected hash values to equal\nGot: %s | Expected: %s", c.name, hash, c.hash)
        }
        if hash == c.hash && !c.expectedSame {
            t.Errorf("Name: %s\nExpected hash values to not equal\nGot: %s | Expected: %s\nString that was hashed: %s", c.name, hash, c.hash, c.str)
        }
    }
}