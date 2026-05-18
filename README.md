# PR3PL (Modernized in Go)

PR3PL es un lenguaje de programación de propósito académico y paradigma funcional. Este proyecto es una refactorización arquitectónica que dota al motor base original (escrito en C++) de una sintaxis moderna, legible y ergonómica, construyendo un pipeline de compilación completo (Frontend) y un intérprete nativo en Go.

Actualmente, el proyecto implementa desde la lectura de caracteres hasta la optimización del Árbol de Sintaxis Abstracta (AST) y la transpilación final.

## 🚀 Características Principales (Pipeline de Compilación)

* **Analizador Léxico Concurrente:** Implementado mediante una Máquina de Estados Finitos (arquitectura propuesta por Rob Pike). Opera como productor en una tubería de concurrencia (Channels y Goroutines).
* **Analizador Sintáctico (Pratt Parser):** Un parser Top-Down de descenso recursivo capaz de resolver matemáticamente la precedencia de operadores lógicos y aritméticos.
* **Analizador Semántico (Type Checker):** Motor de inferencia y validación de tipos que ataja errores lógicos (ej. sumar booleanos con enteros o aplicar parámetros incorrectos a closures) *antes* de permitir la ejecución.
* **Optimizador de AST (Constant Folding):** Fase de optimización intermedia que evalúa y pliega operaciones matemáticas estáticas en tiempo de compilación para  reducir la huella del código transpilado.
* **Evaluador Nativo:** Intérprete interactivo que ejecuta el código directamente en la memoria de Go.
* **Transpilador con :** Genera la traducción inversa hacia el formato Lisp puro requerido por el backend original en C++, aplicando un formateo de indentación cascada para mantener el código resultante (`transpiled.txt`) legible para humanos.

## ✨ Sintaxis Modernizada

La principal razón de este proyecto es eliminar la notación matemática rígida (Lisp-like) del motor original, permitiendo una escritura mucho más humana y limpia, similar a la familia de lenguajes ML o Haskell:

**1. Operaciones Matemáticas y Lógicas**
* *Clásico (C++):* `(add (5) (mult (7) (2)))`
* *Moderno (Go):* `5 + 7 * 2`

**2. Estructuras de Datos (Pares y Listas)**
* *Clásico (C++):* `(pair (1) (pair (2) (pair (3) (unit))))`
* *Moderno (Go):* `[1, 2, 3]`

**3. Variables Globales**
* *Clásico (C++):* `(val (miNumero) (mult (10) (5)))`
* *Moderno (Go):* `val miNumero = 10 * 5`

**4. Condicionales**
* *Clásico (C++):* `(if (lt (x) (y)) (100) (200))`
* *Moderno (Go):* `if x < y then 100 else 200`

**5. Funciones y Bloques Locales**
* *Clásico (C++):* `(let (pow) (fun (cuadrado) (x) (mult (x) (x))) (call (pow) (8)))`
* *Moderno (Go):* `let pow = fun cuadrado(x) = x * x in pow(8) end`

## 🛠️ Cómo ejecutar el proyecto

Asegúrate de tener [Go](https://golang.org/dl/) instalado en tu sistema.

1. **Clonar el repositorio:**
   ```bash
   git clone [https://github.com/tu_usuario/pr3pl-go.git](https://github.com/tu_usuario/pr3pl-go.git)
   cd pr3pl-go

2. **Ejecutar un archivo fuente (Intérprete y Transpilador):**
   ```bash
   go run main.go examples/...
   
3. **Iniciar la consola interactiva (REPL):**
   ```bash
   go run main.go

4. **Iniciar la consola interactiva (REPL):**
    ```bash
   go test ./...
