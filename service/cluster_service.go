/*
 * Copyright 2020-2021 the original author(https://github.com/wj596)
 *
 * <p>
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * </p>
 */
package service

import (
	"github.com/juju/errors"
	"go-mysql-transfer/global"
	"go-mysql-transfer/metrics"
	"go-mysql-transfer/service/election"
	"go-mysql-transfer/storage"
	"go-mysql-transfer/util/logs"
)

type ClusterService struct {
	electionSignal   chan bool //选举信号
	electionService  election.Service
	transferServices []*TransferService
	config_          *global.Config
}

func (s *ClusterService) Nodes() []string {
	return s.electionService.Nodes()
}

func newClusterService(config *global.Config) (*ClusterService, error) {
	var err error
	var positionDao storage.PositionStorageEx

	electionSignal := make(chan bool, 1)

	clusterService := &ClusterService{
		electionSignal:   electionSignal,
		electionService:  election.NewElection(electionSignal),
		transferServices: make([]*TransferService, 0, 8),
	}
	clusterService.config_ = config

	positionDao = storage.NewPositionStorageEx()
	//names := global.Cfg().GetServerNames()
	if err := positionDao.Initialize(); err != nil {
		return nil, errors.Trace(err)
	}

	for i := 0; i < len(config.ServerConfigs); i++ {
		t := newTransferService(config.ServerConfigs[i])
		t.positionDao = positionDao
		clusterService.transferServices = append(clusterService.transferServices, t)
		err = t.initialize()
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	return clusterService, err
}

func (s *ClusterService) boot() error {
	//log.Println("start master election")
	//if s.config_.ServerConfigs
	if s.electionService.IsCluster() {
		err := s.electionService.Elect()
		if err != nil {
			return err
		}
	} else {
		s.electionSignal <- true
	}

	s.startElectListener()

	return nil
}

func (s *ClusterService) startElectListener() {
	go func() {
		for {
			select {
			case selected := <-s.electionSignal:
				global.PrintEx("startElectListener-->selected-->", selected)
				global.SetLeaderNode(s.electionService.Leader())
				global.SetLeaderFlag(selected)
				if selected {
					metrics.SetLeaderState(metrics.LeaderState)
					s.startService()
				} else {
					metrics.SetLeaderState(metrics.FollowerState)
					s.stopService()
				}
			}
		}

	}()
}

func (s *ClusterService) startService() {
	for i := 0; i < len(s.transferServices); i++ {
		s.transferServices[i].StartUp()
		logs.Infof("启动同步服务:%s\r\n", s.transferServices[i].Name)
	}
}

func (s *ClusterService) stopService() {
	for i := 0; i < len(s.transferServices); i++ {
		s.transferServices[i].stopDump()
	}
}

func (s *ClusterService) Close() {
	for i := 0; i < len(s.transferServices); i++ {
		s.transferServices[i].Close()
	}
}
