package shutdown

// Resource Manager — tracks and manages service resources (DB pools, caches, etc.)
//
// Author: Vikram Patel (Infra team)
// Last Modified: 2026-03-25

import (
	"fmt"
	"sync"
)

type ResourceType string

const (
	DatabasePool ResourceType = "database_pool"
	CacheLayer   ResourceType = "cache"
	MessageQueue ResourceType = "message_queue"
	FileHandle   ResourceType = "file_handle"
)

type Resource struct {
	Name     string
	Type     ResourceType
	Active   bool
	Cleanup  func() error
}

type ResourceManager struct {
	mu        sync.RWMutex
	resources map[string]*Resource
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		resources: make(map[string]*Resource),
	}
}

func (rm *ResourceManager) Register(name string, resType ResourceType, cleanup func() error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.resources[name] = &Resource{
		Name:    name,
		Type:    resType,
		Active:  true,
		Cleanup: cleanup,
	}
}

func (rm *ResourceManager) CloseResource(name string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	res, ok := rm.resources[name]
	if !ok {
		return fmt.Errorf("resource %s not found", name)
	}
	if !res.Active {
		return nil
	}

	if res.Cleanup != nil {
		if err := res.Cleanup(); err != nil {
			return fmt.Errorf("failed to close %s: %w", name, err)
		}
	}
	res.Active = false
	return nil
}

func (rm *ResourceManager) CloseAll() map[string]error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	errors := make(map[string]error)
	for name, res := range rm.resources {
		if res.Active && res.Cleanup != nil {
			if err := res.Cleanup(); err != nil {
				errors[name] = err
			} else {
				res.Active = false
			}
		}
	}
	return errors
}

func (rm *ResourceManager) GetActiveCount() int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	count := 0
	for _, res := range rm.resources {
		if res.Active {
			count++
		}
	}
	return count
}
