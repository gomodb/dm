/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gomodb/dm/util"
)

type ExecuteTypeEnum int

const (
	Execute ExecuteTypeEnum = iota
	ExecuteQuery
	ExecuteUpdate
)

var idGenerator int64 = 0

func generateId() string {
	return time.Now().String() + strconv.Itoa(int(atomic.AddInt64(&idGenerator, 1)))
}

func getInt64(counter *int64, reset bool) int64 {
	if reset {
		return atomic.SwapInt64(counter, 0)
	}
	return atomic.LoadInt64(counter)
}

type SqlStatValue struct {
	id string

	sql string

	sqlHash int64

	dataSource string

	dataSourceId string

	executeLastStartTime int64

	executeBatchSizeTotal int64

	executeBatchSizeMax int64

	executeSuccessCount int64

	executeSpanNanoTotal int64

	executeSpanNanoMax int64

	runningCount int64

	concurrentMax int64

	resultSetHoldTimeNano int64

	executeAndResultSetHoldTime int64

	executeNanoSpanMaxOccurTime int64

	executeErrorCount int64

	executeErrorLast error

	executeErrorLastMessage string

	executeErrorLastStackTrace string

	executeErrorLastTime int64

	updateCount int64

	updateCountMax int64

	fetchRowCount int64

	fetchRowCountMax int64

	inTransactionCount int64

	lastSlowParameters string

	clobOpenCount int64

	blobOpenCount int64

	readStringLength int64

	readBytesLength int64

	inputStreamOpenCount int64

	readerOpenCount int64

	histogram_0_1 int64

	histogram_1_10 int64

	histogram_10_100 int64

	histogram_100_1000 int64

	histogram_1000_10000 int64

	histogram_10000_100000 int64

	histogram_100000_1000000 int64

	histogram_1000000_more int64

	executeAndResultHoldTime_0_1 int64

	executeAndResultHoldTime_1_10 int64

	executeAndResultHoldTime_10_100 int64

	executeAndResultHoldTime_100_1000 int64

	executeAndResultHoldTime_1000_10000 int64

	executeAndResultHoldTime_10000_100000 int64

	executeAndResultHoldTime_100000_1000000 int64

	executeAndResultHoldTime_1000000_more int64

	fetchRowCount_0_1 int64

	fetchRowCount_1_10 int64

	fetchRowCount_10_100 int64

	fetchRowCount_100_1000 int64

	fetchRowCount_1000_10000 int64

	fetchRowCount_10000_more int64

	updateCount_0_1 int64

	updateCount_1_10 int64

	updateCount_10_100 int64

	updateCount_100_1000 int64

	updateCount_1000_10000 int64

	updateCount_10000_more int64
}

func newSqlStatValue() *SqlStatValue {
	ssv := new(SqlStatValue)
	return ssv
}

func (ssv *SqlStatValue) getExecuteHistogram() []int64 {
	return []int64{
		ssv.histogram_0_1,
		ssv.histogram_1_10,
		ssv.histogram_10_100,
		ssv.histogram_100_1000,
		ssv.histogram_1000_10000,
		ssv.histogram_10000_100000,
		ssv.histogram_100000_1000000,
		ssv.histogram_1000000_more,
	}
}

func (ssv *SqlStatValue) getExecuteCount() int64 {
	return ssv.executeErrorCount + ssv.executeSuccessCount
}

func (ssv *SqlStatValue) getExecuteMillisMax() int64 {
	return ssv.executeSpanNanoMax / (1000 * 1000)
}

func (ssv *SqlStatValue) getExecuteMillisTotal() int64 {
	return ssv.executeSpanNanoTotal / (1000 * 1000)
}

func (ssv *SqlStatValue) getHistogramValues() []int64 {
	return []int64{

		ssv.histogram_0_1,
		ssv.histogram_1_10,
		ssv.histogram_10_100,
		ssv.histogram_100_1000,
		ssv.histogram_1000_10000,
		ssv.histogram_10000_100000,
		ssv.histogram_100000_1000000,
		ssv.histogram_1000000_more,
	}
}

func (ssv *SqlStatValue) getFetchRowCountHistogramValues() []int64 {
	return []int64{

		ssv.fetchRowCount_0_1,
		ssv.fetchRowCount_1_10,
		ssv.fetchRowCount_10_100,
		ssv.fetchRowCount_100_1000,
		ssv.fetchRowCount_1000_10000,
		ssv.fetchRowCount_10000_more,
	}
}

func (ssv *SqlStatValue) getUpdateCountHistogramValues() []int64 {
	return []int64{

		ssv.updateCount_0_1,
		ssv.updateCount_1_10,
		ssv.updateCount_10_100,
		ssv.updateCount_100_1000,
		ssv.updateCount_1000_10000,
		ssv.updateCount_10000_more,
	}
}

func (ssv *SqlStatValue) getExecuteAndResultHoldTimeHistogramValues() []int64 {
	return []int64{

		ssv.executeAndResultHoldTime_0_1,
		ssv.executeAndResultHoldTime_1_10,
		ssv.executeAndResultHoldTime_10_100,
		ssv.executeAndResultHoldTime_100_1000,
		ssv.executeAndResultHoldTime_1000_10000,
		ssv.executeAndResultHoldTime_10000_100000,
		ssv.executeAndResultHoldTime_100000_1000000,
		ssv.executeAndResultHoldTime_1000000_more,
	}
}

func (ssv *SqlStatValue) getResultSetHoldTimeMilis() int64 {
	return ssv.resultSetHoldTimeNano / (1000 * 1000)
}

func (ssv *SqlStatValue) getExecuteAndResultSetHoldTimeMilis() int64 {
	return ssv.executeAndResultSetHoldTime / (1000 * 1000)
}

func (ssv *SqlStatValue) getData() map[string]any {
	m := make(map[string]any)

	m[idConstStr] = ssv.id
	m[dataSourceConstStr] = ssv.dataSource
	m["DataSourceId"] = ssv.dataSourceId
	m[sqlConstStr] = ssv.sql
	m[executeCountConstStr] = ssv.getExecuteCount()
	m[errorCountConstStr] = ssv.executeErrorCount

	m[totalTimeConstStr] = ssv.getExecuteMillisTotal()
	m["LastTime"] = ssv.executeLastStartTime
	m[maxTimespanConstStr] = ssv.getExecuteMillisMax()
	m["LastError"] = ssv.executeErrorLast
	m[effectedRowCountConstStr] = ssv.updateCount

	m[fetchRowCountConstStr] = ssv.fetchRowCount
	m["MaxTimespanOccurTime"] = ssv.executeNanoSpanMaxOccurTime
	m["BatchSizeMax"] = ssv.executeBatchSizeMax
	m["BatchSizeTotal"] = ssv.executeBatchSizeTotal
	m[concurrentMaxConstStr] = ssv.concurrentMax

	m[runningCountConstStr] = ssv.runningCount

	if ssv.executeErrorLastMessage != "" {
		m["LastErrorMessage"] = ssv.executeErrorLastMessage
		m["LastErrorStackTrace"] = ssv.executeErrorLastStackTrace
		m["LastErrorTime"] = ssv.executeErrorLastTime
	} else {
		m["LastErrorMessage"] = ""
		m["LastErrorClass"] = ""
		m["LastErrorStackTrace"] = ""
		m["LastErrorTime"] = ""
	}

	m[urlConstStr] = ""
	m[inTransactionCountConstStr] = ssv.inTransactionCount

	m["Histogram"] = ssv.getHistogramValues()
	m["LastSlowParameters"] = ssv.lastSlowParameters
	m["ResultSetHoldTime"] = ssv.getResultSetHoldTimeMilis()
	m["ExecuteAndResultSetHoldTime"] = ssv.getExecuteAndResultSetHoldTimeMilis()
	m[fetchRowCountConstStr] = ssv.getFetchRowCountHistogramValues()

	m[effectedRowCountHistogramConstStr] = ssv.getUpdateCountHistogramValues()
	m[executeAndResultHoldTimeHistogramConstStr] = ssv.getExecuteAndResultHoldTimeHistogramValues()
	m["EffectedRowCountMax"] = ssv.updateCountMax
	m["FetchRowCountMax"] = ssv.fetchRowCountMax
	m[clobOpenCountConstStr] = ssv.clobOpenCount

	m[blobOpenCountConstStr] = ssv.blobOpenCount
	m["ReadStringLength"] = ssv.readStringLength
	m["ReadBytesLength"] = ssv.readBytesLength
	m["InputStreamOpenCount"] = ssv.inputStreamOpenCount
	m["ReaderOpenCount"] = ssv.readerOpenCount

	m["HASH"] = ssv.sqlHash

	m[executeHoldTimeHistogramConstStr] = ssv.getExecuteHistogram()

	return m
}

