package main

import (
	"context"
	"fmt"
	"time"

	"github.com/garaekz/tfx/progress" // Importa tu paquete de progreso
	"github.com/garaekz/tfx/runfx"    // Importa tu runner
)

func runFormFXDemo() {
	// 1. Crear el loop de RunFX.
	// Usamos la vía Express para simplicidad.
	loop := runfx.Start()

	// 2. Crear el componente de barra de progreso.
	// También usamos la vía Express. Le damos un total de 100.
	progressBar := progress.Start(progress.ProgressConfig{
		Total: 100,
		Label: "Descargando archivos...",
	})

	// 3. Montar el componente en el loop.
	// runfx se encargará de renderizarlo en cada tick.
	unmount, err := loop.Mount(progressBar)
	if err != nil {
		fmt.Printf("Error al montar la barra de progreso: %v\n", err)
		return
	}
	// Nos aseguramos de desmontar el componente al final.
	defer unmount()

	// 4. Ejecutar el loop en una goroutine (en segundo plano).
	// Usamos un contexto para poder detenerlo más tarde.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Llama a cancel() cuando main() termina para limpiar la goroutine.

	go func() {
		fmt.Println("Iniciando el loop de renderizado en segundo plano...")
		if err := loop.Run(ctx); err != nil {
			// Esto se imprimirá si el loop se detiene con un error,
			// lo cual es normal si el contexto se cancela.
			// fmt.Printf("Loop terminado con error: %v\n", err)
		}
		fmt.Println("Loop de renderizado detenido.")
	}()

	// 5. Simular trabajo en la goroutine principal.
	// Aquí es donde estaría tu lógica de negocio (descargar un archivo, procesar datos, etc.).
	fmt.Println("Iniciando tarea principal...")
	for i := 0; i <= 100; i++ {
		// Actualizamos el estado del componente de progreso.
		// El loop de runfx detectará el cambio y lo re-renderizará automáticamente.
		progressBar.Set(i)
		time.Sleep(30 * time.Millisecond) // Simula una pequeña cantidad de trabajo.
	}

	// 6. Finalizar la barra de progreso y detener el loop.
	progressBar.Finish()

	// Es importante dar un pequeño margen para que el loop renderice el estado final (100%).
	time.Sleep(50 * time.Millisecond)

	// La llamada a `cancel()` (a través del `defer`) detendrá el loop.
	fmt.Println("\nTarea completada.")
}
