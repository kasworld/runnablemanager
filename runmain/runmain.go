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
	"fmt"
	"time"

	"github.com/kasworld/runnablemanager"
)

func main() {
	runnableCount := 10
	jobCount := runnableCount * 10
	runningAI := runnablemanager.New("AI", runnableCount)

	startableJobIDQueue := make(chan int, jobCount)
	for i := 0; i < jobCount; i++ {
		startableJobIDQueue <- i
	}

	ctx, ctxEnd := context.WithCancel(context.Background())

	runningAI.Start(ctx)

	printTicker := time.NewTicker(time.Second)
	defer printTicker.Stop()

	go func() {
	loop1:
		for {
			select {
			case <-ctx.Done():
				break loop1
			case jobID := <-startableJobIDQueue:
				runningAI.GetWaitStartCh() <- NewAI(jobID)
			}
		}
	}()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		case <-printTicker.C:
			fmt.Printf("current %v\n", runningAI)

		case ra := <-runningAI.GetEndedCh():
			endedAI, ok := ra.(*AI)
			if !ok {
				panic(fmt.Sprintf("Invalid AI obj %v", ra))
				break loop
			}
			fmt.Printf("job ended %v\n", endedAI)
			// // enqueue
			// startableJobIDQueue <- endedAI.jobID
			if runningAI.GetRunStat().GetCurrentVal() == 0 {
				break loop
			}
		}
	}
	ctxEnd()
	fmt.Printf("current %v\n", runningAI)
}

type AI struct {
	jobID   int
	DoClose func() `webformhide:"" stringformhide:""`
}

func NewAI(wid int) *AI {
	return &AI{
		jobID: wid,
		DoClose: func() {
			panic(fmt.Sprintf("Too early DoClose call %v", wid))
		},
	}
}

func (cai AI) String() string {
	return fmt.Sprintf("AI-%v", cai.jobID)
}

func (cai *AI) Run(mainctx context.Context) {
	time.Sleep(time.Second * 1)
}
