// Copyright 2017, TCN Inc.
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of TCN Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package templates

const SqlUnaryMethodTemplate = `{{define "sql_unary_method"}}// sql unary {{.GetName}}
func (s* {{.GetServiceName}}Impl) {{.GetName}} (ctx context.Context, req *{{.GetInputType}}) (*{{.GetOutputType}}, error) {
	var (
		{{range $field, $type := .GetFieldsWithLocalTypesFor .GetOutputTypeStruct}}
		{{$field}} {{$type}}{{end}}
		err error
	)

	{{template "before_hook" .}}

	err = s.SqlDB.QueryRow({{Quotes .GetQuery}} {{.GetQueryParamString true}}).
		Scan({{range $index,$t :=.GetTypeDescArrayForStruct .GetOutputTypeStruct}} &{{$t.Name}},
		{{end}})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "%+v doesn't exist", req)
		} else if strings.Contains(err.Error(), "duplicate key") {
			return nil, grpc.Errorf(codes.AlreadyExists, "%+v already exists", req)
		}
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	res := {{.GetOutputType}}{
	{{range $field, $type := .GetTypeDescForFieldsInStruct .GetOutputTypeStruct}}
	{{$field}}: {{template "addr" $type}}{{template "base" $type}}{{template "mapping" $type}},{{end}}
	}

	{{template "after_hook" .}}
	return &res, nil
}
{{end}}`

const SqlServerStreamingMethodTemplate = `{{define "sql_server_streaming_method"}}// sql server streaming {{.GetName}}
func (s *{{.GetServiceName}}Impl) {{.GetName}}(req *{{.GetInputType}}, stream {{.GetFilePackage}}{{.GetServiceName}}_{{.GetName}}Server) error {
	var (
 {{range $field, $type := .GetFieldsWithLocalTypesFor .GetOutputTypeStruct}}
 {{$field}} {{$type}}{{end}}
 	err error
 	)

	{{template "before_hook" .}}

	rows, err := s.SqlDB.Query({{Quotes .GetQuery}} {{.GetQueryParamString true}})

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

		err := rows.Scan({{range $index,$t :=.GetTypeDescArrayForStruct .GetOutputTypeStruct}} &{{$t.Name}},{{end}})
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		res := {{.GetOutputType}}{
		{{range $field, $type := .GetTypeDescForFieldsInStruct  .GetOutputTypeStruct}}
		{{$field}}: {{template "addr" $type}}{{template "base" $type}}{{template "mapping" $type}},{{end}}
		}

		{{template "after_hook" .}}

		stream.Send(&res)
	}
	return nil
}{{end}}`

const SqlClientStreamingMethodTemplate = `{{define "sql_client_streaming_method"}}// sql client streaming {{.GetName}}
func (s *{{.GetServiceName}}Impl) {{.GetName}}(stream {{.GetFilePackage}}{{.GetServiceName}}_{{.GetName}}Server) error {
	var err error
	tx, err := s.SqlDB.Begin()
	if err != nil {
		return err
	}
	stmt, err:= tx.Prepare({{Quotes .GetQuery}})
	if err != nil {
		return err
	}
	{{/* setup the response for aggregation in the after hook */}}
	res := {{.GetOutputType}}{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			tx.Rollback()
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		{{template "before_hook" .}}

		_, err = stmt.Exec({{.GetQueryParamString false}})
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				return grpc.Errorf(codes.NotFound, "%+v doesn't exist", req)
			} else if strings.Contains(err.Error(), "duplicate key") {
				return grpc.Errorf(codes.AlreadyExists, "%+v already exists", req)
			}
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		{{template "after_hook" .}}
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("Commiting transaction failed, rolling back...")
		return grpc.Errorf(codes.Unknown, err.Error())
	}
	stream.SendAndClose(&res)
	return nil
}{{end}}`

const SqlBidiStreamingMethodTemplate = `{{define "sql_bidi_streaming_method"}}// sql bidi streaming {{.GetName}}
func (s *{{.GetServiceName}}Impl) {{.GetName}}(stream {{.GetFilePackage}}{{.GetServiceName}}_{{.GetName}}Server) error {
	var err error
	stmt, err := s.SqlDB.Prepare({{Quotes .GetQuery}})
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
		{{template "before_hook" .}}
		var (
		 {{range $field, $type := .GetFieldsWithLocalTypesFor .GetOutputTypeStruct}}
		 {{$field}} {{$type}}{{end}}
		)
		err = stmt.QueryRow({{.GetQueryParamString false}}).
			Scan({{range $index,$t :=.GetTypeDescArrayForStruct .GetOutputTypeStruct}} &{{$t.Name}},{{end}})
		if err != nil {
			if err == sql.ErrNoRows {
				return grpc.Errorf(codes.NotFound, "%+v doesn't exist", req)
			} else if strings.Contains(err.Error(), "duplicate key") {
				return grpc.Errorf(codes.AlreadyExists, "%+v already exists", req)
			}
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		res := {{.GetOutputType}}{
		{{range $field, $type := .GetTypeDescForFieldsInStruct .GetOutputTypeStruct}}
		{{$field}}: {{template "addr" $type}}{{template "base" $type}}{{template "mapping" $type}},{{end}}
		}

		{{template "after_hook" .}}

		if err := stream.Send(&res); err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	return nil
}

{{end}}`
