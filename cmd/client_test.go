package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net/http"
	"os"
)

func (s *S) TestShouldReturnBodyMessageOnError(c *C) {
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, IsNil)

	client := NewClient(&http.Client{Transport: &transport{msg: "You must be authenticated to execute this command.", status: http.StatusUnauthorized}})
	response, err := client.Do(request)
	c.Assert(response, IsNil)
	c.Assert(err.Error(), Equals, "You must be authenticated to execute this command.")
}

func (s *S) TestShouldReturnErrorWhenServerIsDown(c *C) {
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, IsNil)
	client := NewClient(&http.Client{})
	_, err = client.Do(request)
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "^Server is down\n$")
}

func (s *S) TestShouldNotIncludeTheHeaderAuthorizationWhenTheTsuruTokenFileIsMissing(c *C) {
	os.Remove(os.ExpandEnv("${HOME}/.tsuru_token"))
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, IsNil)
	trans := &transport{msg: "", status: http.StatusOK}
	client := NewClient(&http.Client{Transport: trans})
	_, err = client.Do(request)
	c.Assert(err, IsNil)
	header := map[string][]string(request.Header)
	_, ok := header["Authorization"]
	c.Assert(ok, Equals, false)
}

func (s *S) TestShouldIncludeTheHeaderAuthorizationWhenTsuruTokenFileExists(c *C) {
	token_path := os.ExpandEnv("${HOME}/.tsuru_token")
	defer os.Remove(token_path)
	f, err := os.Create(token_path)
	c.Assert(err, IsNil)
	defer f.Close()
	token := []byte("mytoken")
	n, err := f.Write(token)
	c.Assert(err, IsNil)
	c.Assert(n, Equals, len(token))
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, IsNil)
	trans := &transport{msg: "", status: http.StatusOK}
	client := NewClient(&http.Client{Transport: trans})
	_, err = client.Do(request)
	c.Assert(err, IsNil)
	c.Assert(request.Header.Get("Authorization"), Equals, "mytoken")
}
