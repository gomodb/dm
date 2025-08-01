/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */

package dm

import (
	"bufio"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/gomodb/dm/util"
)

var LogDirDef, _ = os.Getwd()

var StatDirDef, _ = os.Getwd()

const (
	DEFAULT_PORT int32 = 5236

	//log level
	LOG_OFF int = 0

	LOG_ERROR int = 1

	LOG_WARN int = 2

	LOG_SQL int = 3

	LOG_INFO int = 4

	LOG_DEBUG int = 5

	LOG_ALL int = 9

	//stat
	STAT_SQL_REMOVE_LATEST int = 0

	STAT_SQL_REMOVE_OLDEST int = 1

	// 编码字符集
	ENCODING_UTF8 string = "UTF-8"

	ENCODING_EUCKR string = "EUC-KR"

	ENCODING_GB18030 string = "GB18030"

	DbAliveCheckFreqDef = 0

	LocaleDef = 0

	// log
	LogLevelDef = LOG_OFF // 日志级别：off, error, warn, sql, info, all

	LogFlushFreqDef = 10 // 日志刷盘时间s (>=0)

	LogFlushQueueSizeDef = 100 //日志队列大小

	LogBufferSizeDef = 32 * 1024 // 日志缓冲区大小 (>0)

	// stat
	StatEnableDef = false //

	StatFlushFreqDef = 3 // 日志刷盘时间s (>=0)

	StatSlowSqlCountDef = 100 // 慢sql top行数，(0-1000)

	StatHighFreqSqlCountDef = 100 // 高频sql top行数， (0-1000)

	StatSqlMaxCountDef = 100000 // sql 统计最大值(0-100000)

	StatSqlRemoveModeDef = STAT_SQL_REMOVE_LATEST // 记录sql数超过最大值时，sql淘汰方式
)

var (
	DbAliveCheckFreq = DbAliveCheckFreqDef

	Locale = LocaleDef // 0:简体中文 1：英文 2:繁体中文

	// log
	LogLevel = LogLevelDef // 日志级别：off, error, warn, sql, info, all

	LogDir = LogDirDef

	LogFlushFreq = LogFlushFreqDef // 日志刷盘时间s (>=0)

	LogFlushQueueSize = LogFlushQueueSizeDef

	LogBufferSize = LogBufferSizeDef // 日志缓冲区大小 (>0)

	// stat
	StatEnable = StatEnableDef //

	StatDir = StatDirDef // jdbc工作目录,所有生成的文件都在该目录下

	StatFlushFreq = StatFlushFreqDef // 日志刷盘时间s (>=0)

	StatSlowSqlCount = StatSlowSqlCountDef // 慢sql top行数，(0-1000)

	StatHighFreqSqlCount = StatHighFreqSqlCountDef // 高频sql top行数， (0-1000)

	StatSqlMaxCount = StatSqlMaxCountDef // sql 统计最大值(0-100000)

	StatSqlRemoveMode = StatSqlRemoveModeDef // 记录sql数超过最大值时，sql淘汰方式

	/*---------------------------------------------------------------*/
	ServerGroupMap = make(map[string]*epGroup)

	GlobalProperties = NewProperties()
)

