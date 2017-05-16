// This file is generated by protoc-gen-persist
// Source File: examples/sql/basic/basic.proto
// DO NOT EDIT !
package basic

import (
	sql "database/sql"
	fmt "fmt"
	io "io"
	strings "strings"

	hooks "github.com/tcncloud/protoc-gen-persist/examples/hooks"
	mytime "github.com/tcncloud/protoc-gen-persist/examples/mytime"
	test "github.com/tcncloud/protoc-gen-persist/examples/test"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
)

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// WARNING ! WARNING ! WARNING ! WARNING !WARNING ! WARNING !
// In order for your code to work you have to create a file
// in this package with the following content:
//
// type AmazingImpl struct {
// 	SqlDB *sql.DB
// }
// WARNING ! WARNING ! WARNING ! WARNING !WARNING ! WARNING !
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
func NewAmazingImpl(driver, connString string) (*AmazingImpl, error) {
	db, err := sql.Open(driver, connString)
	if err != nil {
		return nil, err
	}
	return &AmazingImpl{SqlDB: db}, nil
}

// sql unary UniarySelect
func (s *AmazingImpl) UniarySelect(ctx context.Context, req *test.PartialTable) (*test.ExampleTable, error) {
	var (
		Id        int64
		Name      string
		StartTime mytime.MyTime
		err       error
	)

	err = s.SqlDB.QueryRow("SELECT * from example_table Where id=$1 AND start_time>$2", req.Id, mytime.MyTime{}.ToSql(req.StartTime)).
		Scan(&Id,
			&StartTime,
			&Name,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "%+v doesn't exist", req)
		} else if strings.Contains(err.Error(), "duplicate key") {
			return nil, grpc.Errorf(codes.AlreadyExists, "%+v already exists", req)
		}
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	res := test.ExampleTable{

		Id:        Id,
		Name:      Name,
		StartTime: StartTime.ToProto(),
	}

	return &res, nil
}

// sql unary UniarySelectWithHooks
func (s *AmazingImpl) UniarySelectWithHooks(ctx context.Context, req *test.PartialTable) (*test.ExampleTable, error) {
	var (
		Id        int64
		Name      string
		StartTime mytime.MyTime
		err       error
	)

	beforeRes, err := hooks.UniarySelectBeforeHook(req)
	if err != nil {

		return nil, grpc.Errorf(codes.Unknown, err.Error())

	}
	if beforeRes != nil {

		return beforeRes, nil

	}

	err = s.SqlDB.QueryRow("SELECT * from example_table Where id=$1 AND start_time>$2", req.Id, mytime.MyTime{}.ToSql(req.StartTime)).
		Scan(&Id,
			&StartTime,
			&Name,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "%+v doesn't exist", req)
		} else if strings.Contains(err.Error(), "duplicate key") {
			return nil, grpc.Errorf(codes.AlreadyExists, "%+v already exists", req)
		}
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	res := test.ExampleTable{

		Id:        Id,
		Name:      Name,
		StartTime: StartTime.ToProto(),
	}

	err = hooks.UniarySelectAfterHook(req, &res)
	if err != nil {

		return nil, grpc.Errorf(codes.Unknown, err.Error())

	}

	return &res, nil
}

