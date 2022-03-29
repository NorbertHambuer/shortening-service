package cache

import (
	"github.com/alicebob/miniredis"
	"log"
	"strings"
	"testing"
)

func TestValidNewRedisCache(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	srvAddr := strings.Split(mr.Addr(), ":")

	_, err = NewRedisCache(srvAddr[0], srvAddr[1], "")
	if err != nil {
		t.Errorf("unable to connect to miniredis server: %s", err.Error())
	}
}

func TestInvalidNewRedisCache(t *testing.T) {
	_, err := NewRedisCache("localhost", "1", "")
	if err == nil {
		t.Errorf("expected error connecting to redis server: %s", err.Error())
	}
}

func TestSetGetShortUrl(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	srvAddr := strings.Split(mr.Addr(), ":")

	client, err := NewRedisCache(srvAddr[0], srvAddr[1], "")
	if err != nil {
		t.Errorf("unable to connect to miniredis server: %s", err.Error())
	}

	err = client.SetShortUrl("test", "www.test.com")
	if err != nil {
		t.Errorf("unable to set short url: %s", err.Error())
	}

	_, err = client.GetShortUrl("test")
	if err != nil {
		t.Errorf("unable to get short url: %s", err.Error())
	}
}

func TestGetShortUrlError(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	srvAddr := strings.Split(mr.Addr(), ":")

	client, err := NewRedisCache(srvAddr[0], srvAddr[1], "")
	if err != nil {
		t.Errorf("unable to connect to miniredis server: %s", err.Error())
	}

	_, err = client.GetShortUrl("test1")
	if err == nil {
		t.Errorf("expected error getting short url, got nil")
	}
}

func TestCacheDisabled(t *testing.T) {
	redisCache, err := NewRedisCache("", "", "")
	if err == nil {
		t.Error("expected error connecting to redis cache, got nil")
	}

	err = redisCache.SetShortUrl("code", "url")
	if err != nil {
		t.Errorf("expected no error, got (%s)", err.Error())
	}

	url, err := redisCache.GetShortUrl("code")
	if url != "" {
		t.Errorf("expected empty url, got (%s)", url)
	}

	if err != nil {
		t.Errorf("expected no error, got (%s)", err.Error())
	}
}
