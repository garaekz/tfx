package runfx

import (
	"bytes"
	"testing"
	"time"
)

func TestMultipathAPI(t *testing.T) {
	tests := []struct {
		name           string
		createLoop     func() Loop
		expectedTick   time.Duration
		expectTestMode bool
	}{
		{
			name: "Beginner - zero config",
			createLoop: func() Loop {
				return Start()
			},
			expectedTick:   50 * time.Millisecond, // default
			expectTestMode: false,
		},
		{
			name: "Beginner - config struct",
			createLoop: func() Loop {
				cfg := Config{
					TickInterval: 100 * time.Millisecond,
					TestMode:     true,
				}
				return Start(cfg)
			},
			expectedTick:   100 * time.Millisecond,
			expectTestMode: true,
		},
		{
			name: "Hardcore - DSL builder",
			createLoop: func() Loop {
				return New().
					TickInterval(75 * time.Millisecond).
					TestMode().
					Start()
			},
			expectedTick:   75 * time.Millisecond,
			expectTestMode: true,
		},
		{
			name: "Hardcore - smooth animation",
			createLoop: func() Loop {
				return New().
					SmoothAnimation().
					TestMode().
					Start()
			},
			expectedTick:   30 * time.Millisecond,
			expectTestMode: true,
		},
		{
			name: "Hardcore - fast animation",
			createLoop: func() Loop {
				return New().
					FastAnimation().
					TestMode().
					Start()
			},
			expectedTick:   100 * time.Millisecond,
			expectTestMode: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loop := tt.createLoop()

			// Cast to MainLoop to access internal fields for testing
			mainLoop, ok := loop.(*MainLoop)
			if !ok {
				t.Fatal("Expected MainLoop instance")
			}

			// Check tick interval
			if mainLoop.eventLoop.interval != tt.expectedTick {
				t.Errorf("Expected tick interval %v, got %v",
					tt.expectedTick, mainLoop.eventLoop.interval)
			}

			// Check test mode
			if mainLoop.testMode != tt.expectTestMode {
				t.Errorf("Expected test mode %v, got %v",
					tt.expectTestMode, mainLoop.testMode)
			}
		})
	}
}

func TestDSLBuilderChaining(t *testing.T) {
	var buf bytes.Buffer

	loop := New().
		TickInterval(80 * time.Millisecond).
		Output(&buf).
		TestMode().
		Start()

	mainLoop, ok := loop.(*MainLoop)
	if !ok {
		t.Fatal("Expected MainLoop instance")
	}

	// Verify all configurations were applied
	if mainLoop.eventLoop.interval != 80*time.Millisecond {
		t.Errorf("Expected tick interval 80ms, got %v", mainLoop.eventLoop.interval)
	}

	if mainLoop.output != &buf {
		t.Error("Expected custom output writer")
	}

	if !mainLoop.testMode {
		t.Error("Expected test mode to be enabled")
	}
}

func TestOverloadAPI(t *testing.T) {
	// Test zero-config
	loop1 := Start()
	mainLoop1 := loop1.(*MainLoop)
	if mainLoop1.eventLoop.interval != 50*time.Millisecond {
		t.Errorf("Zero-config: expected 50ms, got %v", mainLoop1.eventLoop.interval)
	}

	// Test with config struct
	cfg := Config{
		TickInterval: 200 * time.Millisecond,
		TestMode:     true,
	}
	loop2 := Start(cfg)
	mainLoop2 := loop2.(*MainLoop)
	if mainLoop2.eventLoop.interval != 200*time.Millisecond {
		t.Errorf("Config struct: expected 200ms, got %v", mainLoop2.eventLoop.interval)
	}
	if !mainLoop2.testMode {
		t.Error("Config struct: expected test mode enabled")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	expected := Config{
		TickInterval: 50 * time.Millisecond,
		TestMode:     false,
	}

	if cfg.TickInterval != expected.TickInterval {
		t.Errorf("Expected tick interval %v, got %v",
			expected.TickInterval, cfg.TickInterval)
	}

	if cfg.TestMode != expected.TestMode {
		t.Errorf("Expected test mode %v, got %v",
			expected.TestMode, cfg.TestMode)
	}

	if cfg.Output == nil {
		t.Error("Expected default output to be set")
	}
}

func TestFunctionalOptions(t *testing.T) {
	cfg := DefaultConfig()

	// Test WithTickInterval
	WithTickInterval(200 * time.Millisecond)(&cfg)
	if cfg.TickInterval != 200*time.Millisecond {
		t.Errorf("WithTickInterval failed: expected 200ms, got %v", cfg.TickInterval)
	}

	// Test WithTestMode
	WithTestMode()(&cfg)
	if !cfg.TestMode {
		t.Error("WithTestMode failed: expected true")
	}

	// Test WithSmoothAnimation
	WithSmoothAnimation()(&cfg)
	if cfg.TickInterval != 30*time.Millisecond {
		t.Errorf("WithSmoothAnimation failed: expected 30ms, got %v", cfg.TickInterval)
	}

	// Test WithFastAnimation
	WithFastAnimation()(&cfg)
	if cfg.TickInterval != 100*time.Millisecond {
		t.Errorf("WithFastAnimation failed: expected 100ms, got %v", cfg.TickInterval)
	}
}
