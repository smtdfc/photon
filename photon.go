package photon

import (
	photon "github.com/smtdfc/photon/pkg/base"
	echoAdapter "github.com/smtdfc/photon/pkg/echo-adapter"
)

type BaseAdapter = photon.BaseAdapter
type App = photon.App
var NewApp = photon.NewApp
type Module = photon.Module
var NewModule = photon.NewModule
type Request = photon.Request
type Response = photon.Response
type RouteHandler = photon.RouteHandler
type Route = photon.Route
type EchoAdapter = echoAdapter.EchoAdapter
var Init = echoAdapter.Init
type EchoAdapterRequest = echoAdapter.EchoAdapterRequest
var WrapRequest = echoAdapter.WrapRequest
type EchoAdapterResponse = echoAdapter.EchoAdapterResponse
var WrapResponse = echoAdapter.WrapResponse
