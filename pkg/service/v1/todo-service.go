package v1

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stelo/blackmore/pkg/api/v1"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// linksellerServiceServer is implementation of v1.LinksellerServiceServer proto interface
type linksellerServiceServer struct {
	db *sql.DB
}

// NewLinksellerServiceServer creates linkseller service
func NewLinksellerServiceServer(db *sql.DB) v1.LinksellerServiceServer {
	return &linksellerServiceServer{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *linksellerServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// connect returns SQL database connection from the pool
func (s *linksellerServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// Create new linkseller task
func (s *linksellerServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// insert linkseller entity data
	res, err := c.ExecContext(ctx, "INSERT INTO linkseller(`Title`, `Description`, `Reminder`) VALUES(?, ?, ?)",
		req.Linkseller.Person.Type, req.Linkseller.Person.Document, nil)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into linkseller-> "+err.Error())
	}

	// get ID of creates linkseller
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created linkseller-> "+err.Error())
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}