// filePath: dm_svc.conf 文件路径
func load(filePath string) {
	if filePath == "" {
		switch runtime.GOOS {
		case "windows":
			filePath = os.Getenv("SystemRoot") + "\\system32\\dm_svc.conf"
		case "linux":
			filePath = "/etc/dm_svc.conf"
		default:
			return
		}
	}
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	fileReader := bufio.NewReader(file)

	// GlobalProperties = NewProperties()
	var groupProps *Properties
	var line string //dm_svc.conf读取到的一行

	for line, err = fileReader.ReadString('\n'); line != "" && (err == nil || err == io.EOF); line, err = fileReader.ReadString('\n') {
		// 去除#标记的注释
		if notesIndex := strings.IndexByte(line, '#'); notesIndex != -1 {
			line = line[:notesIndex]
		}
		// 去除前后多余的空格
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			groupName := strings.ToLower(line[1 : len(line)-1])
			dbGroup, ok := ServerGroupMap[groupName]
			if groupName == "" || !ok {
				continue
			}
			groupProps = dbGroup.props
			if groupProps.IsNil() {
				groupProps = NewProperties()
				groupProps.SetProperties(GlobalProperties)
				dbGroup.props = groupProps
			}

		} else {
			cfgInfo := strings.Split(line, "=")
			if len(cfgInfo) < 2 {
				continue
			}
			key := strings.TrimSpace(cfgInfo[0])
			value := strings.TrimSpace(cfgInfo[1])
			if strings.HasPrefix(value, "(") && strings.HasSuffix(value, ")") {
				value = strings.TrimSpace(value[1 : len(value)-1])
			}
			if key == "" || value == "" {
				continue
			}
			// 区分属性是全局的还是组的
			var success bool
			if groupProps.IsNil() {
				success = SetServerGroupProperties(GlobalProperties, key, value)
			} else {
				success = SetServerGroupProperties(groupProps, key, value)
			}
			if !success {
				var serverGroup = parseServerName(key, value)
				if serverGroup != nil {
					serverGroup.props = NewProperties()
					serverGroup.props.SetProperties(GlobalProperties)
					ServerGroupMap[strings.ToLower(key)] = serverGroup
				}
			}
		}
	}
}

