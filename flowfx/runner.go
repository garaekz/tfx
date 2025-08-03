package flowfx

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/garaekz/tfx/internal/share"
	"github.com/garaekz/tfx/runfx"
)

// FlowRunner manages flow execution lifecycle with proper cancellation and interrupt handling.
type FlowRunner struct {
	flow         Flow
	name         string
	ctx          context.Context
	cancel       context.CancelFunc
	running      bool
	mu           sync.RWMutex
	onStart      Hook
	onComplete   Hook
	onError      Hook
	onStop       Hook
	handleSigInt bool           // Whether to handle system interrupts
	ttyInfo      *runfx.TTYInfo // TTY detection for non-interactive fallback
}

// RunnerConfig provides configuration for a FlowRunner.
type RunnerConfig struct {
	Name         string
	OnStart      Hook
	OnComplete   Hook
	OnError      Hook
	OnStop       Hook
	HandleSigInt bool // Whether to handle Ctrl+C interrupts
}

// DefaultRunnerConfig returns the default configuration for a FlowRunner.
func DefaultRunnerConfig() RunnerConfig {
	return RunnerConfig{
		Name:         "runner",
		HandleSigInt: true,
	}
}

// newFlowRunner creates a new flow runner with the given configuration and flow.
func newFlowRunner(flow Flow, cfg RunnerConfig) *FlowRunner {
	// Detect TTY capabilities for non-interactive fallback
	ttyInfo := runfx.DetectTTY()

	return &FlowRunner{
		flow:         flow,
		name:         cfg.Name,
		onStart:      cfg.OnStart,
		onComplete:   cfg.OnComplete,
		onError:      cfg.OnError,
		onStop:       cfg.OnStop,
		handleSigInt: cfg.HandleSigInt,
		ttyInfo:      &ttyInfo,
	}
}

// --- MULTIPATH API FUNCTIONS ---

// NewRunner creates a new flow runner with multipath configuration support.
// Supports two usage patterns:
//   - NewRunner(flow)                          // Zero-config, uses defaults
//   - NewRunner(flow, config)                  // Config struct
func NewRunner(flow Flow, args ...any) *FlowRunner {
	cfg := share.Overload(args, DefaultRunnerConfig())
	return newFlowRunner(flow, cfg)
}

// NewRunnerBuilder creates a new RunnerBuilder for DSL chaining.
func NewRunnerBuilder() *RunnerBuilder {
	return &RunnerBuilder{config: DefaultRunnerConfig()}
}

// Start begins flow execution with proper context management and interrupt handling.
// It implements the Runner interface.
func (fr *FlowRunner) Start(ctx context.Context) error {
	fr.mu.Lock()
	if fr.running {
		fr.mu.Unlock()
		return NewFlowError(fr.name, "", ErrAlreadyRunning)
	}

	// Create cancellable context
	fr.ctx, fr.cancel = context.WithCancel(ctx)
	fr.running = true
	fr.mu.Unlock()

	// Call onStart hook if provided
	if fr.onStart != nil {
		fr.onStart(fr.ctx, fr.name, nil)
	}

	// Setup interrupt handling if enabled and in TTY environment
	var sigCh chan os.Signal
	if fr.handleSigInt {
		// In non-TTY environments, interrupt handling may be limited
		if !fr.ttyInfo.IsTTY {
			// Use fallback output for non-interactive environments
			runfx.FallbackOutput("Starting flow execution (non-interactive mode)")
		}

		sigCh = make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		// Handle interrupts in a goroutine
		go func() {
			select {
			case <-sigCh:
				if !fr.ttyInfo.IsTTY {
					// Simple output for non-TTY environments
					runfx.FallbackOutput("Flow execution interrupted")
				}
				fr.Stop()
			case <-fr.ctx.Done():
				// Context cancelled, clean up
			}
		}()
	}

	// Execute the flow
	err := fr.flow.Run(fr.ctx)

	// Cleanup
	fr.mu.Lock()
	fr.running = false
	if fr.cancel != nil {
		fr.cancel()
		fr.cancel = nil
	}
	fr.mu.Unlock()

	// Stop interrupt handling
	if sigCh != nil {
		signal.Stop(sigCh)
		close(sigCh)
	}

	// Handle result
	if err != nil {
		if fr.onError != nil {
			fr.onError(fr.ctx, fr.name, err)
		}
		// In non-TTY environments, provide simple error output
		if !fr.ttyInfo.IsTTY {
			runfx.FallbackOutput("Flow execution failed: " + err.Error())
		}
		return err
	}

	// Success
	if fr.onComplete != nil {
		fr.onComplete(fr.ctx, fr.name, nil)
	}

	// In non-TTY environments, provide simple success output
	if !fr.ttyInfo.IsTTY {
		runfx.FallbackOutput("Flow execution completed successfully")
	}

	return nil
}

