package shutdown

// Shutdown Orchestrator — manages graceful service shutdown.
//
// Executes cleanup tasks (drain connections, flush caches, close DB pools)
// in dependency order with a global timeout.
//
// Author: Vikram Patel (Infra team)
// Last Modified: 2026-03-25

import (
	"fmt"
	"sync"
	"time"
)

type ShutdownTask struct {
	Name         string
	Handler      func() error
	DependsOn    []string // Other tasks that must complete first
	Timeout      time.Duration
}

type TaskResult struct {
	Name       string
	Success    bool
	Duration   time.Duration
	Error      error
}

type ShutdownOrchestrator struct {
	mu             sync.Mutex
	tasks          []ShutdownTask
	results        []TaskResult
	globalTimeout  time.Duration
	shutdownStart  time.Time
}

func NewShutdownOrchestrator(globalTimeout time.Duration) *ShutdownOrchestrator {
	return &ShutdownOrchestrator{
		tasks:         make([]ShutdownTask, 0),
		results:       make([]TaskResult, 0),
		globalTimeout: globalTimeout,
	}
}

func (so *ShutdownOrchestrator) RegisterTask(task ShutdownTask) {
	so.mu.Lock()
	defer so.mu.Unlock()
	so.tasks = append(so.tasks, task)
}

func (so *ShutdownOrchestrator) Shutdown() []TaskResult {
	so.mu.Lock()
	defer so.mu.Unlock()

	// Build dependency-ordered execution plan
	ordered := so.buildExecutionOrder()

	// ordering. Tasks that depend on each other run simultaneously.
	// before starting it. You can check completed tasks in a set.
	var wg sync.WaitGroup
	for _, task := range ordered {
		wg.Add(1)
		go func(t ShutdownTask) {
			defer wg.Done()
			result := so.executeTask(t)
			so.results = append(so.results, result)
		}(task)
	}
	wg.Wait()

	// This means the timeout never triggers during actual execution.
	// Then check elapsed time before/during each task execution.
	so.shutdownStart = time.Now()
	deadline := so.shutdownStart.Add(so.globalTimeout)
	if time.Now().After(deadline) {
		fmt.Println("WARNING: Shutdown exceeded global timeout")
	}

	return so.results
}

func (so *ShutdownOrchestrator) executeTask(task ShutdownTask) TaskResult {
	start := time.Now()

	done := make(chan error, 1)
	go func() {
		done <- task.Handler()
	}()

	select {
	case err := <-done:
		return TaskResult{
			Name:     task.Name,
			Success:  err == nil,
			Duration: time.Since(start),
			Error:    err,
		}
	case <-time.After(task.Timeout):
		return TaskResult{
			Name:     task.Name,
			Success:  false,
			Duration: time.Since(start),
			Error:    fmt.Errorf("task %s timed out after %v", task.Name, task.Timeout),
		}
	}
}

func (so *ShutdownOrchestrator) buildExecutionOrder() []ShutdownTask {
	// Simple topological sort (BFS)
	depCount := make(map[string]int)
	taskMap := make(map[string]ShutdownTask)
	dependents := make(map[string][]string)

	for _, t := range so.tasks {
		taskMap[t.Name] = t
		depCount[t.Name] = len(t.DependsOn)
		for _, dep := range t.DependsOn {
			dependents[dep] = append(dependents[dep], t.Name)
		}
	}

	var ordered []ShutdownTask
	queue := make([]string, 0)

	for name, count := range depCount {
		if count == 0 {
			queue = append(queue, name)
		}
	}

	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		ordered = append(ordered, taskMap[name])

		for _, dep := range dependents[name] {
			depCount[dep]--
			if depCount[dep] == 0 {
				queue = append(queue, dep)
			}
		}
	}

	return ordered
}

func (so *ShutdownOrchestrator) GetResults() []TaskResult {
	return so.results
}
