// Copyright (c) 2022 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package paramobject

import (
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Client sends requests to a server.
type Client struct {
	url  string
	http *http.Client
	log  *zap.Logger
}

// ClientConfig defines the configuration for the client.
type ClientConfig struct {
	URL string
}

// ClientParams defines the parameters necessary to build a client.
// region empty
// region fxin
// region fields
type ClientParams struct {
	// endregion empty
	fx.In
	// endregion fxin

	Config     ClientConfig
	HTTPClient *http.Client
	// region empty
}

// endregion fields
// endregion empty

// NewClient builds a new client.
// region takeparam
// region consume
func NewClient(p ClientParams) (*Client, error) {
	// endregion takeparam
	return &Client{
		url:  p.Config.URL,
		http: p.HTTPClient,
		// ...
	}, nil
	// endregion consume
}
