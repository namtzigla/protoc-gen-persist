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

package generator

import (
	"bytes"
	"strconv"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/tcncloud/protoc-gen-persist/generator/templates"
)

var (
	fileTemplate *template.Template
	TemplateList = map[string]string{
		"import_template":                 templates.ImportTemplate,
		"implement_structs":               templates.StructsTemplate,
		"implement_services":              templates.ServicesTemplate,
		"implement_method":                templates.MethodTemplate,
		"return_convert_helpers":          templates.ReturnConvertHelpers,
		"before_hook":                     templates.BeforeHook,
		"after_hook":                      templates.AfterHook,
		"type_desc_to_mapped":             templates.SpannerHelperTemplates,
		"unary_method":                    templates.UnaryMethodTemplate,
		"client_streaming_method":         templates.ClientStreamingMethodTemplate,
		"server_streaming_method":         templates.ServerStreamingMethodTemplate,
		"bidi_method":                     templates.BidiStreamingMethodTemplate,
		"sql_unary_method":                templates.SqlUnaryMethodTemplate,
		"sql_client_streaming_method":     templates.SqlClientStreamingMethodTemplate,
		"sql_server_streaming_method":     templates.SqlServerStreamingMethodTemplate,
		"sql_bidi_streaming_method":       templates.SqlBidiStreamingMethodTemplate,
		"mongo_unary_method":              templates.MongoUnaryMethodTemplate,
		"mongo_client_streaming_method":   templates.MongoClientStreamingMethodTemplate,
		"mongo_server_streaming_method":   templates.MongoServerStreamingMethodTemplate,
		"mongo_bidi_streaming_method":     templates.MongoBidiStreamingMethodTemplate,
		"spanner_unary_method":            templates.SpannerUnaryMethodTemplate,
		"spanner_unary_delete":            templates.SpannerUnaryDeleteTemplate,
		"spanner_unary_update":            templates.SpannerUnaryUpdateTemplate,
		"spanner_unary_select":            templates.SpannerUnarySelectTemplate,
		"spanner_unary_insert":            templates.SpannerUnaryInsertTemplate,
		"spanner_client_streaming_method": templates.SpannerClientStreamingMethodTemplate,
		"spanner_client_streaming_update": templates.SpannerClientStreamingUpdateTemplate,
		"spanner_client_streaming_insert": templates.SpannerClientStreamingInsertTemplate,
		"spanner_client_streaming_delete": templates.SpannerClientStreamingDeleteTemplate,
		"spanner_server_streaming_method": templates.SpannerServerStreamingMethodTemplate,
		"spanner_bidi_streaming_method":   templates.SpannerBidiStreamingMethodTemplate,
	}
)

func init() {
	logrus.Debug("files package init()")
	var err error
	fileTemplate, err = template.New("fileTemplate").Parse(templates.MainTemplate)
	if err != nil {
		logrus.WithError(err).Fatal("Fail to parse file template")
	}
	fileTemplate := fileTemplate.Funcs(template.FuncMap{
		"Quotes": strconv.Quote,
	})

	for n, tmpl := range TemplateList {
		_, err = fileTemplate.Parse(tmpl)
		if err != nil {
			logrus.WithError(err).Fatalf("Fatal error parsing template %s", n)
		}
	}
}

func ExecuteFileTemplate(fileStruct *FileStruct) []byte {
	var buffer bytes.Buffer
	err := fileTemplate.Execute(&buffer, fileStruct)
	if err != nil {
		logrus.WithError(err).Fatal("Fatal error executing file template")
	}
	return buffer.Bytes()
}
