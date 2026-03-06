# 📺 Series Tracker

Aplicación web en Go con SQLite para gestionar series, actualizar episodios y asignar rating.

## 🗂️ Estructura del Proyecto

Lab5_go/
│
├── main.go
├── go.mod
│
├── handlers/
│   └── handlers.go
│
├── db/
│   ├── models.go
│   ├── queries.go
│   └── series.db
│
├── templates/
│   └── templates.go
│
└── static/
    ├── styles.css
    └── script.js
    
## 📁 Descripción de Carpetas y Archivos
**main.go** 

-Punto de entrada del servidor.
  
-Conexión a la base de datos.
  
-Inicializa el servidor en :8080.

**handlers/**

  -Manejo de rutas HTTP.
  
  -Procesa GET y POST.
  
  -Lógica para crear series, actualizar episodios y guardar rating.

**db/**

  -models.go: definición de estructuras (struct Series).
  
  -queries.go: consultas SQL (SELECT con JOIN, INSERT, UPDATE).
  
  -series.db: base de datos SQLite.

**templates/**

  -Generación dinámica del HTML.
  
  -Renderiza la tabla principal y formularios.

**static/**

  styles.css: estilos visuales.
  
  script.js: funciones JavaScript (ej. fetch para +1 episodio).

## 🚀 Ejecución

Desde root:

  -go run .

Abrir en el navegador:

  -http://localhost:8080

## 📸 Cómo se mira
<img width="1409" height="857" alt="Captura de pantalla 2026-03-05 222313" src="https://github.com/user-attachments/assets/a4d5f49d-0284-485f-b0ef-b07d7c4f95bc" />
