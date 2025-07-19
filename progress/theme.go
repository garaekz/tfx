package progress

// ProgressTheme define colores para la barra y label
// NOTA: puedes integrar con tu sistema de color real más adelante

type ProgressTheme struct {
	CompleteColor   string
	IncompleteColor string
	LabelColor      string
}

var DraculaTheme = ProgressTheme{
	CompleteColor:   "\033[38;5;49m",  // verde
	IncompleteColor: "\033[38;5;240m", // gris
	LabelColor:      "\033[38;5;111m", // azul
}
