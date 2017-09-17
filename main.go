package main

import (
	"flag"
	"os"
	"sync"

	"github.com/RuniVN/find-activation-date/cmd"
	"github.com/Sirupsen/logrus"

	"golang.org/x/sync/syncmap"
)

func main() {
	lf := logrus.Fields{"func": "main"}

	input := flag.String("input", "", "input file location")
	workers := flag.Int("worker", 8, "number of workers")

	flag.Parse()

	if *input == "" {
		logrus.WithFields(lf).Fatal("please input your csv")
	}

	f, err := os.Open(*input)
	if err != nil {
		logrus.WithFields(lf).WithError(err).Fatal("failed to open file")
	}

	logrus.Info("start processing, read by buffer then write to separated files...")
	lFiles, err := cmd.Preprocessing(f)
	if err != nil {
		logrus.Fatal(err)
	}

	// fan out implemented
	logrus.Infof("fan out to %d workers", *workers)
	result := run(workers, lFiles)

	logrus.Info("exporting to csv")
	err = cmd.ExportCSV(result)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("cleaning up...")
	err = cmd.TearDown()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("done.")
}

func run(workers *int, lFiles []string) *syncmap.Map {
	result := &syncmap.Map{}
	var wg sync.WaitGroup

	wg.Add(*workers)
	go pool(&wg, *workers, lFiles, result)
	wg.Wait()

	return result
}

func worker(taskCh <-chan string, wg *sync.WaitGroup, result *syncmap.Map) {
	defer wg.Done()
	for {
		task, ok := <-taskCh
		if !ok {
			return
		}

		r, err := cmd.ProcessOneFile(task)
		if err != nil {
			return
		}
		result.Store(task, r)
	}
}

func pool(wg *sync.WaitGroup, workers int, lFiles []string, result *syncmap.Map) {
	taskCh := make(chan string)

	for i := 0; i < workers; i++ {
		go worker(taskCh, wg, result)
	}

	for i := 0; i < len(lFiles); i++ {
		taskCh <- lFiles[i]
	}

	close(taskCh)
}