// sql server streaming ServerStream
func (s *AmazingImpl) ServerStream(req *test.Name, stream Amazing_ServerStreamServer) error {
	var (
		Id        int64
		Name      string
		StartTime mytime.MyTime
		err       error
	)

	rows, err := s.SqlDB.Query("SELECT * FROM example_table WHERE name=$1", req.Name)
	if err != nil {
		return grpc.Errorf(codes.Unknown, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			if err == sql.ErrNoRows {
				return grpc.Errorf(codes.NotFound, "%+v doesn't exist", req)
			} else if strings.Contains(err.Error(), "duplicate key") {
				return grpc.Errorf(codes.AlreadyExists, "%+v already exists", req)
			}
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		err := rows.Scan(&Id, &StartTime, &Name)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		res := test.ExampleTable{

			Id:        Id,
			Name:      Name,
			StartTime: StartTime.ToProto(),
		}

		stream.Send(&res)
	}
	return nil
}

// sql server streaming ServerStreamWithHooks
func (s *AmazingImpl) ServerStreamWithHooks(req *test.Name, stream Amazing_ServerStreamWithHooksServer) error {
	var (
		Id        int64
		Name      string
		StartTime mytime.MyTime
		err       error
	)

	beforeRes, err := hooks.ServerStreamBeforeHook(req)
	if err != nil {

		return grpc.Errorf(codes.Unknown, err.Error())

	}
	if beforeRes != nil {

		for _, res := range beforeRes {
			err = stream.Send(res)
			if err != nil {
				return err
			}
		}
		return nil

	}

	rows, err := s.SqlDB.Query("SELECT * FROM example_table WHERE name=$1", req.Name)
	if err != nil {
		return grpc.Errorf(codes.Unknown, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			if err == sql.ErrNoRows {
				return grpc.Errorf(codes.NotFound, "%+v doesn't exist", req)
			} else if strings.Contains(err.Error(), "duplicate key") {
				return grpc.Errorf(codes.AlreadyExists, "%+v already exists", req)
			}
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		err := rows.Scan(&Id, &StartTime, &Name)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		res := test.ExampleTable{

			Id:        Id,
			Name:      Name,
			StartTime: StartTime.ToProto(),
		}

		err = hooks.ServerStreamAfterHook(req, &res)
		if err != nil {

			return grpc.Errorf(codes.Unknown, err.Error())

		}

		stream.Send(&res)
	}
	return nil
}

// sql bidi streaming Bidirectional
func (s *AmazingImpl) Bidirectional(stream Amazing_BidirectionalServer) error {
	var err error
	stmt, err := s.SqlDB.Prepare("UPDATE example_table SET (start_time, name) = ($2, $3) WHERE id=$1 RETURNING *")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		var (
			Id        int64
			Name      string
			StartTime mytime.MyTime
		)
		err = stmt.QueryRow(req.Id, mytime.MyTime{}.ToSql(req.StartTime), req.Name).
			Scan(&Id, &StartTime, &Name)
		if err != nil {
			if err == sql.ErrNoRows {
				return grpc.Errorf(codes.NotFound, "%+v doesn't exist", req)
			} else if strings.Contains(err.Error(), "duplicate key") {
				return grpc.Errorf(codes.AlreadyExists, "%+v already exists", req)
			}
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		res := test.ExampleTable{

			Id:        Id,
			Name:      Name,
			StartTime: StartTime.ToProto(),
		}

		if err := stream.Send(&res); err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	return nil
}

// sql bidi streaming BidirectionalWithHooks
func (s *AmazingImpl) BidirectionalWithHooks(stream Amazing_BidirectionalWithHooksServer) error {
	var err error
	stmt, err := s.SqlDB.Prepare("UPDATE example_table SET (start_time, name) = ($2, $3) WHERE id=$1 RETURNING *")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		beforeRes, err := hooks.BidirectionalBeforeHook(req)
		if err != nil {

			return grpc.Errorf(codes.Unknown, err.Error())

		}
		if beforeRes != nil {

			err = stream.Send(beforeRes)
			if err != nil {
				return grpc.Errorf(codes.Unknown, err.Error())
			}
			continue

		}

		var (
			Id        int64
			Name      string
			StartTime mytime.MyTime
		)
		err = stmt.QueryRow(req.Id, mytime.MyTime{}.ToSql(req.StartTime), req.Name).
			Scan(&Id, &StartTime, &Name)
		if err != nil {
			if err == sql.ErrNoRows {
				return grpc.Errorf(codes.NotFound, "%+v doesn't exist", req)
			} else if strings.Contains(err.Error(), "duplicate key") {
				return grpc.Errorf(codes.AlreadyExists, "%+v already exists", req)
			}
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		res := test.ExampleTable{

			Id:        Id,
			Name:      Name,
			StartTime: StartTime.ToProto(),
		}

		err = hooks.BidirectionalAfterHook(req, &res)
		if err != nil {

			return grpc.Errorf(codes.Unknown, err.Error())

		}

		if err := stream.Send(&res); err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	return nil
}

// sql client streaming ClientStream
func (s *AmazingImpl) ClientStream(stream Amazing_ClientStreamServer) error {
	var err error
	tx, err := s.SqlDB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO example_table (id, start_time, name) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	res := test.NumRows{}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			tx.Rollback()
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		_, err = stmt.Exec(req.Id, mytime.MyTime{}.ToSql(req.StartTime), req.Name)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				return grpc.Errorf(codes.NotFound, "%+v doesn't exist", req)
			} else if strings.Contains(err.Error(), "duplicate key") {
				return grpc.Errorf(codes.AlreadyExists, "%+v already exists", req)
			}
			return grpc.Errorf(codes.Unknown, err.Error())
		}

	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("Commiting transaction failed, rolling back...")
		return grpc.Errorf(codes.Unknown, err.Error())
	}
	stream.SendAndClose(&res)
	return nil
}

// sql client streaming ClientStreamWithHook
func (s *AmazingImpl) ClientStreamWithHook(stream Amazing_ClientStreamWithHookServer) error {
	var err error
	tx, err := s.SqlDB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO example_table (id, start_time, name) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	res := test.Ids{}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			tx.Rollback()
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		beforeRes, err := hooks.ClientStreamBeforeHook(req)
		if err != nil {

			tx.Rollback()
			return grpc.Errorf(codes.Unknown, err.Error())

		}
		if beforeRes != nil {

			continue

		}

		_, err = stmt.Exec(req.Id, mytime.MyTime{}.ToSql(req.StartTime), req.Name)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				return grpc.Errorf(codes.NotFound, "%+v doesn't exist", req)
			} else if strings.Contains(err.Error(), "duplicate key") {
				return grpc.Errorf(codes.AlreadyExists, "%+v already exists", req)
			}
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		err = hooks.ClientStreamAfterHook(req, &res)
		if err != nil {

			tx.Rollback()
			return grpc.Errorf(codes.Unknown, err.Error())

		}

	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("Commiting transaction failed, rolling back...")
		return grpc.Errorf(codes.Unknown, err.Error())
	}
	stream.SendAndClose(&res)
	return nil
}
