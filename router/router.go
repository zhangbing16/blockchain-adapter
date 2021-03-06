/*
Copyright ArxanChain Ltd. 2020 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package router

import (
	"fmt"
	"net/http"

	"github.com/csiabb/blockchain-adapter/common/log"
	srvctx "github.com/csiabb/blockchain-adapter/context"
	"github.com/csiabb/blockchain-adapter/controllers/blockchain"
	"github.com/csiabb/blockchain-adapter/controllers/callback"
	"github.com/csiabb/blockchain-adapter/controllers/version"
	"github.com/csiabb/blockchain-adapter/middleware"

	"github.com/gin-gonic/gin"
)

// api version
const (
	APIVersion = "v1"
)

var (
	logger = log.MustGetLogger("router")

	// checkServerVerion
	checkVersionURL = "version"

	// main url prefix
	apiPrefix = fmt.Sprintf("api/%s", APIVersion)

	// controller url
	accountsURL   = "blockchain/accounts"
	publicitesURL = "blockchain/publicities"

	// blockchain callback url
	arxanchainCallbckURL = "blockchain/callback/arxan"
)

// url path
const ()

// Router service router
type Router struct {
	context           *srvctx.Context
	versionHandler    *version.RestHandler
	blockchainHandler *blockchain.RestHandler
	callbackHandler   *callback.RestHandler
}

// InitRouter init router
func (r *Router) InitRouter(ctx *srvctx.Context) error {
	if nil == ctx {
		return fmt.Errorf("param is nil")
	}

	r.context = ctx

	// Init version handler
	var err error
	r.versionHandler, err = version.NewRestHandler(r.context)
	if err != nil {
		logger.Errorf("Failed to create version rest http handler instance, %+v", err)
		return err
	}

	// Init blockchain handler
	r.blockchainHandler, err = blockchain.NewRestHandler(r.context)
	if err != nil {
		logger.Errorf("Failed to create blockchain rest http handler instance, %+v", err)
		return err
	}

	// Init blockchain callback handler
	r.callbackHandler, err = callback.NewRestHandler(r.context)
	if err != nil {
		logger.Errorf("Failed to create blockchain callback notify rest http handler instance, %+v", err)
		return err
	}

	return nil
}

// SetupRouter add routes for rest api server
func (r *Router) SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Delims("{{", "}}")
	router.Use(Cors())

	// service version
	router.GET(checkVersionURL, r.versionHandler.Version)

	// v1 group api
	apiPrefix := router.Group(apiPrefix)
	{
		// log reponse and request
		apiPrefix.Use(middleware.RequestResponseLogger())

		apiPrefix.POST(accountsURL, r.blockchainHandler.CreateAccount)
		apiPrefix.POST(publicitesURL, r.blockchainHandler.PublicityData)

		apiPrefix.POST(arxanchainCallbckURL, r.callbackHandler.ArxanchainCallback)
	}
	return router
}

// Cors ...
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()

	}
}
