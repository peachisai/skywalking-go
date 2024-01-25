// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package goelasticsearchv7

import (
	"fmt"
	"github.com/labstack/echo/v4"

	"github.com/apache/skywalking-go/plugins/core/log"
	"github.com/apache/skywalking-go/plugins/core/operator"
	"github.com/apache/skywalking-go/plugins/core/tracing"
	es "github.com/elastic/go-elasticsearch/v7"
)

type EchoInterceptor struct{}

// BeforeInvoke would be called before the target method invocation.
func (h *EchoInterceptor) BeforeInvoke(invocation operator.Invocation) error {
	config := invocation.Args()[0].(es.Config)
	addresses := config.Addresses
	span, err := tracing.CreateExitSpan("testGolang", addresses[0], func(headerKey, headerValue string) error {
		return nil
	}, tracing.WithComponent(42),
		tracing.WithLayer(tracing.SpanLayerDatabase),
		tracing.WithTag(tracing.TagDBType, "Elasticsearch"))

	if err != nil {
		log.Warnf("cannot create exit span on elasticsearch client: %v", err)
		return nil
	}
	invocation.SetContext(span)
	return nil
}

// AfterInvoke would be called after the target method invocation.
func (h *EchoInterceptor) AfterInvoke(invocation operator.Invocation, result ...interface{}) error {
	span := invocation.GetContext().(tracing.Span)
	span.End()
	return nil
}

func middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			span, err := tracing.CreateEntrySpan(
				request.Method+":"+c.Path(),
				func(headerKey string) (string, error) {
					return request.Header.Get(headerKey), nil
				},
				tracing.WithLayer(tracing.SpanLayerHTTP),
				tracing.WithTag(tracing.TagHTTPMethod, request.Method),
				tracing.WithTag(tracing.TagURL, request.Host+request.URL.Path),
				tracing.WithComponent(5015))
			if err != nil {
				return err
			}

			// serve the request to the next middleware
			if err = next(c); err != nil {
				span.Error(err.Error())
				// invokes the registered HTTP error handler
				c.Error(err)
			}
			span.Tag(tracing.TagStatusCode, fmt.Sprintf("%d", c.Response().Status))
			span.End()
			return nil
		}
	}
}