// Stop cancels the running flow execution.
// It implements the Runner interface.
func (fr *FlowRunner) Stop() error {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	if !fr.running {
		return NewFlowError(fr.name, "", ErrNotRunning)
	}

	if fr.cancel != nil {
		fr.cancel()
	}

	// Call onStop hook if provided
	if fr.onStop != nil {
		fr.onStop(fr.ctx, fr.name, nil)
	}

	return nil
}

// IsRunning returns true if the flow is currently executing.
// It implements the Runner interface.
func (fr *FlowRunner) IsRunning() bool {
	fr.mu.RLock()
	defer fr.mu.RUnlock()
	return fr.running
}

// GetFlow returns the flow being managed by this runner.
func (fr *FlowRunner) GetFlow() Flow {
	return fr.flow
}

// IsInteractive returns true if the runner is operating in an interactive TTY environment.
func (fr *FlowRunner) IsInteractive() bool {
	return fr.ttyInfo.IsTTY
}

// GetTTYInfo returns detailed information about the terminal environment.
func (fr *FlowRunner) GetTTYInfo() runfx.TTYInfo {
	if fr.ttyInfo == nil {
		// Fallback if somehow not initialized
		ttyInfo := runfx.DetectTTY()
		return ttyInfo
	}
	return *fr.ttyInfo
}

// --- DSL BUILDER ---

// RunnerBuilder provides a fluent API for building flow runners.
type RunnerBuilder struct {
	config RunnerConfig
	flow   Flow
}

// Name sets the name of the runner.
func (rb *RunnerBuilder) Name(name string) *RunnerBuilder {
	rb.config.Name = name
	return rb
}

// OnStart sets the start hook.
func (rb *RunnerBuilder) OnStart(hook Hook) *RunnerBuilder {
	rb.config.OnStart = hook
	return rb
}

// OnComplete sets the complete hook.
func (rb *RunnerBuilder) OnComplete(hook Hook) *RunnerBuilder {
	rb.config.OnComplete = hook
	return rb
}

// OnError sets the error hook.
func (rb *RunnerBuilder) OnError(hook Hook) *RunnerBuilder {
	rb.config.OnError = hook
	return rb
}

// OnStop sets the stop hook.
func (rb *RunnerBuilder) OnStop(hook Hook) *RunnerBuilder {
	rb.config.OnStop = hook
	return rb
}

// HandleSigInt enables or disables system interrupt handling.
func (rb *RunnerBuilder) HandleSigInt(enabled bool) *RunnerBuilder {
	rb.config.HandleSigInt = enabled
	return rb
}

// Flow sets the flow to be managed by the runner.
func (rb *RunnerBuilder) Flow(flow Flow) *RunnerBuilder {
	rb.flow = flow
	return rb
}

// Build creates a new FlowRunner instance.
func (rb *RunnerBuilder) Build() *FlowRunner {
	return newFlowRunner(rb.flow, rb.config)
}

// Start creates and starts the flow runner.
func (rb *RunnerBuilder) Start(ctx context.Context) error {
	return rb.Build().Start(ctx)
}

// --- CONVENIENCE FUNCTIONS ---

// RunFlow is a convenience function to run a flow with default runner settings.
func RunFlow(ctx context.Context, flow Flow) error {
	runner := NewRunner(flow)
	return runner.Start(ctx)
}

// RunFlowWithInterrupts is a convenience function to run a flow with interrupt handling.
func RunFlowWithInterrupts(ctx context.Context, flow Flow) error {
	config := RunnerConfig{
		Name:         "flow-runner",
		HandleSigInt: true,
	}
	runner := NewRunner(flow, config)
	return runner.Start(ctx)
}

// RunFlowNonInteractive is a convenience function to run a flow in non-interactive mode.
// This is particularly useful in CI/CD environments, automated scripts, or when stdout/stderr are redirected.
func RunFlowNonInteractive(ctx context.Context, flow Flow) error {
	// Detect TTY to automatically configure for non-interactive environments
	ttyInfo := runfx.DetectTTY()

	config := RunnerConfig{
		Name:         "non-interactive-runner",
		HandleSigInt: ttyInfo.IsTTY, // Only handle interrupts if we're in a TTY
	}

	runner := NewRunner(flow, config)

	// Provide feedback for non-TTY environments
	if !ttyInfo.IsTTY {
		runfx.FallbackOutput("Executing flow in non-interactive mode...")
	}

	return runner.Start(ctx)
}

// --- ENVIRONMENT DETECTION ---

// IsInteractiveEnvironment returns true if running in an interactive TTY environment.
// This is useful for flows that need to adapt their behavior based on the execution context.
func IsInteractiveEnvironment() bool {
	ttyInfo := runfx.DetectTTY()
	return ttyInfo.IsTTY
}

// GetEnvironmentInfo returns detailed information about the current execution environment.
func GetEnvironmentInfo() runfx.TTYInfo {
	return runfx.DetectTTY()
}
