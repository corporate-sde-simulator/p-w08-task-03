package main

import "testing"

func TestProcess(t *testing.T) {
	p := NewProcessor()
	_, err := p.Process(map[string]interface{}{"key": "val"})
	if err != nil { t.Fatalf("unexpected error: %v", err) }
}
func TestNilInput(t *testing.T) {
	p := NewProcessor()
	_, err := p.Process(nil)
	if err == nil { t.Error("expected error for nil") }
}
func TestStats(t *testing.T) {
	p := NewProcessor()
	p.Process(map[string]interface{}{"x": 1})
	if p.GetStats()["processed"] != 1 { t.Error("wrong count") }
}
func TestManager(t *testing.T) {
	m := NewManager()
	_, err := m.Process(map[string]interface{}{"data": "test"})
	if err != nil { t.Fatalf("unexpected error: %v", err) }
}
