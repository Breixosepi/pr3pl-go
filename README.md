# PR3PL (Modernized in Go)

PR3PL es un lenguaje de programación de propósito académico y paradigma funcional. Este proyecto es una refactorización arquitectónica que dota al motor base original (escrito en C++) de una sintaxis moderna, legible y ergonómica, separando la sintaxis concreta de la sintaxis abstracta.

Actualmente, el proyecto implementa Analizador Léxico y Analizador Sintáctico construida en Go.

## 🚀 Características Principales

* **Analizador Léxico Concurrente:** Implementado mediante una Máquina de Estados Finitos que opera como productor en una tubería de concurrencia (Channels y Goroutines), emitiendo tokens bajo demanda para optimizar el uso de memoria.
* **Analizador Sintáctico (Pratt Parser):** Un parser Top-Down de descenso recursivo capaz de resolver matemáticamente la precedencia de operadores lógicos y aritméticos sin saturar la pila de llamadas.
* **Árbol de Sintaxis Abstracta (AST):** Reconstrucción jerárquica del código fuente mediante interfaces.
* **Soporte REPL Interactivo:** Consola interactiva con soporte para bloques lógicos multilínea.

## ✨ Sintaxis Modernizada

La principal razon de este proyecto es eliminar la notación matemática rígida (Lisp-like) requerida por el motor original, permitiendo una escritura mucho más humana, limpia y parecida a lenguajes funcionales modernos (como ML o Haskell):

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

1.  **Clonar el repositorio:**
    ```bash
    git clone [https://github.com/tu_usuario/pr3pl-go.git](https://github.com/tu_usuario/pr3pl-go.git)
    cd pr3pl-go
    ```

2.  **Iniciar la consola interactiva (REPL):**
    ```bash
    go run main.go
    ```
    *(Escribe código de PR3PL directamente en la consola para ver cómo el Parser construye el Árbol AST. Presiona `Ctrl+C` para salir).*

3.  **Ejecutar las pruebas unitarias:**
    ```bash
    go test ./...
    ```
