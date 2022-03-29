package grpc

import (
	"context"
	"fmt"
	"github.com/norby7/shortening-service/entities"
	"github.com/norby7/shortening-service/interfaceAdapters/grpc/protocol"
	"github.com/norby7/shortening-service/usecases/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"os"
	"testing"
)

const bufSize = 1024 * 1024

var (
	lis          *bufconn.Listener
	createError  = fmt.Errorf("unable to create the url")
	deleteError  = fmt.Errorf("unable to delete the url")
	getError     = fmt.Errorf("unable to fetch the url")
	counterError = fmt.Errorf("unable to increment counter")
)

type ServiceMock struct{}

func (s *ServiceMock) Create(u *entities.Url) error {
	if u.Url == "http://www.invalidUrl.com" || u.Url == "" {
		return createError
	}

	if u.Code == "d4jn8dsf" {
		return service.ErrCodeAlreadyExists
	}

	return nil
}

func (s *ServiceMock) Delete(id int64) error {
	if id == 0 {
		return deleteError
	}

	return nil
}

func (s *ServiceMock) GetUrlByCode(code string) (string, error) {
	if code == "invalidCode" {
		return "", getError
	}

	if code != "84gfj4i9" {
		return "", nil
	}

	return "https://google.com", nil
}

func (s *ServiceMock) GetById(id int64) (entities.Url, error) {
	if id == 0 {
		return entities.Url{}, getError
	}

	if id != 1 {
		return entities.Url{}, nil
	}

	return entities.Url{
		Id:       1,
		Code:     "84gfj4i9",
		Url:      "https://google.com",
		ShortUrl: "http://localhost/84gfj4i9",
		Domain:   "http://localhost",
		Counter:  1,
	}, nil
}

func (s *ServiceMock) IncrementCounter(string) {

}

func init() {
	serviceMock := ServiceMock{}
	l := log.New(os.Stdout, "urls-api", log.LstdFlags)
	urlService := NewUrlGrpcService(&serviceMock, l)

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	protocol.RegisterUrlServiceServer(s, urlService)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func TestAdd(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			t.Errorf(err.Error())
		}
	}()

	client := protocol.NewUrlServiceClient(conn)

	testCases := []struct {
		name          string
		input         protocol.Url
		expectedError bool
	}{
		{
			name:          "empty url",
			input:         protocol.Url{},
			expectedError: true,
		},
		{
			name: "create service error",
			input: protocol.Url{
				Id:       1,
				Code:     "84gfj4i9",
				Url:      "http://www.invalidUrl.com",
				ShortUrl: "http://localhost/84gfj4i9",
				Domain:   "http://localhost",
				Counter:  1,
			},
			expectedError: true,
		},
		{
			name: "valid request",
			input: protocol.Url{
				Id:       1,
				Code:     "84gfj4i9",
				Url:      "https://google.com",
				ShortUrl: "http://localhost/84gfj4i9",
				Domain:   "http://localhost",
				Counter:  1,
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Add(ctx, &tc.input)

			if (err != nil) != tc.expectedError {
				t.Errorf("expected error (%v), got (%v) with response: (%v)", tc.expectedError, err, resp.String())
			}
		})
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			t.Errorf(err.Error())
		}
	}()

	client := protocol.NewUrlServiceClient(conn)

	testCases := []struct {
		name          string
		input         protocol.UrlId
		expectedError bool
	}{
		{
			name:          "empty url id",
			input:         protocol.UrlId{},
			expectedError: true,
		},
		{
			name:          "delete service error",
			input:         protocol.UrlId{Value: 0},
			expectedError: true,
		},
		{
			name:          "valid request",
			input:         protocol.UrlId{Value: 1},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Delete(ctx, &tc.input)

			if (err != nil) != tc.expectedError {
				t.Errorf("expected error (%v), got (%v) with response: (%v)", tc.expectedError, err, resp.String())
			}
		})
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			t.Errorf(err.Error())
		}
	}()

	client := protocol.NewUrlServiceClient(conn)

	testCases := []struct {
		name          string
		input         protocol.UrlId
		expectedError bool
	}{
		{
			name:          "empty url id",
			input:         protocol.UrlId{},
			expectedError: true,
		},
		{
			name:          "get service error",
			input:         protocol.UrlId{Value: 0},
			expectedError: true,
		},
		{
			name:          "valid request",
			input:         protocol.UrlId{Value: 1},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Get(ctx, &tc.input)

			if (err != nil) != tc.expectedError {
				t.Errorf("expected error (%v), got (%v) with response: (%v)", tc.expectedError, err, resp.String())
			}
		})
	}
}

func TestGetCounter(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			t.Errorf(err.Error())
		}
	}()

	client := protocol.NewUrlServiceClient(conn)

	testCases := []struct {
		name          string
		input         protocol.UrlId
		expectedError bool
	}{
		{
			name:          "empty url id",
			input:         protocol.UrlId{},
			expectedError: true,
		},
		{
			name:          "get service error",
			input:         protocol.UrlId{Value: 0},
			expectedError: true,
		},
		{
			name:          "valid request",
			input:         protocol.UrlId{Value: 1},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.GetCounter(ctx, &tc.input)

			if (err != nil) != tc.expectedError {
				t.Errorf("expected error (%v), got (%v) with response: (%v)", tc.expectedError, err, resp.String())
			}
		})
	}
}
