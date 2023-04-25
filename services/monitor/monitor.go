package monitor

import (
	"sync"
	"time"

	"github.com/staking4all/celestia-monitoring-bot/services"
	"github.com/staking4all/celestia-monitoring-bot/services/cosmos"
	"github.com/staking4all/celestia-monitoring-bot/services/models"
	"go.uber.org/zap"
)

type monitorService struct {
	alertState     map[string]map[string]*models.ValidatorAlertState
	valStats       map[string]models.ValidatorStats
	alertStateLock sync.Mutex

	config models.Config
	client services.CosmosClient
	ns     services.NotificationService
}

func NewMonitorService(config models.Config) (services.MonitorService, error) {
	m := &monitorService{
		alertState:     make(map[string]map[string]*models.ValidatorAlertState),
		alertStateLock: sync.Mutex{},

		config: config,
	}

	client, err := cosmos.NewCosmosClient(config.ValidatorsMonitor.RPC, config.ValidatorsMonitor.ChainID)
	if err != nil {
		return nil, err
	}

	m.client = client

	return m, nil
}

func (m *monitorService) Add(userID string, validator *models.Validator) {
	m.alertStateLock.Lock()
	defer m.alertStateLock.Unlock()

	if m.alertState[validator.Address] == nil {
		m.alertState[validator.Address] = make(map[string]*models.ValidatorAlertState)
	}

	m.alertState[validator.Address][userID] = &models.ValidatorAlertState{
		UserValidator:              validator,
		AlertTypeCounts:            make(map[models.AlertType]int64),
		SentryGRPCErrorCounts:      make(map[string]int64),
		SentryOutOfSyncErrorCounts: make(map[string]int64),
		SentryHaltErrorCounts:      make(map[string]int64),
		SentryLatestHeight:         make(map[string]int64),
	}
}

func (m *monitorService) Remove(userID string, address string) {
	m.alertStateLock.Lock()
	defer m.alertStateLock.Unlock()

	if m.alertState[address] != nil {
		delete(m.alertState[address], userID)

		if len(m.alertState[address]) == 0 {
			delete(m.alertState, address)
		}
	}
}

func (m *monitorService) Stop() error {
	// TODO: implement stop function
	return nil
}

func (m *monitorService) Run() error {
	ticker := time.NewTicker(10 * time.Second)

	concurrentGoroutines := make(chan struct{}, m.config.ValidatorsMonitor.MaxNbConcurrentGoroutines)

	zap.L().Info("starting node monitoring")
	for range ticker.C {
		zap.L().Info("checking data")

		slashingInfo, err := m.client.GetSlashingInfo()
		if err != nil {
			zap.L().Error("error retriving slashing info", zap.Error(err))
			continue
		}

		status, err := m.client.GetNodeStatus()
		if err != nil {
			zap.L().Error("error retriving node status", zap.Error(err))
			continue
		}

		// check node data
		if status.SyncInfo.CatchingUp {
			zap.L().Error("error retriving node status: node out of sync")
			continue
		}

		timeSinceLastBlock := time.Now().UnixNano() - status.SyncInfo.LatestBlockTime.UnixNano()
		if timeSinceLastBlock > m.config.ValidatorsMonitor.HaltThresholdNanoseconds {
			zap.L().Error("error retriving node status: chain halt", zap.Int64("time", timeSinceLastBlock))
			continue
		}

		m.valStats = make(map[string]models.ValidatorStats)
		statsLock := sync.Mutex{}
		wg := sync.WaitGroup{}
		for addr := range m.alertState {
			zap.L().Debug("checking node", zap.String("address", addr))
			wg.Add(1)
			go func(addr string) {
				defer func() {
					wg.Done()
					<-concurrentGoroutines
				}()
				concurrentGoroutines <- struct{}{}
				valStats, err := m.getData(addr, slashingInfo, status)
				if err != nil {
					zap.L().Error("error getting validator info", zap.Error(err))
					return
				}
				valStats.DetermineAggregatedErrorsAndAlertLevel()
				statsLock.Lock()
				m.valStats[addr] = *valStats
				statsLock.Unlock()
			}(addr)
		}
		wg.Wait()

		for addr, stats := range m.valStats {

			// get users subscribe
			for userID, val := range m.alertState[addr] {
				notification := stats.GetAlertNotification(val, stats.Errs)
				if notification != nil {
					m.ns.SendValidatorAlertNotification(userID, val.UserValidator, stats, notification)
				}
			}
		}
	}

	return nil
}
