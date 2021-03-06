// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package ingress

import (
	"fmt"
	"net/http"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/apache/apisix-ingress-controller/test/e2e/scaffold"
	"github.com/onsi/ginkgo"
)

var _ = ginkgo.Describe("support ingress.networking/v1", func() {
	s := scaffold.NewDefaultScaffold()

	ginkgo.It("path exact match", func() {
		backendSvc, backendPort := s.DefaultHTTPBackend()
		ing := fmt.Sprintf(`
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
  name: ingress-v1
spec:
  rules:
  - host: httpbin.org
    http:
      paths:
      - path: /ip
        pathType: Exact
        backend:
          service:
            name: %s
            port:
              number: %d
`, backendSvc, backendPort[0])
		err := s.CreateResourceFromString(ing)
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")
		time.Sleep(5 * time.Second)

		_ = s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").Expect().Status(http.StatusOK)
		// Exact path, doesn't match /ip/aha
		_ = s.NewAPISIXClient().GET("/ip/aha").WithHeader("Host", "httpbin.org").Expect().Status(http.StatusNotFound)
		// Mismatched host
		_ = s.NewAPISIXClient().GET("/ip/aha").WithHeader("Host", "a.httpbin.org").Expect().Status(http.StatusNotFound)
	})

	ginkgo.It("path prefix match", func() {
		backendSvc, backendPort := s.DefaultHTTPBackend()
		ing := fmt.Sprintf(`
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
  name: ingress-v1
spec:
  rules:
  - host: httpbin.org
    http:
      paths:
      - path: /status
        pathType: Prefix
        backend:
          service:
            name: %s
            port:
              number: %d
`, backendSvc, backendPort[0])
		err := s.CreateResourceFromString(ing)
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")
		time.Sleep(5 * time.Second)

		_ = s.NewAPISIXClient().GET("/status/500").WithHeader("Host", "httpbin.org").Expect().Status(http.StatusInternalServerError)
		_ = s.NewAPISIXClient().GET("/status/504").WithHeader("Host", "httpbin.org").Expect().Status(http.StatusGatewayTimeout)
		_ = s.NewAPISIXClient().GET("/statusaaa").WithHeader("Host", "httpbin.org").Expect().Status(http.StatusNotFound).Body().Contains("404 Route Not Found")
		// Mismatched host
		_ = s.NewAPISIXClient().GET("/status/200").WithHeader("Host", "a.httpbin.org").Expect().Status(http.StatusNotFound)
	})
})

var _ = ginkgo.Describe("support ingress.networking/v1beta1", func() {
	s := scaffold.NewDefaultScaffold()

	ginkgo.It("path exact match", func() {
		backendSvc, backendPort := s.DefaultHTTPBackend()
		ing := fmt.Sprintf(`
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
  name: ingress-v1beta1
spec:
  rules:
  - host: httpbin.org
    http:
      paths:
      - path: /ip
        pathType: Exact
        backend:
          serviceName: %s
          servicePort: %d
`, backendSvc, backendPort[0])
		err := s.CreateResourceFromString(ing)
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")
		time.Sleep(5 * time.Second)

		_ = s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").Expect().Status(http.StatusOK)
		// Exact path, doesn't match /ip/aha
		_ = s.NewAPISIXClient().GET("/ip/aha").WithHeader("Host", "httpbin.org").Expect().Status(http.StatusNotFound)
		// Mismatched host
		_ = s.NewAPISIXClient().GET("/ip/aha").WithHeader("Host", "a.httpbin.org").Expect().Status(http.StatusNotFound)
	})

	ginkgo.It("path prefix match", func() {
		backendSvc, backendPort := s.DefaultHTTPBackend()
		ing := fmt.Sprintf(`
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
  name: ingress-v1beta1
spec:
  rules:
  - host: httpbin.org
    http:
      paths:
      - path: /status
        pathType: Prefix
        backend:
          serviceName: %s
          servicePort: %d
`, backendSvc, backendPort[0])
		err := s.CreateResourceFromString(ing)
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")
		time.Sleep(5 * time.Second)

		_ = s.NewAPISIXClient().GET("/status/500").WithHeader("Host", "httpbin.org").Expect().Status(http.StatusInternalServerError)
		_ = s.NewAPISIXClient().GET("/status/504").WithHeader("Host", "httpbin.org").Expect().Status(http.StatusGatewayTimeout)
		_ = s.NewAPISIXClient().GET("/statusaaa").WithHeader("Host", "httpbin.org").Expect().Status(http.StatusNotFound).Body().Contains("404 Route Not Found")
		// Mismatched host
		_ = s.NewAPISIXClient().GET("/status/200").WithHeader("Host", "a.httpbin.org").Expect().Status(http.StatusNotFound)
	})
})
