package config

import "github.com/unionj-cloud/go-doudou/framework/configmgr"

const FrameworkName = "Go-doudou"

const (
	// Default configs for framework component
	DefaultGddBanner             = true
	DefaultGddBannerText         = FrameworkName
	DefaultGddLogLevel           = "info"
	DefaultGddLogFormat          = "text"
	DefaultGddLogReqEnable       = false
	DefaultGddGraceTimeout       = "15s"
	DefaultGddWriteTimeout       = "15s"
	DefaultGddReadTimeout        = "15s"
	DefaultGddIdleTimeout        = "60s"
	DefaultGddServiceName        = ""
	DefaultGddRouteRootPath      = ""
	DefaultGddHost               = ""
	DefaultGddPort               = 6060
	DefaultGddRetryCount         = 0
	DefaultGddManage             = true
	DefaultGddManageUser         = "admin"
	DefaultGddManagePass         = "admin"
	DefaultGddTracingMetricsRoot = FrameworkName
	DefaultGddWeight             = 1

	// Default configs for memberlist component
	DefaultGddMemSeed           = ""
	DefaultGddMemPort           = 7946
	DefaultGddMemDeadTimeout    = "60s"
	DefaultGddMemSyncInterval   = "60s"
	DefaultGddMemReclaimTimeout = "3s"
	DefaultGddMemProbeInterval  = "5s"
	DefaultGddMemProbeTimeout   = "3s"
	DefaultGddMemSuspicionMult  = 6
	DefaultGddMemRetransmitMult = 4
	DefaultGddMemGossipNodes    = 4
	DefaultGddMemGossipInterval = "500ms"
	DefaultGddMemTCPTimeout     = "30s"
	DefaultGddMemIndirectChecks = 3
	DefaultGddMemWeight         = 1
	DefaultGddMemWeightInterval = 0
	DefaultGddMemName           = ""
	DefaultGddMemHost           = ""
	DefaultGddMemCIDRsAllowed   = ""
	DefaultGddMemLogDisable     = false

	DefaultGddServiceDiscoveryMode = "memberlist"

	DefaultGddNacosNamespaceId         = "public"
	DefaultGddNacosTimeoutMs           = 10000
	DefaultGddNacosNotLoadCacheAtStart = false
	DefaultGddNacosLogDir              = "/tmp/nacos/log"
	DefaultGddNacosCacheDir            = "/tmp/nacos/cache"
	DefaultGddNacosLogLevel            = "info"
	DefaultGddNacosServerAddr          = ""
	DefaultGddNacosRegisterHost        = ""
	DefaultGddNacosClusterName         = "DEFAULT"
	DefaultGddNacosGroupName           = "DEFAULT_GROUP"

	DefaultGddNacosConfigFormat = configmgr.DotenvConfigFormat
	DefaultGddNacosConfigGroup  = "DEFAULT_GROUP"
	DefaultGddNacosConfigDataid = ""

	DefaultGddEnableResponseGzip = true
	DefaultGddAppType            = "rest"
	DefaultGddConfigRemoteType   = ""

	DefaultGddApolloCluster      = "default"
	DefaultGddApolloAddr         = ""
	DefaultGddApolloNamespace    = "application.properties"
	DefaultGddApolloBackupEnable = true
	DefaultGddApolloBackupPath   = ""
	DefaultGddApolloSecret       = ""
	DefaultGddApolloMuststart    = false
	DefaultGddApolloLogEnable    = false

	// DefaultGddSqlLogEnable only for doc purpose
	DefaultGddSqlLogEnable = false
)
