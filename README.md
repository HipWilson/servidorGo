# Series Tracker — Lab 5

Servidor HTTP hecho desde cero en Go usando TCP, sin usar el paquete `net/http`. Permite registrar y hacer seguimiento de series de televisión, con una base de datos SQLite.

<img width="1919" height="1024" alt="image" src="https://github.com/user-attachments/assets/58edceee-36a2-42a2-a443-a1438a3fd705" />


## Como correr el proyecto

```bash
go run .
```

Luego abrir el navegador en: `http://localhost:8080`


## Estructura de archivos

```
servidor-tcp-go/
├── main.go         
├── handlers.go     
├── db.go           
├── templates.go    
├── static.go       
├── static/
│   ├── style.css   
│   ├── app.js     
│   └── favicon.svg 
├── series.db       
└── README.md
```

## Rutas implementadas

| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| GET | `/` | Muestra la tabla de series |
| GET | `/create` | Formulario para agregar una serie |
| POST | `/create` | Recibe el formulario e inserta en la DB |
| POST | `/update?id=X` | Suma +1 al episodio actual |
| POST | `/decrement?id=X` | Resta -1 al episodio actual |
| DELETE | `/delete?id=X` | Elimina una serie |
| GET | `/static/*` | Sirve archivos estaticos |

## Challenges implementados

- **Estilos y CSS** — Pagina con estilos propios en `static/style.css`
- **Go ordenado en archivos** — Codigo separado en `main.go`, `handlers.go`, `db.go`, `templates.go` y `static.go`
- **JavaScript en archivo separado** — Todo el JS del cliente esta en `static/app.js`
- **Barra de progreso** — Muestra visualmente cuantos episodios se han visto vs el total
- **Serie completada** — Las series terminadas muestran el texto "(Completada)" junto al nombre
- **Boton -1** — Permite restar un episodio con `POST /decrement`
- **Favicon** — Icono SVG servido desde `/static/favicon.svg`
- **Eliminar serie** — Usa el metodo HTTP `DELETE` desde JavaScript con `fetch()`
- **Validacion en servidor** — Se valida que el nombre no este vacio, que los episodios sean numeros validos y que el episodio actual no sea mayor al total
- **Responsive** — La tabla se adapta a pantallas pequenas

## Patrones usados

**POST/Redirect/GET** — Despues de insertar una serie con POST, el servidor responde con `303 See Other` redirigiendo a `/`. Esto evita que al recargar la pagina se vuelva a enviar el formulario.

**fetch() para acciones sin recargar** — Los botones +1, -1 y Eliminar usan `fetch()` desde JavaScript para enviar la peticion al servidor. Al terminar llaman `location.reload()` para actualizar la tabla.
