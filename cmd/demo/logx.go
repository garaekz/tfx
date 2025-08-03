package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/garaekz/tfx/color"
	"github.com/garaekz/tfx/logx"
)

func runLogxDemo() {
	// 🎬 Opening Cinematic
	fmt.Println("\n🌟 ═══════════════════════════════════════════════════════════════")
	fmt.Println("🎨                   TFX LOGX SHOWCASE                       🎨")
	fmt.Println("🌟 ═══════════════════════════════════════════════════════════════")
	fmt.Println("🚀 Building the future of terminal logging, one badge at a time")
	fmt.Println()

	// 🎭 Chapter 1: The Art of Visual Communication
	fmt.Println("\n🎭 Chapter 1: THE ART OF VISUAL COMMUNICATION")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	logx.Trace("🔍 Deep system introspection - following the electron trails")
	logx.Debug("🛠️  Engineering insights - watching the gears turn")
	logx.Info("ℹ️  Narrative unfolding - the story your system tells")
	logx.Success("✨ Victory achieved - dreams becoming reality")
	logx.Warn("⚠️  Storm approaching - wisdom from the edge")
	logx.Error("🔥 Phoenix moment - rising from digital ashes")

	time.Sleep(400 * time.Millisecond)

	// 🏆 Chapter 2: Badge Mastery - Where Form Meets Function
	fmt.Println("\n🏆 Chapter 2: BADGE MASTERY - WHERE FORM MEETS FUNCTION")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	logx.SuccessBadge("API", "🌐 Quantum entanglement with external dimensions")
	logx.ErrorBadge("DB", "🗄️  Digital ocean levels critical - send backup whales")
	logx.WarnBadge("CACHE", "⚡ Memory palace reorganizing - philosophers confused")
	logx.InfoBadge("SYS", "🔄 Cosmic alignment achieved - servers humming in harmony")
	logx.DebugBadge("AUTH", "🔐 Digital handshakes verified - trust protocols engaged")

	// ✨ Chapter 2.5: Visual Badge Magic
	fmt.Println("\n✨ Visual Badge Magic:")
	logx.BadgeWithOptions("DEPLOY", "🚀 Rocket ship departing for production", logx.BadgeOptions{
		Gradient: []color.Color{color.NewHex("D38312"), color.NewHex("A83279")},
	})
	logx.BadgeWithOptions("NEON", "💎 Cyberpunk dreams materializing", logx.BadgeOptions{
		Neon: true,
	})
	logx.BadgeWithOptions("THEME", "🎨 Canvas painted with midnight blues", logx.BadgeOptions{
		Theme: "blue",
		Bold:  true,
	})
	logx.BadgeWithOptions("PULSE", "💓 Digital heartbeat detected", logx.BadgeOptions{
		Blink:      true,
		Foreground: color.NewHex("FF6B6B"),
	})
	logx.BadgeWithOptions("EPIC", "🌈 Where typography becomes art", logx.BadgeOptions{
		Bold:       true,
		Italic:     true,
		Underline:  true,
		Foreground: color.NewHex("4ECDC4"),
	})

	time.Sleep(500 * time.Millisecond)

	// 🌈 Chapter 3: The Startup Symphony
	fmt.Println("\n🌈 Chapter 3: THE STARTUP SYMPHONY")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	logx.APIBadge("🌐 GraphQL mutations dancing through fiber optics", true)
	logx.DatabaseBadge("🗄️  MongoDB clusters singing in perfect harmony", true)
	logx.AuthBadge("🔐 OAuth2 tokens born from digital stardust", true)
	logx.CacheBadge("⚡ Redis pipelines conducting lightning symphonies", true)
	logx.SystemBadge("🔄 Kubernetes orchestrating the cloud ballet")
	logx.InfoBadge("STARTUP", "💡 Where unicorns meet terminal rainbows")
	logx.WarnBadge("GROWTH", "📈 Scaling faster than coffee consumption")
	logx.ErrorBadge("HUSTLE", "💪 Failing fast, learning faster")

	time.Sleep(400 * time.Millisecond)

	// 🔮 Chapter 4: The Wisdom of Conditional Logic
	fmt.Println("\n🔮 Chapter 4: THE WISDOM OF CONDITIONAL LOGIC")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	networkErr := errors.New("cosmic interference detected")
	dataErr := errors.New("reality validation failure")
	var noErr error

	if logx.ErrorIf(networkErr, "🌌 Network spirits are restless tonight") {
		fmt.Println("   ✨ Error captured in digital amber")
	}

	if logx.WarnIf(dataErr, "⚠️  The data oracle speaks of inconsistencies") {
		fmt.Println("   ✨ Warning whispered to the terminal winds")
	}

	if !logx.InfoIf(noErr, "This won't be logged") {
		fmt.Println("   ✨ Silence is golden - no error, no noise")
	}

	time.Sleep(400 * time.Millisecond)

	// 🏛️ Chapter 5: Enterprise Theater
	fmt.Println("\n🏛️ Chapter 5: ENTERPRISE THEATER")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	logx.DatabaseBadge("🛒 Order #12345 crystallized in digital vaults", true)
	logx.APIBadge("💳 Payment electrons successfully transferred", true)
	logx.CacheBadge("📦 Product catalog materialized in memory palace", true)
	logx.DatabaseBadge("🌪️  User preferences scattered by digital winds", false)
	logx.APIBadge("🚫 External service vanished into the void", false)

	time.Sleep(400 * time.Millisecond)

	// 🌊 Chapter 6: The Context Rivers
	fmt.Println("\n🌊 Chapter 6: THE CONTEXT RIVERS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	ctx := logx.WithFields(map[string]any{
		"request_id": "req-cosmic-123",
		"user_id":    "user-stardust-789",
		"operation":  "soul_update",
		"region":     "multiverse-west-∞",
		"dimension":  "digital",
	})

	ctx.Info("🚀 API request born from quantum thoughts")
	ctx.Success("✨ User essence successfully transformed")

	ctxWithFile := ctx.WithField("file_size", "2.3MB of pure magic")
	ctxWithFile.Info("📁 Digital artifact upload commenced")

	uploadErr := errors.New("artifact too powerful for this realm")
	if ctxWithFile.ErrorIf(uploadErr, "📤 Artifact rejected by reality") {
		ctxWithFile.Warn("🔄 Compressing magic for mortal consumption")
	}

	time.Sleep(400 * time.Millisecond)

	// ⚡ Chapter 7: The Three Paths of Creation
	fmt.Println("\n⚡ Chapter 7: THE THREE PATHS OF CREATION")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	expressLogger := logx.Log()
	expressLogger.Info("🏃 Express path - swift as lightning")

	configLogger := logx.LogWithConfig(logx.LogOptions{
		Level:     logx.LevelDebug,
		Timestamp: true,
	})
	configLogger.Debug("🏗️  Config path - engineered with precision")

	fluentLogger := logx.LogWith(
		logx.WithLevel(logx.LevelInfo),
		logx.WithTimestamp(false),
		logx.WithDevelopment(),
	)
	fluentLogger.Success("🌊 Fluent path - flowing like liquid poetry")

	time.Sleep(400 * time.Millisecond)

	// 🎪 Chapter 8: The Logger Personalities
	fmt.Println("\n🎪 Chapter 8: THE LOGGER PERSONALITIES")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	devLogger := logx.DevLogger()
	devLogger.Debug("🛠️  Development oracle - seeing through code veils")

	consoleLogger := logx.ConsoleLogger()
	consoleLogger.Info("🖥️  Console poet - painting words on terminal canvas")

	structuredLogger := logx.StructuredLogger()
	structuredLogger.Info("📋 Structured sage - speaking in machine tongues")

	time.Sleep(400 * time.Millisecond)

	// 🌀 Chapter 9: The Fluent Spells
	fmt.Println("\n🌀 Chapter 9: THE FLUENT SPELLS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	testErr := errors.New("cosmic service disruption")

	logx.If(testErr).AsError().Msg("🌌 Service portal sealed by interdimensional forces")
	logx.If(testErr).
		AsWarn().
		WithField("retry_count", 3).
		Msg("🔄 Attempting quantum tunnel reconnection")

	logx.If(testErr).
		AsError().
		WithField("service", "soul-auth").
		WithField("timeout", "30s of eternity").
		WithField("attempt", "first of many").
		Msg("🔐 Authentication temple temporarily closed")

	time.Sleep(400 * time.Millisecond)

	// 💫 Chapter 10: The Performance Symphony
	fmt.Println("\n💫 Chapter 10: THE PERFORMANCE SYMPHONY")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	logx.SystemBadge("💾 Memory rivers flowing at 234MB - crystal clear")
	logx.SystemBadge("🖥️  CPU dancing at 12% - energy in perfect harmony")
	logx.SystemBadge("🌐 Network whispers at 45ms - messages swift as wind")
	logx.SystemBadge("⚡ Response lightning at 120ms - speed of thought")

	// 🎆 The Grand Finale
	fmt.Println("\n🎆 ═══════════════════════════════════════════════════════════════")
	fmt.Println("🎊                    GRAND FINALE                             🎊")
	fmt.Println("🎆 ═══════════════════════════════════════════════════════════════")

	// Final showcase of mastery
	logx.InfoBadge("FINALE", "🎭 All systems awakening from digital dreams")
	logx.WarnBadge("MEMORY", "💭 Cache poets writing verses in silicon")
	logx.ErrorBadge("PHOENIX", "🔥 Errors transforming into wisdom")

	time.Sleep(500 * time.Millisecond)

	// 🌟 The Epic Conclusion
	fmt.Println("\n🌟 ═══════════════════════════════════════════════════════════════")
	fmt.Println("✨                  THE LOGX MANIFESTO                         ✨")
	fmt.Println("🌟 ═══════════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("🎨 ARTISTRY ACHIEVED:")
	fmt.Println("   • 🌈 Visual poetry written in terminal light")
	fmt.Println("   • 🏆 Badge mastery - where function wears beauty")
	fmt.Println("   • 🎯 Context flows like rivers through digital landscapes")
	fmt.Println("   • 🌊 Fluent APIs that speak in developer dreams")
	fmt.Println("   • ⚡ Three paths of creation - choose your journey")
	fmt.Println("   • 📊 Performance metrics that sing of optimization")
	fmt.Println("   • 🎭 Personalities for every logging soul")
	fmt.Println("")
	fmt.Println("🚀 WELCOME TO THE FUTURE OF TERMINAL EXPRESSION")
	fmt.Println("💫 Where every log entry is a work of art")
	fmt.Println("🌟 Where developers become digital poets")
	fmt.Println("✨ Where terminals transform into galleries")
	fmt.Println("")
	fmt.Println("🎉 TFX LogX - Painting dreams in terminal light since 2025")
}