func SetServerGroupProperties(props *Properties, key string, value string) bool {
	if util.StringUtil.EqualsIgnoreCase(key, "ADDRESS_REMAP") {
		tmp := props.GetString(AddressRemapKey, "")
		props.Set(AddressRemapKey, tmp+"("+value+")")
	} else if util.StringUtil.EqualsIgnoreCase(key, "ALWAYS_ALLOW_COMMIT") {
		props.Set(AlwayseAllowCommitKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "APP_NAME") {
		props.Set(AppNameKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "AUTO_COMMIT") {
		props.Set(AutoCommitKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "BATCH_ALLOW_MAX_ERRORS") {
		props.Set(BatchAllowMaxErrorsKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "BATCH_CONTINUE_ON_ERROR") ||
		util.StringUtil.EqualsIgnoreCase(key, "CONTINUE_BATCH_ON_ERROR") {
		props.Set(ContinueBatchOnErrorKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "BATCH_NOT_ON_CALL") {
		props.Set(BatchNotOnCallKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "BATCH_TYPE") {
		props.Set(BatchTypeKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "BUF_PREFETCH") {
		props.Set(BufPrefetchKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "CIPHER_PATH") {
		props.Set(CipherPathKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "CLUSTER") {
		props.Set(ClusterKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "COLUMN_NAME_UPPER_CASE") {
		props.Set(ColumnNameUpperCaseKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "COLUMN_NAME_CASE") {
		props.Set(ColumnNameCaseKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "COMPATIBLE_MODE") {
		props.Set(CompatibleModeKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "COMPRESS") ||
		util.StringUtil.EqualsIgnoreCase(key, "COMPRESS_MSG") {
		props.Set(CompressKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "COMPRESS_ID") {
		props.Set(CompressIdKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "CONNECT_TIMEOUT") {
		props.Set(ConnectTimeoutKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "DO_SWITCH") ||
		util.StringUtil.EqualsIgnoreCase(key, "AUTO_RECONNECT") {
		props.Set(DoSwitchKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "ENABLE_RS_CACHE") {
		props.Set(EnRsCacheKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "EP_SELECTION") {
		props.Set(EpSelectorKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "ESCAPE_PROCESS") {
		props.Set(EscapeProcessKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "IS_BDTA_RS") {
		props.Set(IsBdtaRSKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "KEY_WORDS") ||
		util.StringUtil.EqualsIgnoreCase(key, "KEYWORDS") {
		props.Set(KeywordsKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LANGUAGE") {
		props.Set(LanguageKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOB_MODE") {
		props.Set(LobModeKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOG_DIR") {
		props.Set(LogDirKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOG_FLUSH_FREQ") {
		props.Set(LogFlushFreqKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOG_LEVEL") {
		props.Set(LogLevelKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOGIN_DSC_CTRL") {
		props.Set(LoginDscCtrlKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOGIN_ENCRYPT") {
		props.Set(LoginEncryptKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOGIN_MODE") {
		props.Set(LoginModeKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOGIN_STATUS") {
		props.Set(LoginStatusKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "MAX_ROWS") {
		props.Set(MaxRowsKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "MPP_LOCAL") {
		props.Set(MppLocalKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "OS_NAME") {
		props.Set(OsNameKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "RS_CACHE_SIZE") {
		props.Set(RsCacheSizeKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "RS_REFRESH_FREQ") {
		props.Set(RsRefreshFreqKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "RW_HA") {
		props.Set(RwHAKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "RW_IGNORE_SQL") {
		props.Set(RwIgnoreSqlKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "RW_PERCENT") {
		props.Set(RwPercentKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "RW_SEPARATE") {
		props.Set(RwSeparateKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "RW_STANDBY_RECOVER_TIME") {
		props.Set(RwStandbyRecoverTimeKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "SCHEMA") {
		props.Set(SchemaKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "SESS_ENCODE") {
		if IsSupportedCharset(value) {
			props.Set("sessEncode", value)
		}
	} else if util.StringUtil.EqualsIgnoreCase(key, "SESSION_TIMEOUT") {
		props.Set(SessionTimeoutKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "SOCKET_TIMEOUT") {
		props.Set(SocketTimeoutKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "SSL_FILES_PATH") {
		props.Set(SslFilesPathKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "STAT_DIR") {
		props.Set(StatDirKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "STAT_ENABLE") {
		props.Set(StatEnableKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "STAT_FLUSH_FREQ") {
		props.Set(StatFlushFreqKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "STAT_HIGH_FREQ_SQL_COUNT") {
		props.Set(StatHighFreqSqlCountKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "STAT_SLOW_SQL_COUNT") {
		props.Set(StatSlowSqlCountKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "STAT_SQL_MAX_COUNT") {
		props.Set(StatSqlMaxCountKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "STAT_SQL_REMOVE_MODE") {
		props.Set(StatSqlRemoveModeKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "SWITCH_INTERVAL") {
		props.Set(SwitchIntervalKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "SWITCH_TIME") ||
		util.StringUtil.EqualsIgnoreCase(key, "SWITCH_TIMES") {
		props.Set(SwitchTimesKey, value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "TIME_ZONE") {
		props.Set(TimeZoneKey, value)
		props.Set("localTimezone", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "USER_REMAP") {
		tmp := props.GetString(UserRemapKey, "")
		props.Set(UserRemapKey, tmp+"("+value+")")
	} else {
		return false
	}
	return true
}

func parseServerName(name string, value string) *epGroup {
	values := strings.Split(value, ",")

	var tmpVals []string
	var tmpName string
	var tmpPort int
	var svrList = make([]*ep, 0, len(values))

	for _, v := range values {

		var tmp *ep
		// 先查找IPV6,以[]包括
		begin := strings.IndexByte(v, '[')
		end := -1
		if begin != -1 {
			end = strings.IndexByte(v[begin:], ']')
		}
		if end != -1 {
			tmpName = v[begin+1 : end]

			// port
			if portIndex := strings.IndexByte(v[end:], ':'); portIndex != -1 {
				tmpPort, _ = strconv.Atoi(strings.TrimSpace(v[portIndex+1:]))
			} else {
				tmpPort = int(DEFAULT_PORT)
			}
			tmp = newEP(tmpName, int32(tmpPort))
			svrList = append(svrList, tmp)
			continue
		}
		// IPV4
		tmpVals = strings.Split(v, ":")
		tmpName = strings.TrimSpace(tmpVals[0])
		if len(tmpVals) >= 2 {
			tmpPort, _ = strconv.Atoi(tmpVals[1])
		} else {
			tmpPort = int(DEFAULT_PORT)
		}
		tmp = newEP(tmpName, int32(tmpPort))
		svrList = append(svrList, tmp)
	}

	if len(svrList) == 0 {
		return nil
	}
	return newEPGroup(name, svrList)
}

func setDriverAttributes(props *Properties) {
	if props == nil || props.Len() == 0 {
		return
	}

	parseLanguage(props.GetString(LanguageKey, "cn"))
	DbAliveCheckFreq = props.GetInt(DbAliveCheckFreqKey, DbAliveCheckFreqDef, 1, int(INT32_MAX))

	//// log
	//LogLevel = ParseLogLevel(props)
	//LogDir = util.StringUtil.FormatDir(props.GetTrimString(LogDirKey, LogDirDef))
	//LogBufferSize = props.GetInt(LogBufferSizeKey, LogBufferSizeDef, 1, int(INT32_MAX))
	//LogFlushFreq = props.GetInt(LogFlushFreqKey, LogFlushFreqDef, 1, int(INT32_MAX))
	//LogFlushQueueSize = props.GetInt(LogFlusherQueueSizeKey, LogFlushQueueSizeDef, 1, int(INT32_MAX))
	//
	//// stat
	//StatEnable = props.GetBool(StatEnableKey, StatEnableDef)
	//StatDir = util.StringUtil.FormatDir(props.GetTrimString(StatDirKey, StatDirDef))
	//StatFlushFreq = props.GetInt(StatFlushFreqKey, StatFlushFreqDef, 1, int(INT32_MAX))
	//StatHighFreqSqlCount = props.GetInt(StatHighFreqSqlCountKey, StatHighFreqSqlCountDef, 0, 1000)
	//StatSlowSqlCount = props.GetInt(StatSlowSqlCountKey, StatSlowSqlCountDef, 0, 1000)
	//StatSqlMaxCount = props.GetInt(StatSqlMaxCountKey, StatSqlMaxCountDef, 0, 100000)
	//parseStatSqlRemoveMode(props)
}

func parseLanguage(value string) {
	if util.StringUtil.EqualsIgnoreCase("cn", value) {
		Locale = 0
	} else if util.StringUtil.EqualsIgnoreCase("en", value) {
		Locale = 1
	}
}

func IsSupportedCharset(charset string) bool {
	if util.StringUtil.EqualsIgnoreCase(ENCODING_UTF8, charset) || util.StringUtil.EqualsIgnoreCase(ENCODING_GB18030, charset) || util.StringUtil.EqualsIgnoreCase(ENCODING_EUCKR, charset) {
		return true
	}
	return false
}

func ParseLogLevel(props *Properties) int {
	logLevel := LOG_OFF
	value := props.GetString(LogLevelKey, "")
	if value != "" && !util.StringUtil.IsDigit(value) {
		if util.StringUtil.EqualsIgnoreCase("debug", value) {
			logLevel = LOG_DEBUG
		} else if util.StringUtil.EqualsIgnoreCase("info", value) {
			logLevel = LOG_INFO
		} else if util.StringUtil.EqualsIgnoreCase("sql", value) {
			logLevel = LOG_SQL
		} else if util.StringUtil.EqualsIgnoreCase("warn", value) {
			logLevel = LOG_WARN
		} else if util.StringUtil.EqualsIgnoreCase("error", value) {
			logLevel = LOG_ERROR
		} else if util.StringUtil.EqualsIgnoreCase("off", value) {
			logLevel = LOG_OFF
		} else if util.StringUtil.EqualsIgnoreCase("all", value) {
			logLevel = LOG_ALL
		}
	} else {
		logLevel = props.GetInt(LogLevelKey, logLevel, LOG_OFF, LOG_INFO)
	}

	return logLevel
}
