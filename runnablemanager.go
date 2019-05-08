// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package runnablemanager manage limit concurrent running
package runnablemanager

import (
	"context"
	"fmt"

	"github.com/kasworld/rangestat"
)

type RunnableI interface {
	Run(ctx context.Context)
	GetRunnableID() int
}

func (raman RunnableManager) String() string {
	return fmt.Sprintf("RunnableManager[%v %v]",
		raman.name,
		raman.runStat)
}

type RunnableManager struct {
	name            string
	concurrentCount int
	waitStartCh     chan RunnableI
	runnableList    []RunnableI
	endedCh         chan RunnableI
	runStat         *rangestat.RangeStat
	logger          loggerI
}

func New(name string, concurrentCount int, logger loggerI) *RunnableManager {
	raman := &RunnableManager{
		name:            name,
		concurrentCount: concurrentCount,
		waitStartCh:     make(chan RunnableI, concurrentCount),
		endedCh:         make(chan RunnableI, concurrentCount),
		runnableList:    make([]RunnableI, concurrentCount),
		runStat:         rangestat.New(name, 0, concurrentCount),
		logger:          logger,
	}
	return raman
}

func (raman *RunnableManager) Start(ctx context.Context) {
	for i := 0; i < raman.concurrentCount; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case ra, ok := <-raman.waitStartCh:
					if !ok {
						return
					}
					if !raman.runStat.Inc() {
						raman.logger.Fatal("fail to inc")
					}
					raman.runnableList[ra.GetRunnableID()] = ra
					ra.Run(ctx)
					if !raman.runStat.Dec() {
						raman.logger.Fatal("fail to dec")
					}
					raman.runnableList[ra.GetRunnableID()] = nil
					raman.endedCh <- ra
				}
			}
		}()
	}
}
