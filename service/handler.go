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
	"github.com/asim/mq/common"
	"go-mysql-transfer/metrics"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"

	"go-mysql-transfer/global"
	"go-mysql-transfer/model"
	"go-mysql-transfer/util/logs"
)

type handler struct {
	queue           chan interface{}
	stop            chan struct{}
	transferService *TransferService
}

func newHandler(t *TransferService) *handler {
	return &handler{
		queue:           make(chan interface{}, 4096),
		stop:            make(chan struct{}, 1),
		transferService: t,
	}
}

func (s *handler) OnRotate(e *replication.RotateEvent) error {
	s.queue <- model.PosRequest{
		Name:  string(e.NextLogName),
		Pos:   uint32(e.Position),
		Force: true,
	}
	return nil
}

func (s *handler) OnTableChanged(schema, table string) error {
	err := s.transferService.updateRule(schema, table)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (s *handler) OnDDL(nextPos mysql.Position, _ *replication.QueryEvent) error {
	s.queue <- model.PosRequest{
		Name:  nextPos.Name,
		Pos:   nextPos.Pos,
		Force: true,
	}
	return nil
}

func (s *handler) OnXID(nextPos mysql.Position) error {
	s.queue <- model.PosRequest{
		Name:  nextPos.Name,
		Pos:   nextPos.Pos,
		Force: false,
	}
	return nil
}

func (s *handler) OnRow(e *canal.RowsEvent) error {
	//logs.Infof("ServerID--->%d",e.Header.ServerID)
	//复制程序产生的binlog，不进行处理
	if e.Header.ServerID > common.MIN_REPLICATION_SLAVE &&
		e.Header.ServerID < common.MAX_REPLICATION_SLAVE &&
		e.Header.ServerID%2 == 0 {
		logs.Infof("服务ID不满足要求:0x%s", strings.ToUpper(strconv.FormatInt(int64(e.Header.ServerID), 16)))
		return nil
	}
	ruleKey := global.RuleKey(e.Table.Schema, e.Table.Name)
	if !s.transferService.ruleInsExist(ruleKey) {
		return nil
	}

	var requests []*model.RowRequest
	raw := e.ReplicationRowsEvent.Raw
	requests = make([]*model.RowRequest, 1, 1)

	r := &model.RowRequest{}
	r.RuleKey = ruleKey
	r.Action = e.Action
	r.Timestamp = e.Header.Timestamp
	r.Raw = raw
	r.ReplicationRowsEvent = e.ReplicationRowsEvent
	requests[0] = r

	s.queue <- requests

	return nil
}

func (s *handler) OnRow2(e *canal.RowsEvent) error {
	ruleKey := global.RuleKey(e.Table.Schema, e.Table.Name)
	if !s.transferService.ruleInsExist(ruleKey) {
		return nil
	}

	var requests []*model.RowRequest
	if e.Action != canal.UpdateAction {
		// 定长分配
		requests = make([]*model.RowRequest, 0, len(e.Rows))
	}

	if e.Action == canal.UpdateAction {
		for i := 0; i < len(e.Rows); i++ {
			if (i+1)%2 == 0 {
				v := new(model.RowRequest)
				v.RuleKey = ruleKey
				v.Action = e.Action
				v.Timestamp = e.Header.Timestamp
				if global.Cfg().IsReserveRawData() {
					v.Old = e.Rows[i-1]
				}
				v.Row = e.Rows[i]
				requests = append(requests, v)
			}
		}
	} else {
		for _, row := range e.Rows {
			v := new(model.RowRequest)
			v.RuleKey = ruleKey
			v.Action = e.Action
			v.Timestamp = e.Header.Timestamp
			v.Row = row
			requests = append(requests, v)
		}
	}
	s.queue <- requests

	return nil
}

func (s *handler) OnGTID(gtid mysql.GTIDSet) error {
	return nil
}

func (s *handler) OnPosSynced(pos mysql.Position, set mysql.GTIDSet, force bool) error {
	return nil
}

func (s *handler) String() string {
	return "TransferHandler"
}

func (s *handler) startListener() {
	go func() {
		interval := time.Duration(global.Cfg().FlushBulkInterval)
		bulkSize := global.Cfg().BulkSize
		ticker := time.NewTicker(time.Millisecond * interval)
		defer ticker.Stop()

		lastSavedTime := time.Now()
		requests := make([]*model.RowRequest, 0, bulkSize)
		var current mysql.Position
		from := s.transferService.CurrentPos
		for {
			needFlush := false
			needSavePos := false
			select {
			case v := <-s.queue:
				//global.PrintEx("handler startListener-->queue-->153-->",v)
				switch v := v.(type) {
				case model.PosRequest:
					now := time.Now()
					if v.Force || now.Sub(lastSavedTime) > 3*time.Second {
						lastSavedTime = now
						needFlush = true
						needSavePos = true
						current = mysql.Position{
							Name: v.Name,
							Pos:  v.Pos,
						}
					}
					//global.PrintEx("handler startListener-->model.PosRequest-->166-->", v)

				case []*model.RowRequest:
					requests = append(requests, v...)
					needFlush = int64(len(requests)) >= global.Cfg().BulkSize
					//global.PrintEx("handler startListener-->model.RowRequest-->171-->", v)
				}
			case <-ticker.C:
				needFlush = true
			case <-s.stop:
				return
			}

			if needFlush && len(requests) > 0 && s.transferService.endpointEnable.Load() {

				err := s.transferService.endpoint.Consume(s.transferService.Name, from, requests, s.transferService.ruleMap)
				if err != nil {
					s.transferService.endpointEnable.Store(false)
					metrics.SetDestState(metrics.DestStateFail)
					logs.Error(err.Error())
					go s.transferService.stopDump()
				}
				requests = requests[0:0]
			}
			if needSavePos && s.transferService.endpointEnable.Load() {
				//logs.Infof("save position %s %d", current.Name, current.Pos)
				//if err := s.transferService.SavePosition(current); err != nil {
				//	logs.Errorf("save sync position %s err %v, close sync", current, err)
				//	s.transferService.Close()
				//	return
				//}
				from = current
			}
		}
	}()
}

func (s *handler) stopListener() {
	log.Println("transfer stop")
	s.stop <- struct{}{}
}