type sqlStat struct {
	Sql string

	SqlHash int64

	Id string

	ExecuteLastStartTime int64

	ExecuteBatchSizeTotal int64

	ExecuteBatchSizeMax int64

	ExecuteSuccessCount int64

	ExecuteSpanNanoTotal int64

	ExecuteSpanNanoMax int64

	RunningCount int64

	ConcurrentMax int64

	ResultSetHoldTimeNano int64

	ExecuteAndResultSetHoldTime int64

	DataSource string

	File string

	ExecuteNanoSpanMaxOccurTime int64

	ExecuteErrorCount int64

	ExecuteErrorLast error

	ExecuteErrorLastTime int64

	UpdateCount int64

	UpdateCountMax int64

	FetchRowCount int64

	FetchRowCountMax int64

	InTransactionCount int64

	LastSlowParameters string

	Removed int64

	ClobOpenCount int64

	BlobOpenCount int64

	ReadStringLength int64

	ReadBytesLength int64

	InputStreamOpenCount int64

	ReaderOpenCount int64

	Histogram_0_1 int64

	Histogram_1_10 int64

	Histogram_10_100 int64

	Histogram_100_1000 int64

	Histogram_1000_10000 int64

	Histogram_10000_100000 int64

	Histogram_100000_1000000 int64

	Histogram_1000000_more int64

	ExecuteAndResultHoldTime_0_1 int64

	ExecuteAndResultHoldTime_1_10 int64

	ExecuteAndResultHoldTime_10_100 int64

	ExecuteAndResultHoldTime_100_1000 int64

	ExecuteAndResultHoldTime_1000_10000 int64

	ExecuteAndResultHoldTime_10000_100000 int64

	ExecuteAndResultHoldTime_100000_1000000 int64

	ExecuteAndResultHoldTime_1000000_more int64

	FetchRowCount_0_1 int64

	FetchRowCount_1_10 int64

	FetchRowCount_10_100 int64

	FetchRowCount_100_1000 int64

	FetchRowCount_1000_10000 int64

	FetchRowCount_10000_more int64

	UpdateCount_0_1 int64

	UpdateCount_1_10 int64

	UpdateCount_10_100 int64

	UpdateCount_100_1000 int64

	UpdateCount_1000_10000 int64

	UpdateCount_10000_more int64

	DataSourceId string
}

func NewSqlStat(sql string) *sqlStat {
	s := new(sqlStat)
	s.Sql = sql
	s.Id = "SQL" + generateId()
	return s
}

func (s *sqlStat) getValue(reset bool) *SqlStatValue {
	ssv := newSqlStatValue()
	ssv.dataSource = s.DataSource
	ssv.dataSourceId = s.DataSourceId
	ssv.sql = s.Sql
	ssv.sqlHash = s.SqlHash
	ssv.id = s.Id
	ssv.executeLastStartTime = s.ExecuteLastStartTime
	if reset {
		s.ExecuteLastStartTime = 0
	}

	ssv.executeBatchSizeTotal = getInt64(&s.ExecuteBatchSizeTotal, reset)
	ssv.executeBatchSizeMax = getInt64(&s.ExecuteBatchSizeMax, reset)
	ssv.executeSuccessCount = getInt64(&s.ExecuteSuccessCount, reset)
	ssv.executeSpanNanoTotal = getInt64(&s.ExecuteSpanNanoTotal, reset)
	ssv.executeSpanNanoMax = getInt64(&s.ExecuteSpanNanoMax, reset)
	ssv.executeNanoSpanMaxOccurTime = s.ExecuteNanoSpanMaxOccurTime
	if reset {
		s.ExecuteNanoSpanMaxOccurTime = 0
	}

	ssv.runningCount = s.RunningCount
	ssv.concurrentMax = getInt64(&s.ConcurrentMax, reset)
	ssv.executeErrorCount = getInt64(&s.ExecuteErrorCount, reset)
	ssv.executeErrorLast = s.ExecuteErrorLast
	if reset {
		s.ExecuteErrorLast = nil
	}

	ssv.executeErrorLastTime = s.ExecuteErrorLastTime
	if reset {
		ssv.executeErrorLastTime = 0
	}

	ssv.updateCount = getInt64(&s.UpdateCount, reset)
	ssv.updateCountMax = getInt64(&s.UpdateCountMax, reset)
	ssv.fetchRowCount = getInt64(&s.FetchRowCount, reset)
	ssv.fetchRowCountMax = getInt64(&s.FetchRowCountMax, reset)
	ssv.histogram_0_1 = getInt64(&s.Histogram_0_1, reset)
	ssv.histogram_1_10 = getInt64(&s.Histogram_1_10, reset)
	ssv.histogram_10_100 = getInt64(&s.Histogram_10_100, reset)
	ssv.histogram_100_1000 = getInt64(&s.Histogram_100_1000, reset)
	ssv.histogram_1000_10000 = getInt64(&s.Histogram_1000_10000, reset)
	ssv.histogram_10000_100000 = getInt64(&s.Histogram_10000_100000, reset)
	ssv.histogram_100000_1000000 = getInt64(&s.Histogram_100000_1000000, reset)
	ssv.histogram_1000000_more = getInt64(&s.Histogram_1000000_more, reset)
	ssv.lastSlowParameters = s.LastSlowParameters
	if reset {
		s.LastSlowParameters = ""
	}

	ssv.inTransactionCount = getInt64(&s.InTransactionCount, reset)
	ssv.resultSetHoldTimeNano = getInt64(&s.ResultSetHoldTimeNano, reset)
	ssv.executeAndResultSetHoldTime = getInt64(&s.ExecuteAndResultSetHoldTime, reset)
	ssv.fetchRowCount_0_1 = getInt64(&s.FetchRowCount_0_1, reset)
	ssv.fetchRowCount_1_10 = getInt64(&s.FetchRowCount_1_10, reset)
	ssv.fetchRowCount_10_100 = getInt64(&s.FetchRowCount_10_100, reset)
	ssv.fetchRowCount_100_1000 = getInt64(&s.FetchRowCount_100_1000, reset)
	ssv.fetchRowCount_1000_10000 = getInt64(&s.FetchRowCount_1000_10000, reset)
	ssv.fetchRowCount_10000_more = getInt64(&s.FetchRowCount_10000_more, reset)
	ssv.updateCount_0_1 = getInt64(&s.UpdateCount_0_1, reset)
	ssv.updateCount_1_10 = getInt64(&s.UpdateCount_1_10, reset)
	ssv.updateCount_10_100 = getInt64(&s.UpdateCount_10_100, reset)
	ssv.updateCount_100_1000 = getInt64(&s.UpdateCount_100_1000, reset)
	ssv.updateCount_1000_10000 = getInt64(&s.UpdateCount_1000_10000, reset)
	ssv.updateCount_10000_more = getInt64(&s.UpdateCount_10000_more, reset)
	ssv.executeAndResultHoldTime_0_1 = getInt64(&s.ExecuteAndResultHoldTime_0_1, reset)
	ssv.executeAndResultHoldTime_1_10 = getInt64(&s.ExecuteAndResultHoldTime_1_10, reset)
	ssv.executeAndResultHoldTime_10_100 = getInt64(&s.ExecuteAndResultHoldTime_10_100, reset)
	ssv.executeAndResultHoldTime_100_1000 = getInt64(&s.ExecuteAndResultHoldTime_100_1000, reset)
	ssv.executeAndResultHoldTime_1000_10000 = getInt64(&s.ExecuteAndResultHoldTime_1000_10000, reset)
	ssv.executeAndResultHoldTime_10000_100000 = getInt64(&s.ExecuteAndResultHoldTime_10000_100000, reset)
	ssv.executeAndResultHoldTime_100000_1000000 = getInt64(&s.ExecuteAndResultHoldTime_100000_1000000, reset)
	ssv.executeAndResultHoldTime_1000000_more = getInt64(&s.ExecuteAndResultHoldTime_1000000_more, reset)
	ssv.blobOpenCount = getInt64(&s.BlobOpenCount, reset)
	ssv.clobOpenCount = getInt64(&s.ClobOpenCount, reset)
	ssv.readStringLength = getInt64(&s.ReadStringLength, reset)
	ssv.readBytesLength = getInt64(&s.ReadBytesLength, reset)
	ssv.inputStreamOpenCount = getInt64(&s.InputStreamOpenCount, reset)
	ssv.readerOpenCount = getInt64(&s.ReaderOpenCount, reset)
	return ssv
}

