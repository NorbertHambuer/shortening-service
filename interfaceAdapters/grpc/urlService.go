package grpc

import (
	"context"
	"github.com/norby7/shortening-service/entities"
	"github.com/norby7/shortening-service/interfaceAdapters/grpc/protocol"
	"github.com/norby7/shortening-service/usecases/service"
	"log"
)

type UrlGrpcService struct {
	Service service.Interactor
	Logger  *log.Logger
	protocol.UnimplementedUrlServiceServer
}

// NewUrlGrpcService returns a new UrlGrpcService object address
func NewUrlGrpcService(s service.Interactor, l *log.Logger) *UrlGrpcService {
	return &UrlGrpcService{Service: s, Logger: l}
}

// Add creates a new Url and inserts it into the database
func (us *UrlGrpcService) Add(ctx context.Context, u *protocol.Url) (*protocol.Url, error) {
	us.Logger.Println("UrlGrpcService:Add called")

	url := ProtoUrlToUrl(u)

	err := us.Service.Create(url)
	if err != nil {
		return &protocol.Url{}, err
	}

	return UrlToProtoUrl(url), nil
}

// Delete removes a url from the database based on the given ID
func (us *UrlGrpcService) Delete(ctx context.Context, id *protocol.UrlId) (*protocol.VoidResponse, error) {
	us.Logger.Println("UrlGrpcService:Delete called")

	err := us.Service.Delete(id.Value)
	if err != nil {
		return &protocol.VoidResponse{}, err
	}

	return &protocol.VoidResponse{}, nil
}

// Get returns a url from the database based on the given ID
func (us *UrlGrpcService) Get(ctx context.Context, id *protocol.UrlId) (*protocol.Url, error) {
	us.Logger.Println("UrlGrpcService:Get called")

	u, err := us.Service.GetById(id.Value)
	if err != nil {
		return &protocol.Url{}, err
	}

	return UrlToProtoUrl(&u), nil
}

// GetCounter returns the redirections counter for the given ID
func (us *UrlGrpcService) GetCounter(ctx context.Context, id *protocol.UrlId) (*protocol.Counter, error) {
	us.Logger.Println("UrlGrpcService:GetCounter called")

	u, err := us.Service.GetById(id.Value)
	if err != nil {
		return &protocol.Counter{}, err
	}

	return &protocol.Counter{Value: u.Counter}, nil
}

// ProtoUrlToUrl converts a *protocol.Url object into a *entities.Url object
func ProtoUrlToUrl(u *protocol.Url) *entities.Url {
	return &entities.Url{
		Id:       u.Id,
		Code:     u.Code,
		Url:      u.Url,
		ShortUrl: u.ShortUrl,
		Domain:   u.Domain,
		Counter:  u.Counter,
	}
}

// UrlToProtoUrl converts a *entities.Url object into a *protocol.Url object
func UrlToProtoUrl(u *entities.Url) *protocol.Url {
	return &protocol.Url{
		Id:       u.Id,
		Code:     u.Code,
		Url:      u.Url,
		ShortUrl: u.ShortUrl,
		Domain:   u.Domain,
		Counter:  u.Counter,
	}
}
