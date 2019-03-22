package services

import (
	"sync"
	"sync/atomic"

	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/alekns/tinyrstats/pkg/helpers/str"
)

type statsBucket = map[string]int

type statsServiceImpl struct {
	bucketsCount int
	buckets      []statsBucket
	mtxs         []sync.RWMutex

	callResponseMax int32
	callResponseMin int32
}

// init .
func (si *statsServiceImpl) init() *statsServiceImpl {
	si.mtxs = make([]sync.RWMutex, si.bucketsCount)
	si.buckets = make([]statsBucket, si.bucketsCount)
	for bucketIndex := 0; bucketIndex < si.bucketsCount; bucketIndex++ {
		si.buckets[bucketIndex] = make(statsBucket)
	}
	return si
}

// GetAllHosts .
func (si *statsServiceImpl) GetAllHosts() monitor.StatsHostsInfo {
	result := make(monitor.StatsHostsInfo, 0)

	for bucketIndex := 0; bucketIndex < si.bucketsCount; bucketIndex++ {
		si.mtxs[bucketIndex].RLock()
		for host, count := range si.buckets[bucketIndex] {
			result[host] = count
		}
		si.mtxs[bucketIndex].RUnlock()
	}

	return result
}

// GetMinMax .
func (si *statsServiceImpl) GetMinMax() (int32, int32) {
	return si.callResponseMin, si.callResponseMax
}

// AddHost .
func (si *statsServiceImpl) AddHost(host string, delta int) {
	bucketIndex := str.BasicStrHash(host) % si.bucketsCount
	si.mtxs[bucketIndex].Lock()
	si.buckets[bucketIndex][host] = si.buckets[bucketIndex][host] + delta
	si.mtxs[bucketIndex].Unlock()
}

// AddMinMax .
func (si *statsServiceImpl) AddMinMax(isMax bool, delta int32) {
	if isMax {
		atomic.AddInt32(&si.callResponseMax, delta)
		return
	}

	atomic.AddInt32(&si.callResponseMin, delta)
}

// DeleteHost .
func (si *statsServiceImpl) DeleteHost(host string) {
	bucketIndex := str.BasicStrHash(host) % si.bucketsCount
	si.mtxs[bucketIndex].Lock()
	if _, ok := si.buckets[bucketIndex][host]; ok {
		delete(si.buckets[bucketIndex], host)
	}
	si.mtxs[bucketIndex].Unlock()
}

// NewStatsServiceImpl .
func NewStatsServiceImpl(settings *config.StatsSettings) monitor.StatsService {
	instance := &statsServiceImpl{
		bucketsCount: settings.BucketsCount,
	}
	return instance.init()
}