func (s *sqlStat) addUpdateCount(delta int64) {
	if delta > 0 {
		atomic.AddInt64(&s.UpdateCount, delta)
	}

	for {
		max := atomic.LoadInt64(&s.UpdateCountMax)
		if delta <= max {
			break
		}
		if atomic.CompareAndSwapInt64(&s.UpdateCountMax, max, delta) {
			break
		}
	}

	if delta < 1 {
		atomic.AddInt64(&s.UpdateCount_0_1, 1)
	} else if delta < 10 {
		atomic.AddInt64(&s.UpdateCount_1_10, 1)
	} else if delta < 100 {
		atomic.AddInt64(&s.UpdateCount_10_100, 1)
	} else if delta < 1000 {
		atomic.AddInt64(&s.UpdateCount_100_1000, 1)
	} else if delta < 10000 {
		atomic.AddInt64(&s.UpdateCount_1000_10000, 1)
	} else {
		atomic.AddInt64(&s.UpdateCount_10000_more, 1)
	}
}

func (s *sqlStat) addStringReadLength(length int64) {
	atomic.AddInt64(&s.ReadStringLength, length)
}

func (s *sqlStat) addReadBytesLength(length int64) {
	atomic.AddInt64(&s.ReadBytesLength, length)
}

func (s *sqlStat) addReaderOpenCount(count int64) {
	atomic.AddInt64(&s.ReaderOpenCount, count)
}

func (s *sqlStat) addInputStreamOpenCount(count int64) {
	atomic.AddInt64(&s.InputStreamOpenCount, count)
}

func (s *sqlStat) addFetchRowCount(delta int64) {
	atomic.AddInt64(&s.FetchRowCount, delta)
	for {
		max := atomic.LoadInt64(&s.FetchRowCountMax)
		if delta <= max {
			break
		}
		if atomic.CompareAndSwapInt64(&s.FetchRowCountMax, max, delta) {
			break
		}
	}

	if delta < 1 {
		atomic.AddInt64(&s.FetchRowCount_0_1, 1)
	} else if delta < 10 {
		atomic.AddInt64(&s.FetchRowCount_1_10, 1)
	} else if delta < 100 {
		atomic.AddInt64(&s.FetchRowCount_10_100, 1)
	} else if delta < 1000 {
		atomic.AddInt64(&s.FetchRowCount_100_1000, 1)
	} else if delta < 10000 {
		atomic.AddInt64(&s.FetchRowCount_1000_10000, 1)
	} else {
		atomic.AddInt64(&s.FetchRowCount_10000_more, 1)
	}

}

func (s *sqlStat) incrementExecuteSuccessCount() {
	atomic.AddInt64(&s.ExecuteSuccessCount, 1)
}

func (s *sqlStat) incrementRunningCount() {
	val := atomic.AddInt64(&s.RunningCount, 1)

	for {
		max := atomic.LoadInt64(&s.ConcurrentMax)
		if val > max {
			if atomic.CompareAndSwapInt64(&s.ConcurrentMax, max, val) {
				break
			} else {
				continue
			}
		} else {
			break
		}
	}
}

func (s *sqlStat) decrementRunningCount() {
	atomic.AddInt64(&s.RunningCount, -1)
}

func (s *sqlStat) addExecuteTimeAndResultHoldTimeHistogramRecord(executeType ExecuteTypeEnum, firstResultSet bool, nanoSpan int64, parameters string) {
	s.addExecuteTime(nanoSpan, parameters)

	if ExecuteQuery != executeType && !firstResultSet {
		s.executeAndResultHoldTimeHistogramRecord(nanoSpan)
	}
}

func (s *sqlStat) executeAndResultHoldTimeHistogramRecord(nanoSpan int64) {
	millis := nanoSpan / 1000 / 1000

	if millis < 1 {
		atomic.AddInt64(&s.ExecuteAndResultHoldTime_0_1, 1)
	} else if millis < 10 {
		atomic.AddInt64(&s.ExecuteAndResultHoldTime_1_10, 1)
	} else if millis < 100 {
		atomic.AddInt64(&s.ExecuteAndResultHoldTime_10_100, 1)
	} else if millis < 1000 {
		atomic.AddInt64(&s.ExecuteAndResultHoldTime_100_1000, 1)
	} else if millis < 10000 {
		atomic.AddInt64(&s.ExecuteAndResultHoldTime_1000_10000, 1)
	} else if millis < 100000 {
		atomic.AddInt64(&s.ExecuteAndResultHoldTime_10000_100000, 1)
	} else if millis < 1000000 {
		atomic.AddInt64(&s.ExecuteAndResultHoldTime_100000_1000000, 1)
	} else {
		atomic.AddInt64(&s.ExecuteAndResultHoldTime_1000000_more, 1)
	}
}

func (s *sqlStat) histogramRecord(nanoSpan int64) {
	millis := nanoSpan / 1000 / 1000

	if millis < 1 {
		atomic.AddInt64(&s.Histogram_0_1, 1)
	} else if millis < 10 {
		atomic.AddInt64(&s.Histogram_1_10, 1)
	} else if millis < 100 {
		atomic.AddInt64(&s.Histogram_10_100, 1)
	} else if millis < 1000 {
		atomic.AddInt64(&s.Histogram_100_1000, 1)
	} else if millis < 10000 {
		atomic.AddInt64(&s.Histogram_1000_10000, 1)
	} else if millis < 100000 {
		atomic.AddInt64(&s.Histogram_10000_100000, 1)
	} else if millis < 1000000 {
		atomic.AddInt64(&s.Histogram_100000_1000000, 1)
	} else {
		atomic.AddInt64(&s.Histogram_1000000_more, 1)
	}
}

