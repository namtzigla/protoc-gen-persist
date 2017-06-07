// This file is generated by protoc-gen-persist
// Source File: examples/spanner/delete/delete_examples.proto
// DO NOT EDIT !
package delete

import (
	strings "strings"

	"cloud.google.com/go/spanner"
	context "golang.org/x/net/context"
	"google.golang.org/api/option"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
)

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// WARNING ! WARNING ! WARNING ! WARNING !WARNING ! WARNING !
// In order for your code to work you have to create a file
// in this package with the following content:
//
// type DeleteImpl struct {
// 	SpannerDB *spanner.Client
// }
// WARNING ! WARNING ! WARNING ! WARNING !WARNING ! WARNING !
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
func NewDeleteImpl(d string, conf *spanner.ClientConfig, opts ...option.ClientOption) (*DeleteImpl, error) {
	var client *spanner.Client
	var err error
	if conf != nil {
		client, err = spanner.NewClientWithConfig(context.Background(), d, *conf, opts...)
	} else {
		client, err = spanner.NewClient(context.Background(), d, opts...)
	}
	if err != nil {
		return nil, err
	}
	return &DeleteImpl{SpannerDB: client}, nil
}

// spanner unary select DeleteEquals
func (s *DeleteImpl) DeleteEquals(ctx context.Context, req *Empty) (*Empty, error) {
	var err error

	start := make([]interface{}, 0)
	end := make([]interface{}, 0)
	start = append(start, 1)
	end = append(end, 1)
	key := spanner.KeyRange{
		Start: start,
		End:   end,
		Kind:  spanner.ClosedClosed,
	}
	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.Delete("mytable", key)
	_, err = s.SpannerDB.Apply(ctx, muts)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, grpc.Errorf(codes.NotFound, err.Error())
		}
	}
	res := Empty{}

	return &res, nil
}

// spanner unary select DeleteGreater
func (s *DeleteImpl) DeleteGreater(ctx context.Context, req *Empty) (*Empty, error) {
	var err error

	start := make([]interface{}, 0)
	end := make([]interface{}, 0)
	start = append(start, 1)
	key := spanner.KeyRange{
		Start: start,
		End:   end,
		Kind:  spanner.ClosedOpen,
	}
	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.Delete("mytable", key)
	_, err = s.SpannerDB.Apply(ctx, muts)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, grpc.Errorf(codes.NotFound, err.Error())
		}
	}
	res := Empty{}

	return &res, nil
}

// spanner unary select DeleteLess
func (s *DeleteImpl) DeleteLess(ctx context.Context, req *Empty) (*Empty, error) {
	var err error

	start := make([]interface{}, 0)
	end := make([]interface{}, 0)
	start = append(start, "b123b_asdf")
	start = append(start, 1.1)
	end = append(end, 10)
	key := spanner.KeyRange{
		Start: start,
		End:   end,
		Kind:  spanner.OpenClosed,
	}
	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.Delete("mytable", key)
	_, err = s.SpannerDB.Apply(ctx, muts)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, grpc.Errorf(codes.NotFound, err.Error())
		}
	}
	res := Empty{}

	return &res, nil
}

// spanner unary select DeleteMultiEquals
func (s *DeleteImpl) DeleteMultiEquals(ctx context.Context, req *Empty) (*Empty, error) {
	var err error

	start := make([]interface{}, 0)
	end := make([]interface{}, 0)
	start = append(start, 1)
	start = append(start, 50)
	end = append(end, 1)
	end = append(end, 50)
	key := spanner.KeyRange{
		Start: start,
		End:   end,
		Kind:  spanner.ClosedClosed,
	}
	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.Delete("mytable", key)
	_, err = s.SpannerDB.Apply(ctx, muts)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, grpc.Errorf(codes.NotFound, err.Error())
		}
	}
	res := Empty{}

	return &res, nil
}
