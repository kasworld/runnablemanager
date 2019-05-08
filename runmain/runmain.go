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

package main

import (
	"context"
	"time"

	"github.com/kasworld/log/loggen/basiclog"
	"github.com/kasworld/runnablemanager"
)

func main() {
	Concurrent := 10
	runningAI := runnablemanager.New(
		"AI",
		Concurrent,
		basiclog.GlobalLogger)

	ctx := context.Background()
	runningAI.Start(ctx)

	startableWorkerIDQueue := make(chan int, Concurrent)
	for i := 0; i < Concurrent; i++ {
		startableWorkerIDQueue <- i
	}

	printTicker := time.NewTicker(time.Second)
	defer printTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return

		case <-printTicker.C:
			basiclog.Debug("current %v", runningAI)

		case ra := <-runningAI.GetEndedCh():
			endedAI, ok := ra.(*AI)
			if !ok {
				basiclog.Fatal("Invalid AI obj %v", ra)
				return
			}
			// start next
			startableWorkerIDQueue <- endedAI.GetRunnableID()

		case workerid := <-startableWorkerIDQueue:
			time.Sleep(time.Millisecond)
			startingAI := NewAI(workerid)
			runningAI.GetWaitStartCh() <- startingAI
		}
	}

}

type AI struct {
	runnableid int
	DoClose    func() `webformhide:"" stringformhide:""`
}

func NewAI(rid int) *AI {
	return &AI{
		runnableid: rid,
		DoClose: func() {
			basiclog.Fatal("Too early DoClose call %v", rid)
		},
	}
}

func (cai *AI) Run(mainctx context.Context) {
	ctx, closeCtx := context.WithCancel(mainctx)
	cai.DoClose = closeCtx
	defer cai.DoClose()

	timerInfoTk := time.NewTicker(time.Second)
	defer timerInfoTk.Stop()
	timerEndTk := time.NewTicker(time.Second * 1)
	defer timerEndTk.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		case <-timerInfoTk.C:

		case <-timerEndTk.C:
			break loop
		}
	}
}

func (cai *AI) GetRunnableID() int {
	return cai.runnableid
}