func (s *sqlStat) addExecuteTime(nanoSpan int64, parameters string) {
	atomic.AddInt64(&s.ExecuteSpanNanoTotal, nanoSpan)

	for {
		current := atomic.LoadInt64(&s.ExecuteSpanNanoMax)
		if current < nanoSpan {
			if atomic.CompareAndSwapInt64(&s.ExecuteSpanNanoMax, current, nanoSpan) {

				s.ExecuteNanoSpanMaxOccurTime = time.Now().UnixNano()
				s.LastSlowParameters = parameters

				break
			} else {
				continue
			}
		} else {
			break
		}
	}

	s.histogramRecord(nanoSpan)
}

func (s *sqlStat) incrementInTransactionCount() {
	atomic.AddInt64(&s.InTransactionCount, 1)
}

func (s *sqlStat) getExecuteCount() int64 {
	return s.ExecuteErrorCount + s.ExecuteSuccessCount
}

func (s *sqlStat) getData() map[string]any {
	return s.getValue(false).getData()
}

func (s *sqlStat) error(err error) {
	atomic.AddInt64(&s.ExecuteErrorCount, 1)
	s.ExecuteErrorLastTime = time.Now().UnixNano()
	s.ExecuteErrorLast = err
}

func (s *sqlStat) addResultSetHoldTimeNano2(statementExecuteNano int64, resultHoldTimeNano int64) {
	atomic.AddInt64(&s.ResultSetHoldTimeNano, resultHoldTimeNano)
	atomic.AddInt64(&s.ExecuteAndResultSetHoldTime, statementExecuteNano+resultHoldTimeNano)
	s.executeAndResultHoldTimeHistogramRecord((statementExecuteNano + resultHoldTimeNano) / 1000 / 1000)
	atomic.AddInt64(&s.UpdateCount_0_1, 1)
}

type connectionStatValue struct {
	id string

	url string

	connCount int64

	activeConnCount int64

	maxActiveConnCount int64

	executeCount int64

	errorCount int64

	stmtCount int64

	activeStmtCount int64

	maxActiveStmtCount int64

	commitCount int64

	rollbackCount int64

	clobOpenCount int64

	blobOpenCount int64

	properties string
}

func newConnectionStatValue() *connectionStatValue {
	csv := new(connectionStatValue)
	return csv
}

func (csv *connectionStatValue) getData() map[string]any {
	m := make(map[string]any)
	m[idConstStr] = csv.id
	m[urlConstStr] = csv.url
	m[connCountConstStr] = csv.connCount
	m[activeConnCountConstStr] = csv.activeConnCount
	m[maxActiveConnCountConstStr] = csv.maxActiveConnCount

	m[stmtCountConstStr] = csv.stmtCount
	m[activeStmtCountConstStr] = csv.activeStmtCount
	m[maxActiveStmtCountConstStr] = csv.maxActiveStmtCount

	m[executeCountConstStr] = csv.executeCount
	m[errorCountConstStr] = csv.errorCount
	m[commitCountConstStr] = csv.commitCount
	m[rollbackCountConstStr] = csv.rollbackCount

	m[clobOpenCountConstStr] = csv.clobOpenCount
	m[blobOpenCountConstStr] = csv.blobOpenCount

	m[propertiesConstStr] = csv.properties
	return m
}

type connectionStat struct {
	id string

	url string

	connCount int64

	activeConnCount int64

	maxActiveConnCount int64

	executeCount int64

	errorCount int64

	stmtCount int64

	activeStmtCount int64

	maxActiveStmtCount int64

	commitCount int64

	rollbackCount int64

	clobOpenCount int64

	blobOpenCount int64

	sqlStatMap map[string]*sqlStat

	maxSqlSize int

	skipSqlCount int64

	lock sync.RWMutex

	properties string
}

func newConnectionStat(url string) *connectionStat {
	cs := new(connectionStat)
	cs.maxSqlSize = StatSqlMaxCount
	cs.id = "DS" + generateId()
	cs.url = url
	cs.sqlStatMap = make(map[string]*sqlStat, 200)
	return cs
}

func (cs *connectionStat) createSqlStat(sql string) *sqlStat {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	sqlStat, ok := cs.sqlStatMap[sql]
	if !ok {
		sqlStat := NewSqlStat(sql)
		sqlStat.DataSource = cs.url
		sqlStat.DataSourceId = cs.id
		if cs.putSqlStat(sqlStat) {
			return sqlStat
		} else {
			return nil
		}
	}

	return sqlStat

}

func (cs *connectionStat) putSqlStat(sqlStat *sqlStat) bool {
	if cs.maxSqlSize > 0 && len(cs.sqlStatMap) == cs.maxSqlSize {
		if StatSqlRemoveMode == STAT_SQL_REMOVE_OLDEST {
			removeSqlStat := cs.eliminateSqlStat()
			if removeSqlStat.RunningCount > 0 || removeSqlStat.getExecuteCount() > 0 {
				atomic.AddInt64(&cs.skipSqlCount, 1)
			}
			cs.sqlStatMap[sqlStat.Sql] = sqlStat
			return true
		} else {
			if sqlStat.RunningCount > 0 || sqlStat.getExecuteCount() > 0 {
				atomic.AddInt64(&cs.skipSqlCount, 1)
			}
			return false
		}
	} else {
		cs.sqlStatMap[sqlStat.Sql] = sqlStat
		return true
	}
}

func (cs *connectionStat) eliminateSqlStat() *sqlStat {
	if cs.maxSqlSize > 0 && len(cs.sqlStatMap) == cs.maxSqlSize {
		if StatSqlRemoveMode == STAT_SQL_REMOVE_OLDEST {
			for s, item := range cs.sqlStatMap {
				if item != nil {
					delete(cs.sqlStatMap, s)
					return item
				}
			}
		}
	}
	return nil
}

func (cs *connectionStat) getSqlStatMap() map[string]*sqlStat {
	m := make(map[string]*sqlStat, len(cs.sqlStatMap))
	cs.lock.Lock()
	defer cs.lock.Unlock()
	for s, item := range cs.sqlStatMap {
		m[s] = item
	}
	return m
}

func (cs *connectionStat) incrementConn() {
	atomic.AddInt64(&cs.connCount, 1)
	atomic.AddInt64(&cs.activeConnCount, 1)
	count := atomic.LoadInt64(&cs.activeConnCount)
	if count > atomic.LoadInt64(&cs.maxActiveConnCount) {
		atomic.StoreInt64(&cs.maxActiveConnCount, count)
	}
}

func (cs *connectionStat) decrementConn() {
	atomic.AddInt64(&cs.activeConnCount, -1)
}

func (cs *connectionStat) incrementStmt() {
	atomic.AddInt64(&cs.stmtCount, 1)
	atomic.AddInt64(&cs.activeStmtCount, 1)
	count := atomic.LoadInt64(&cs.activeStmtCount)
	if count > atomic.LoadInt64(&cs.maxActiveStmtCount) {
		atomic.StoreInt64(&cs.maxActiveStmtCount, count)
	}
}

func (cs *connectionStat) decrementStmt() {
	atomic.AddInt64(&cs.activeStmtCount, -1)
}

func (cs *connectionStat) decrementStmtByActiveStmtCount(activeStmtCount int64) {
	atomic.AddInt64(&cs.activeStmtCount, -activeStmtCount)
}

func (cs *connectionStat) incrementExecuteCount() {
	atomic.AddInt64(&cs.executeCount, 1)
}

func (cs *connectionStat) incrementErrorCount() {
	atomic.AddInt64(&cs.errorCount, 1)
}

