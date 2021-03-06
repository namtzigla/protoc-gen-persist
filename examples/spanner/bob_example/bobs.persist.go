// This file is generated by protoc-gen-persist
// Source File: examples/spanner/bob_example/bobs.proto
// DO NOT EDIT !
package bob_example

import (
	io "io"
	strings "strings"

	"cloud.google.com/go/spanner"
	mytime "github.com/tcncloud/protoc-gen-persist/examples/mytime"
	pb "github.com/tcncloud/protoc-gen-persist/examples/spanner/bob_example"
	context "golang.org/x/net/context"
	iterator "google.golang.org/api/iterator"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
)

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// WARNING ! WARNING ! WARNING ! WARNING !WARNING ! WARNING !
// In order for your code to work you have to create a file
// in this package with the following content:
//
// type BobsImpl struct {
// 	SpannerDB *spanner.Client
// }
// WARNING ! WARNING ! WARNING ! WARNING !WARNING ! WARNING !
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// spanner unary select DeleteBobs
func (s *BobsImpl) DeleteBobs(ctx context.Context, req *pb.Bob) (*pb.Empty, error) {
	var err error

	start := make([]interface{}, 0)
	end := make([]interface{}, 0)
	var conv interface{}
	start = append(start, "Bob")
	end = append(end, "Bob")
	conv, err = mytime.MyTime{}.ToSpanner(req.StartTime).SpannerValue()
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	end = append(end, conv)
	key := spanner.KeyRange{
		Start: start,
		End:   end,
		Kind:  spanner.ClosedOpen,
	}
	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.Delete("bob_table", key)
	_, err = s.SpannerDB.Apply(ctx, muts)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, grpc.Errorf(codes.NotFound, err.Error())
		}
	}
	res := pb.Empty{}

	return &res, nil
}

// spanner client streaming PutBobs
func (s *BobsImpl) PutBobs(stream pb.Bobs_PutBobsServer) error {
	var err error
	res := pb.NumRows{}

	muts := make([]*spanner.Mutation, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		//spanner client streaming insert
		params := make([]interface{}, 0)
		var conv interface{}

		conv = req.Id

		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		params = append(params, conv)

		conv = req.Name

		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		params = append(params, conv)

		conv, err = mytime.MyTime{}.ToSpanner(req.StartTime).SpannerValue()

		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		params = append(params, conv)
		muts = append(muts, spanner.Insert("bob_table", []string{"id", "name", "start_time"}, params))

		////////////////////////////// NOTE //////////////////////////////////////
		// In the future, we might do apply if muts gets really big,  but for now,
		// we only do one apply on the database with all the records stored in muts
		//////////////////////////////////////////////////////////////////////////
	}
	_, err = s.SpannerDB.Apply(context.Background(), muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
	}

	stream.SendAndClose(&res)
	return nil
}

// spanner server streaming GetBobs
func (s *BobsImpl) GetBobs(req *pb.Empty, stream pb.Bobs_GetBobsServer) error {
	var (
		Id        int64
		Name      string
		StartTime mytime.MyTime
	)

	params := make(map[string]interface{})
	stmt := spanner.Statement{SQL: "SELECT * from bob_table", Params: params}
	tx := s.SpannerDB.Single()
	defer tx.Close()
	iter := tx.Query(context.Background(), stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		// scan our values out of the row

		err = row.ColumnByName("id", &Id)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		gcv1 := new(spanner.GenericColumnValue)
		err = row.ColumnByName("start_time", gcv1)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		err = StartTime.SpannerScan(gcv1)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		err = row.ColumnByName("name", &Name)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		res := pb.Bob{

			Id:        Id,
			Name:      Name,
			StartTime: StartTime.ToProto(),
		}

		stream.Send(&res)
	}
	return nil
}

// spanner server streaming GetPeopleFromNames
func (s *BobsImpl) GetPeopleFromNames(req *pb.Names, stream pb.Bobs_GetPeopleFromNamesServer) error {
	var (
		Id        int64
		Name      string
		StartTime mytime.MyTime
	)

	var err error

	params := make(map[string]interface{})
	var conv interface{}

	conv = req.Names

	if err != nil {
		return grpc.Errorf(codes.Unknown, err.Error())
	}
	params["string0"] = conv
	stmt := spanner.Statement{SQL: "SELECT * FROM bob_table WHERE name IN UNNEST( @string0)", Params: params}
	tx := s.SpannerDB.Single()
	defer tx.Close()
	iter := tx.Query(context.Background(), stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		// scan our values out of the row

		err = row.ColumnByName("id", &Id)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		gcv1 := new(spanner.GenericColumnValue)
		err = row.ColumnByName("start_time", gcv1)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		err = StartTime.SpannerScan(gcv1)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		err = row.ColumnByName("name", &Name)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		res := pb.Bob{

			Id:        Id,
			Name:      Name,
			StartTime: StartTime.ToProto(),
		}

		stream.Send(&res)
	}
	return nil
}
