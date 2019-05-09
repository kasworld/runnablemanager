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

	"github.com/kasworld/rangestat"
)

type RunnableI interface {
	Run(ctx context.Context)
}

func (raman RunnableManager) String() string {
	return raman.runStat.String()
}

type RunnableManager struct {
	name            string
	concurrentCount int
	waitStartCh     chan RunnableI
	endedCh         chan RunnableI
	runStat         *rangestat.RangeStat
}

func New(name string, concurrentCount int) *RunnableManager {
	raman := &RunnableManager{
		name:            name,
		concurrentCount: concurrentCount,
		waitStartCh:     make(chan RunnableI, concurrentCount),
		endedCh:         make(chan RunnableI, concurrentCount),
		runStat:         rangestat.New(name, 0, concurrentCount),
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
						// fmt.Println("not ok")
						return
					}
					if !raman.runStat.Inc() {
						panic("fail to inc")
					}
					ra.Run(ctx)
					if !raman.runStat.Dec() {
						panic("fail to dec")
					}
					raman.endedCh <- ra
				}
			}
		}()
	}
}

func (raman *RunnableManager) GetWaitStartCh() chan<- RunnableI {
	return raman.waitStartCh
}
func (raman *RunnableManager) GetEndedCh() <-chan RunnableI {
	return raman.endedCh
}
func (raman *RunnableManager) GetRunStat() *rangestat.RangeStat {
	return raman.runStat
}