func (cs *connectionStat) incrementCommitCount() {
	atomic.AddInt64(&cs.commitCount, 1)
}

func (cs *connectionStat) incrementRollbackCount() {
	atomic.AddInt64(&cs.rollbackCount, 1)
}

func (cs *connectionStat) getValue(reset bool) *connectionStatValue {
	val := newConnectionStatValue()
	val.id = cs.id
	val.url = cs.url

	val.connCount = getInt64(&cs.connCount, reset)
	val.activeConnCount = getInt64(&cs.activeConnCount, false)
	val.maxActiveConnCount = getInt64(&cs.maxActiveConnCount, false)

	val.stmtCount = getInt64(&cs.stmtCount, reset)
	val.activeStmtCount = getInt64(&cs.activeStmtCount, false)
	val.maxActiveStmtCount = getInt64(&cs.maxActiveStmtCount, false)

	val.commitCount = getInt64(&cs.commitCount, reset)
	val.rollbackCount = getInt64(&cs.rollbackCount, reset)
	val.executeCount = getInt64(&cs.executeCount, reset)
	val.errorCount = getInt64(&cs.errorCount, reset)

	val.blobOpenCount = getInt64(&cs.blobOpenCount, reset)
	val.clobOpenCount = getInt64(&cs.clobOpenCount, reset)

	val.properties = cs.properties
	return val
}

func (cs *connectionStat) getData() map[string]any {
	return cs.getValue(false).getData()
}

type GoStat struct {
	connStatMap map[string]*connectionStat

	lock sync.RWMutex

	maxConnSize int

	skipConnCount int64
}

func newGoStat(maxConnSize int) *GoStat {
	gs := new(GoStat)
	if maxConnSize > 0 {
		gs.maxConnSize = maxConnSize
	} else {
		gs.maxConnSize = 1000
	}

	gs.connStatMap = make(map[string]*connectionStat, 16)
	return gs
}

func (gs *GoStat) createConnStat(conn *DmConnection) *connectionStat {
	url := conn.dmConnector.host + ":" + strconv.Itoa(int(conn.dmConnector.port))
	gs.lock.Lock()
	defer gs.lock.Unlock()
	connstat, ok := gs.connStatMap[url]
	if !ok {
		connstat = newConnectionStat(url)

		remove := len(gs.connStatMap) > gs.maxConnSize
		if remove && connstat.activeConnCount > 0 {
			atomic.AddInt64(&gs.skipConnCount, 1)
		}

		gs.connStatMap[url] = connstat
	}

	return connstat
}

func (gs *GoStat) getConnStatMap() map[string]*connectionStat {
	m := make(map[string]*connectionStat, len(gs.connStatMap))
	gs.lock.Lock()
	defer gs.lock.Unlock()

	for s, stat := range gs.connStatMap {
		m[s] = stat
	}
	return m
}

var sqlRowField = []string{rowNumConstStr, dataSourceConstStr, sqlConstStr, executeCountConstStr,
	totalTimeConstStr, maxTimespanConstStr, inTransactionCountConstStr, errorCountConstStr, effectedRowCountConstStr,
	fetchRowCountConstStr, runningCountConstStr, concurrentMaxConstStr, executeHoldTimeHistogramConstStr,
	executeAndResultHoldTimeHistogramConstStr, fetchRowCountHistogramConstStr, effectedRowCountHistogramConstStr}

const (
	rowNumConstStr                            = "rowNum"
	idConstStr                                = "ID"
	urlConstStr                               = "Url"
	connCountConstStr                         = "ConnCount"
	activeConnCountConstStr                   = "ActiveConnCount"
	maxActiveConnCountConstStr                = "MaxActiveConnCount"
	stmtCountConstStr                         = "StmtCount"
	activeStmtCountConstStr                   = "ActiveStmtCount"
	maxActiveStmtCountConstStr                = "MaxActiveStmtCount"
	executeCountConstStr                      = "ExecuteCount"
	errorCountConstStr                        = "ErrorCount"
	commitCountConstStr                       = "CommitCount"
	rollbackCountConstStr                     = "RollbackCount"
	clobOpenCountConstStr                     = "ClobOpenCount"
	blobOpenCountConstStr                     = "BlobOpenCount"
	propertiesConstStr                        = "Properties"
	dataSourceConstStr                        = "DataSource"
	sqlConstStr                               = "SQL"
	totalTimeConstStr                         = "TotalTime"
	maxTimespanConstStr                       = "MaxTimespan"
	inTransactionCountConstStr                = "InTransactionCount"
	effectedRowCountConstStr                  = "EffectedRowCount"
	fetchRowCountConstStr                     = "FetchRowCount"
	runningCountConstStr                      = "RunningCount"
	concurrentMaxConstStr                     = "ConcurrentMax"
	executeHoldTimeHistogramConstStr          = "ExecuteHoldTimeHistogram"
	executeAndResultHoldTimeHistogramConstStr = "ExecuteAndResultHoldTimeHistogram"
	fetchRowCountHistogramConstStr            = "FetchRowCountHistogram"
	effectedRowCountHistogramConstStr         = "EffectedRowCountHistogram"
)

var dsRowField = []string{rowNumConstStr, urlConstStr, activeConnCountConstStr,
	maxActiveConnCountConstStr, activeStmtCountConstStr, maxActiveStmtCountConstStr, executeCountConstStr, errorCountConstStr,
	commitCountConstStr, rollbackCountConstStr}

const (
	PROP_NAME_SORT            = "sort"
	PROP_NAME_SORT_FIELD      = "field"
	PROP_NAME_SORT_TYPE       = "direction"
	PROP_NAME_SEARCH          = "search"
	PROP_NAME_PAGE_NUM        = "pageNum"
	PROP_NAME_PAGE_SIZE       = "pageSize"
	PROP_NAME_PAGE_COUNT      = "pageCount"
	PROP_NAME_TOTAL_ROW_COUNT = "totalRowCount"
	PROP_NAME_FLUSH_FREQ      = "flushFreq"
	PROP_NAME_DATASOURCE_ID   = "dataSourceId"
	PROP_NAME_SQL_ID          = "sqlId"

	URL_SQL               = "sql"
	URL_SQL_DETAIL        = "sqlDetail"
	URL_DATASOURCE        = "dataSource"
	URL_DATASOURCE_DETAIL = "dataSourceDetail"

	RESULT_CODE_SUCCESS = 1
	RESULT_CODE_ERROR   = -1
	DEFAULT_PAGE_NUM    = 1
	DEFAULT_PAGE_SIZE   = int(INT32_MAX)
	DEFAULT_ORDER_TYPE  = "asc"
	DEFAULT_ORDERBY     = "DataSourceId"
)

type StatReader struct {
	connStat []map[string]any

	connStatColLens []int

	highFreqSqlStat []map[string]any

	highFreqSqlStatColLens []int

	slowSqlStat []map[string]any

	slowSqlStatColLens []int
}

func newStatReader() *StatReader {
	sr := new(StatReader)
	return sr
}

func (sr *StatReader) readConnStat(retList []string, maxCount int) (bool, []string) {
	fields := dsRowField
	isAppend := false
	if sr.connStat == nil {
		sr.connStat = sr.getConnStat("", fields)
		sr.connStatColLens = calcColLens(sr.connStat, fields, COL_MAX_LEN)
		isAppend = false
	} else {
		isAppend = true
	}
	var retContent []map[string]any
	if maxCount > 0 && len(sr.connStat) > maxCount {
		retContent = sr.connStat[0:maxCount]
		sr.connStat = sr.connStat[maxCount:len(sr.connStat)]
	} else {
		retContent = sr.connStat
		sr.connStat = nil
	}
	retList = append(retList, sr.getFormattedOutput(retContent, fields, sr.connStatColLens, isAppend))
	return sr.connStat != nil, retList
}

