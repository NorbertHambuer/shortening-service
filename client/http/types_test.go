package http

import "testing"

func TestValidateUrl(t *testing.T){
	testCases := []struct{
		name string
		input Url
		isError bool
	}{
		{
			name:    "valid url",
			input:   Url{
				Id:       1,
				Code:     "84gfj4i9",
				Url:      "https://google.com",
				ShortUrl: "http://localhost/84gfj4i9",
				Domain:   "http://localhost",
				Counter:  1,
			},
			isError: false,
		},
		{
			name:    "empty url",
			input:   Url{},
			isError: true,
		},
		{
			name:    "invalid code length",
			input:   Url{
				Code:    "84gfj4i9gdfg34g34",
				Url:      "https://google.com",
				ShortUrl: "http://localhost/84gfj4i9",
				Domain:   "http://localhost",
				Counter:  1,
			},
			isError: true,
		},
		{
			name:    "invalid url length",
			input:   Url{
				Code:    "84gfj4i9",
				Url:      "ole.com",
				ShortUrl: "http://localhost/84gfj4i9",
				Domain:   "http://localhost",
				Counter:  1,
			},
			isError: true,
		},
		{
			name:    "invalid url structure pattern",
			input:   Url{
				Code:    "84gfj4i9",
				Url:      "postgres://user:abc{DEf1=ghi@example.com:5432/db?sslmode=require",
				ShortUrl: "http://localhost/84gfj4i9",
				Domain:   "http://localhost",
				Counter:  1,
			},
			isError: true,
		},
		{
			name:    "invalid shortUrl structure pattern",
			input:   Url{
				Code:    "84gfj4i9",
				Url:      "www.google.com",
				ShortUrl: "postgres://user:abc{DEf1=ghi@example.com:5432/db?sslmode=require",
				Domain:   "http://localhost",
				Counter:  1,
			},
			isError: true,
		},
	}

	for _, tc := range testCases{
		t.Run(tc.name, func(t *testing.T){
			err := tc.input.Validate()

			if (err != nil) != tc.isError{
				t.Errorf("expected error (%v), got (%v)", tc.isError, err)
			}
		})
	}
}
