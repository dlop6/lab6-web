# La Liga Tracker

## Documentación de API
[![Swagger](https://img.shields.io/badge/Swagger-Documentación-blue)](http://localhost:8080/swagger/index.html)


## Endpoints PATCH
| Método | Ruta                     | Body Request Example      | Descripción               |
|--------|--------------------------|---------------------------|---------------------------|
| PATCH  | /matches/:id/goals       | `{"team":"team1","goals":2}` | Actualiza goles           |
| PATCH  | /matches/:id/yellowcards | -                         | +1 tarjeta amarilla       |
| PATCH  | /matches/:id/redcards    | -                         | +1 tarjeta roja           |
| PATCH  | /matches/:id/extratime   | `{"minutes":5}`           | Establece tiempo extra    |

## Screenshots
| Vista          | Imagen                      |
|----------------|-----------------------------|
| Listar partidos | ![Listado](/ss/cargar.png)  |
| Crear partido  | ![Creación](/ss/crear.png)  |
| Actualizar partido  | ![Creación](/ss/update.png)  |
| Eliminar partido  | ![Creación](/ss/delete.png)  |