func (sr *StatReader) readHighFreqSqlStat(retList []string, maxCount int) (bool, []string) {
	isAppend := false
	if sr.highFreqSqlStat == nil {
		sr.highFreqSqlStat = sr.getHighFreqSqlStat(StatHighFreqSqlCount, -1, sqlRowField)
		sr.highFreqSqlStatColLens = calcColLens(sr.highFreqSqlStat, sqlRowField, COL_MAX_LEN)
		isAppend = false
	} else {
		isAppend = true
	}
	var retContent []map[string]any
	if maxCount > 0 && len(sr.highFreqSqlStat) > maxCount {
		retContent = sr.highFreqSqlStat[0:maxCount]
		sr.highFreqSqlStat = sr.highFreqSqlStat[maxCount:len(sr.highFreqSqlStat)]
	} else {
		retContent = sr.highFreqSqlStat
		sr.highFreqSqlStat = nil
	}
	retList = append(retList, sr.getFormattedOutput(retContent, sqlRowField, sr.highFreqSqlStatColLens, isAppend))
	return sr.highFreqSqlStat != nil, retList
}

func (sr *StatReader) getHighFreqSqlStat(topCount int, sqlId int,
	fields []string) []map[string]any {
	var content []map[string]any

	if topCount != 0 {
		parameters := NewProperties()
		parameters.Set(PROP_NAME_SORT_FIELD, "ExecuteCount")
		parameters.Set(PROP_NAME_SORT_TYPE, "desc")
		parameters.Set(PROP_NAME_PAGE_NUM, "1")
		parameters.Set(PROP_NAME_PAGE_SIZE, strconv.Itoa(topCount))
		content = sr.service(URL_SQL, parameters)
		if sqlId != -1 {
			matchedContent := make([]map[string]any, 0)
			for _, sqlStat := range content {
				idStr := sqlStat["ID"]
				if idStr == sqlId {
					matchedContent = append(matchedContent, sqlStat)
					break
				}
			}
			content = matchedContent
		}
	}

	if content == nil {
		content = make([]map[string]any, 0)
	} else {
		i := 1
		for _, m := range content {
			m[rowNumConstStr] = i
			i++
		}
	}
	content = addTitles(content, fields)
	return content
}

func (sr *StatReader) readSlowSqlStat(retList []string, maxCount int) (bool, []string) {
	isAppend := false
	if sr.slowSqlStat == nil {
		sr.slowSqlStat = sr.getSlowSqlStat(StatSlowSqlCount, -1, sqlRowField)
		sr.slowSqlStatColLens = calcColLens(sr.slowSqlStat, sqlRowField,
			COL_MAX_LEN)
		isAppend = false
	} else {
		isAppend = true
	}
	var retContent []map[string]any
	if maxCount > 0 && len(sr.slowSqlStat) > maxCount {
		retContent = sr.slowSqlStat[0:maxCount]
		sr.slowSqlStat = sr.slowSqlStat[maxCount:len(sr.slowSqlStat)]
	} else {
		retContent = sr.slowSqlStat
		sr.slowSqlStat = nil
	}
	retList = append(retList, sr.getFormattedOutput(retContent, sqlRowField, sr.slowSqlStatColLens, isAppend))
	return sr.slowSqlStat != nil, retList
}

func (sr *StatReader) getSlowSqlStat(topCount int, sqlId int, fields []string) []map[string]any {
	var content []map[string]any

	if topCount != 0 {
		parameters := NewProperties()
		parameters.Set(PROP_NAME_SORT_FIELD, "MaxTimespan")
		parameters.Set(PROP_NAME_SORT_TYPE, "desc")
		parameters.Set(PROP_NAME_PAGE_NUM, "1")
		parameters.Set(PROP_NAME_PAGE_SIZE, strconv.Itoa(topCount))

		content = sr.service(URL_SQL, parameters)
		if sqlId != -1 {
			matchedContent := make([]map[string]any, 0)
			for _, sqlStat := range content {
				idStr := sqlStat["ID"]
				if idStr == sqlId {
					matchedContent = append(matchedContent, sqlStat)
					break
				}
			}
			content = matchedContent
		}
	}

	if content == nil {
		content = make([]map[string]any, 0)
	} else {
		i := 1
		for _, m := range content {
			m["rowNum"] = i
			i++
		}
	}
	content = addTitles(content, fields)
	return content
}

func (sr *StatReader) getConnStat(connId string, fields []string) []map[string]any {
	content := sr.service(URL_DATASOURCE, nil)
	if connId != "" {
		matchedContent := make([]map[string]any, 0)
		for _, dsStat := range content {
			idStr := dsStat["Identity"]
			if connId == idStr {
				matchedContent = append(matchedContent, dsStat)
				break
			}
		}
		content = matchedContent
	}
	if content == nil {
		content = make([]map[string]any, 0)
	} else {
		i := 1
		for _, m := range content {
			m["rowNum"] = i
			i++
		}
	}
	content = addTitles(content, fields)
	return content
}

func (sr *StatReader) getFormattedOutput(content []map[string]any, fields []string, colLens []int,
	isAppend bool) string {
	return toTable(content, fields, colLens, true, isAppend)
}

func (sr *StatReader) parseUrl(url string) *Properties {
	parameters := NewProperties()

	if url == "" || len(strings.TrimSpace(url)) == 0 {
		return parameters
	}

	parametersStr := util.StringUtil.SubstringBetween(url, "?", "")
	if parametersStr == "" || len(parametersStr) == 0 {
		return parameters
	}

	parametersArray := strings.Split(parametersStr, "&")

	for _, parameterStr := range parametersArray {
		index := strings.Index(parametersStr, "=")
		if index <= 0 {
			continue
		}

		name := parameterStr[0:index]
		value := parameterStr[index+1:]
		parameters.Set(name, value)
	}
	return parameters
}

func (sr *StatReader) service(url string, params *Properties) []map[string]any {
	if params != nil {
		params.SetProperties(sr.parseUrl(url))
	} else {
		params = sr.parseUrl(url)
	}

	if strings.Index(url, URL_SQL) == 0 {
		array := sr.getSqlStatList(params)
		array = sr.comparatorOrderBy(array, params)
		params.Set(PROP_NAME_FLUSH_FREQ, strconv.Itoa(StatFlushFreq))
		return array
	} else if strings.Index(url, URL_SQL_DETAIL) == 0 {
		array := sr.getSqlStatDetailList(params)
		return array
	} else if strings.Index(url, URL_DATASOURCE) == 0 {
		array := sr.getConnStatList(params)
		array = sr.comparatorOrderBy(array, params)
		params.Set(PROP_NAME_FLUSH_FREQ, strconv.Itoa(StatFlushFreq))
		return array
	} else if strings.Index(url, URL_DATASOURCE_DETAIL) == 0 {
		array := sr.getConnStatDetailList(params)
		return array
	} else {
		return nil
	}
}

func (sr *StatReader) getSqlStatList(_ *Properties) []map[string]any {
	array := make([]map[string]any, 0)
	connStatMap := goStat.getConnStatMap()
	var sqlStatMap map[string]*sqlStat
	for _, connStat := range connStatMap {
		sqlStatMap = connStat.getSqlStatMap()
		for _, sqlStat := range sqlStatMap {
			data := sqlStat.getData()
			executeCount := data[executeCountConstStr]
			runningCount := data[runningCountConstStr]
			if executeCount == 0 && runningCount == 0 {
				continue
			}

			array = append(array, data)
		}
	}

	return array
}

func (sr *StatReader) getSqlStatDetailList(params *Properties) []map[string]any {
	array := make([]map[string]any, 0)
	connStatMap := goStat.getConnStatMap()
	var data *sqlStat
	sqlId := ""
	dsId := ""
	if v := params.GetString(PROP_NAME_SQL_ID, ""); v != "" {
		sqlId = v
	}
	if v := params.GetString(PROP_NAME_DATASOURCE_ID, ""); v != "" {
		dsId = v
	}
	if sqlId != "" && dsId != "" {
		for _, connStat := range connStatMap {
			if dsId != connStat.id {
				continue
			} else {
				sqlStatMap := connStat.getSqlStatMap()
				for _, sqlStat := range sqlStatMap {

					if sqlId == sqlStat.Id {
						data = sqlStat
						break
					}
				}
			}
			break
		}
	}
	if data != nil {

		array = append(array, data.getData())

	}
	return array
}

func (sr *StatReader) getConnStatList(params *Properties) []map[string]any {
	array := make([]map[string]any, 0)
	connStatMap := goStat.getConnStatMap()
	id := ""
	if v := params.GetString(PROP_NAME_DATASOURCE_ID, ""); v != "" {
		id = v
	}
	for _, connStat := range connStatMap {
		data := connStat.getData()

		connCount := data["ConnCount"]

		if connCount == 0 {
			continue
		}

		if id != "" {
			if id == connStat.id {
				array = append(array, data)
				break
			} else {
				continue
			}
		} else {

			array = append(array, data)
		}

	}
	return array
}

func (sr *StatReader) getConnStatDetailList(params *Properties) []map[string]any {
	array := make([]map[string]any, 0)
	var data *connectionStat
	connStatMap := goStat.getConnStatMap()
	id := ""
	if v := params.GetString(PROP_NAME_DATASOURCE_ID, ""); v != "" {
		id = v
	}
	if id != "" {
		for _, connStat := range connStatMap {
			if id == connStat.id {
				data = connStat
				break
			}
		}
	}
	if data != nil {
		dataValue := data.getValue(false)
		m := make(map[string]any, 2)
		m["name"] = "数据源"
		m["value"] = dataValue.url
		array = append(array, m)

		m = make(map[string]any, 2)
		m["name"] = "总会话数"
		m["value"] = dataValue.connCount
		array = append(array, m)

		m = make(map[string]any, 2)
		m["name"] = "活动会话数"
		m["value"] = dataValue.activeConnCount
		array = append(array, m)

		m = make(map[string]any, 2)
		m["name"] = "活动会话数峰值"
		m["value"] = dataValue.maxActiveStmtCount
		array = append(array, m)

		m = make(map[string]any, 2)
		m["name"] = "总句柄数"
		m["value"] = dataValue.stmtCount
		array = append(array, m)

		m = make(map[string]any, 2)
		m["name"] = "活动句柄数"
		m["value"] = dataValue.activeStmtCount
		array = append(array, m)

		m = make(map[string]any, 2)
		m["name"] = "活动句柄数峰值"
		m["value"] = dataValue.maxActiveStmtCount
		array = append(array, m)

		m = make(map[string]any, 2)
		m["name"] = "执行次数"
		m["value"] = dataValue.executeCount
		array = append(array, m)

		m = make(map[string]any, 2)
		m["name"] = "执行出错次数"
		m["value"] = dataValue.errorCount
		array = append(array, m)

		m = make(map[string]any, 2)
		m["name"] = "提交次数"
		m["value"] = dataValue.commitCount
		array = append(array, m)

		m = make(map[string]any, 2)
		m["name"] = "回滚次数"
		m["value"] = dataValue.rollbackCount
		array = append(array, m)

	}
	return array
}

type mapSlice struct {
	m          []map[string]any
	isDesc     bool
	orderByKey string
}

func newMapSlice(m []map[string]any, isDesc bool, orderByKey string) *mapSlice {
	ms := new(mapSlice)
	ms.m = m
	ms.isDesc = isDesc
	ms.orderByKey = orderByKey
	return ms
}

func (ms mapSlice) Len() int { return len(ms.m) }

func (ms mapSlice) Less(i, j int) bool {
	m1 := ms.m[i]
	m2 := ms.m[j]
	v1 := m1[ms.orderByKey]
	v2 := m2[ms.orderByKey]
	if v1 == nil {
		return true
	} else if v2 == nil {
		return false
	}

	switch v1.(type) {
	case int64:
		return v1.(int64) < v2.(int64)
	case float64:
		return v1.(float64) < v2.(float64)
	default:
		return true
	}
}

func (ms mapSlice) Swap(i, j int) {
	ms.m[i], ms.m[j] = ms.m[j], ms.m[i]
}

func (sr *StatReader) comparatorOrderBy(array []map[string]any, params *Properties) []map[string]any {
	if array == nil {
		array = make([]map[string]any, 0)
	}

	orderBy := DEFAULT_ORDERBY
	orderType := DEFAULT_ORDER_TYPE
	pageNum := DEFAULT_PAGE_NUM
	pageSize := DEFAULT_PAGE_SIZE
	if params != nil {
		if v := params.GetTrimString(PROP_NAME_SORT_FIELD, ""); v != "" {
			orderBy = v
		}

		if v := params.GetTrimString(PROP_NAME_SORT_TYPE, ""); v != "" {
			orderType = v
		}

		if v := params.GetTrimString(PROP_NAME_PAGE_NUM, ""); v != "" {
			var err error
			pageNum, err = strconv.Atoi(v)
			if err != nil {
				pageNum = DEFAULT_PAGE_NUM
			}
		}
		if v := params.GetTrimString(PROP_NAME_PAGE_SIZE, ""); v != "" {
			var err error
			pageSize, err = strconv.Atoi(v)
			if err != nil {
				pageSize = DEFAULT_PAGE_SIZE
			}
		}
	}

	rowCount := len(array)
	pageCount := max(int(math.Ceil(float64(rowCount*1.0/pageSize))), 1)

	if pageNum > pageCount {
		pageNum = pageCount
	}

	if len(array) > 0 {

		if orderBy != "" {
			sort.Sort(newMapSlice(array, DEFAULT_ORDER_TYPE != orderType, orderBy))
		}

		fromIndex := (pageNum - 1) * pageSize

		toIndex := pageNum * pageSize
		if toIndex > rowCount {
			toIndex = rowCount
		}
		array = array[fromIndex:toIndex]
	}
	sr.resetPageInfo(params, rowCount, pageCount, pageNum)
	return array
}

func (sr *StatReader) resetPageInfo(params *Properties, rowCount int, pageCount int, pageNum int) {

	if params != nil {
		v := params.GetString(PROP_NAME_PAGE_SIZE, "")
		if v != "" {

			params.Set(PROP_NAME_PAGE_COUNT, strconv.Itoa(pageCount))
			params.Set(PROP_NAME_TOTAL_ROW_COUNT, strconv.Itoa(rowCount))
			params.Set(PROP_NAME_PAGE_NUM, strconv.Itoa(pageNum))
		}
	}
}

const COL_MAX_LEN = 32

func calcColLens(objList []map[string]any, fields []string, maxColLen int) []int {

	colLen := 0
	colVal := ""
	colLens := make([]int, len(fields))
	for _, obj := range objList {
		for i := 0; i < len(fields); i++ {
			colVal = getColValue(obj[fields[i]])
			colLen = len(colVal)
			if colLen > colLens[i] {
				colLens[i] = colLen
			}
		}
	}
	if maxColLen > 0 {
		for i := 0; i < len(fields); i++ {
			if colLens[i] > maxColLen {
				colLens[i] = maxColLen
			}
		}
	}
	return colLens
}

func addTitles(objList []map[string]any, fields []string) []map[string]any {
	titleMap := make(map[string]any, len(fields))
	for i := 0; i < len(fields); i++ {
		titleMap[fields[i]] = fields[i]
	}

	dst := append(objList, titleMap)
	copy(dst[1:], dst[:len(dst)-1])
	dst[0] = titleMap
	return dst
}

func toTable(objList []map[string]any, fields []string, colLens []int,
	showAll bool, append bool) string {
	if fields == nil || objList == nil {
		return ""
	}

	if colLens == nil {
		colLens = calcColLens(objList, fields, COL_MAX_LEN)
	}

	output := &strings.Builder{}
	if !append {
		sepLine(output, colLens)
	}

	for _, obj := range objList {
		objMore := obj
		for objMore != nil {
			objMore = formateLine(output, objMore, fields, colLens, showAll)
		}
		sepLine(output, colLens)
	}

	return output.String()
}

func formateLine(output *strings.Builder, obj map[string]any, fields []string, colLens []int,
	showAll bool) map[string]any {
	hasMore := false
	objMore := make(map[string]any)
	colLen := 0
	colVal := ""
	for i := 0; i < len(fields); i++ {
		colVal = getColValue(obj[fields[i]])

		colLen = len(colVal)
		if colLen <= colLens[i] {
			output.WriteString("|")
			output.WriteString(colVal)
			blanks(output, colLens[i]-colLen)
			if showAll {
				objMore[fields[i]] = ""
			}
		} else {
			output.WriteString("|")
			if showAll {
				output.WriteString(colVal[0:colLens[i]])
				objMore[fields[i]] = colVal[colLens[i]:]
				hasMore = true
			} else {
				output.WriteString(colVal[0:colLens[i]-3] + "...")
			}
		}
	}
	output.WriteString("|")
	output.WriteString(util.StringUtil.LineSeparator())

	if hasMore {
		return objMore
	} else {
		return nil
	}
}

func sepLine(output *strings.Builder, colLens []int) {
	output.WriteString("+")
	for _, colLen := range colLens {
		for i := 0; i < colLen; i++ {
			output.WriteString("+")
		}
		output.WriteString("+")
	}
	output.WriteString(util.StringUtil.LineSeparator())
}

func blanks(output *strings.Builder, count int) {
	for count > 0 {
		output.WriteString(" ")
		count--
	}
}

func getColValue(colObj any) string {
	var colVal string
	if colObj == nil {
		colVal = ""
	} else {
		colVal = fmt.Sprint(colObj)
	}

	colVal = strings.ReplaceAll(colVal, "\t", "")
	colVal = strings.ReplaceAll(colVal, "\n", "")
	colVal = strings.ReplaceAll(colVal, "\r", "")

	return colVal
}

const (
	READ_MAX_SIZE = 100
)

type statFlusher struct {
	sr         *StatReader
	logList    []string
	date       string
	logFile    *os.File
	flushFreq  int
	filePath   string
	filePrefix string
	buffer     *Dm_build_283
}

func newStatFlusher() *statFlusher {
	sf := new(statFlusher)
	sf.sr = newStatReader()
	sf.logList = make([]string, 0, 32)
	sf.date = time.Now().Format("2006-01-02")
	sf.flushFreq = StatFlushFreq
	sf.filePath = StatDir
	sf.filePrefix = "dm_go_stat"
	sf.buffer = Dm_build_287()
	return sf
}

func (sf *statFlusher) isConnStatEnabled() bool {
	return StatEnable
}

func (sf *statFlusher) isSlowSqlStatEnabled() bool {
	return StatEnable
}

func (sf *statFlusher) isHighFreqSqlStatEnabled() bool {
	return StatEnable
}

func (sf *statFlusher) doRun() {

	for {
		if len(goStat.connStatMap) > 0 {
			sf.logList = append(sf.logList, time.Now().String())
			if sf.isConnStatEnabled() {
				sf.logList = append(sf.logList, "#connection stat")
				hasMore := true
				for hasMore {
					hasMore, sf.logList = sf.sr.readConnStat(sf.logList, READ_MAX_SIZE)
					sf.writeAndFlush(sf.logList, 0, len(sf.logList))
					sf.logList = sf.logList[0:0]
				}
			}
			if sf.isHighFreqSqlStatEnabled() {
				sf.logList = append(sf.logList, "#top "+strconv.Itoa(StatHighFreqSqlCount)+" high freq sql stat")
				hasMore := true
				for hasMore {
					hasMore, sf.logList = sf.sr.readHighFreqSqlStat(sf.logList, READ_MAX_SIZE)
					sf.writeAndFlush(sf.logList, 0, len(sf.logList))
					sf.logList = sf.logList[0:0]
				}
			}
			if sf.isSlowSqlStatEnabled() {
				sf.logList = append(sf.logList, "#top "+strconv.Itoa(StatSlowSqlCount)+" slow sql stat")
				hasMore := true
				for hasMore {
					hasMore, sf.logList = sf.sr.readSlowSqlStat(sf.logList, READ_MAX_SIZE)
					sf.writeAndFlush(sf.logList, 0, len(sf.logList))
					sf.logList = sf.logList[0:0]
				}
			}
			sf.logList = append(sf.logList, util.StringUtil.LineSeparator())
			sf.logList = append(sf.logList, util.StringUtil.LineSeparator())
			sf.writeAndFlush(sf.logList, 0, len(sf.logList))
			sf.logList = sf.logList[0:0]
			time.Sleep(time.Duration(StatFlushFreq) * time.Second)
		}
	}
}

func (sf *statFlusher) writeAndFlush(logs []string, startOff int, l int) {
	var bytes []byte
	for i := startOff; i < startOff+l; i++ {
		bytes = []byte(logs[i] + util.StringUtil.LineSeparator())

		sf.buffer.Dm_build_309(bytes, 0, len(bytes))

		if sf.buffer.Dm_build_288() >= FLUSH_SIZE {
			sf.doFlush(sf.buffer)
		}
	}

	if sf.buffer.Dm_build_288() > 0 {
		sf.doFlush(sf.buffer)
	}
}

func (sf *statFlusher) doFlush(buffer *Dm_build_283) {
	if sf.needCreateNewFile() {
		sf.closeCurrentFile()
		sf.logFile = sf.createNewFile()
	}
	if sf.logFile != nil {
		buffer.Dm_build_303(sf.logFile, buffer.Dm_build_288())
	}
}
func (sf *statFlusher) closeCurrentFile() {
	if sf.logFile != nil {
		_ = sf.logFile.Close()
		sf.logFile = nil
	}
}
func (sf *statFlusher) createNewFile() *os.File {
	sf.date = time.Now().Format("2006-01-02")
	fileName := sf.filePrefix + "_" + sf.date + "_" + strconv.Itoa(time.Now().Nanosecond()) + ".txt"
	sf.filePath = StatDir
	if len(sf.filePath) > 0 {
		if _, err := os.Stat(sf.filePath); err != nil {
			os.MkdirAll(sf.filePath, 0755)
		}
		if _, err := os.Stat(sf.filePath + fileName); err != nil {
			logFile, err := os.Create(sf.filePath + fileName)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			return logFile
		}
	}
	return nil
}
func (sf *statFlusher) needCreateNewFile() bool {
	now := time.Now().Format("2006-01-02")
	fileInfo, err := sf.logFile.Stat()
	return now != sf.date || err != nil || sf.logFile == nil || fileInfo.Size() > int64(MAX_FILE_SIZE)
}
