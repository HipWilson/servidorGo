package main

import "fmt"

func buildIndexPage(rows string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Series Tracker</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <div class="container">
        <h1>Series Tracker</h1>
        <a href="/create" class="btn-add">+ Agregar serie</a>

        <table id="seriesTable">
            <thead>
                <tr>
                    <th>#</th>
                    <th>Nombre</th>
                    <th>Episodios</th>
                    <th>Progreso</th>
                    <th>Acciones</th>
                </tr>
            </thead>
            <tbody>
                %s
            </tbody>
        </table>
    </div>
    <script src="/static/app.js"></script>
</body>
</html>`, rows)
}

func buildCreatePage(errorMsg string, name string, current string, total string) string {
	errorHTML := ""
	if errorMsg != "" {
		errorHTML = fmt.Sprintf(`<p class="error">%s</p>`, errorMsg)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Agregar Serie</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <div class="container">
        <h1>Agregar nueva serie</h1>
        <a href="/" class="btn-back">Volver</a>

        %s

        <form method="POST" action="/create">
            <div class="form-group">
                <label for="series_name">Nombre de la serie</label>
                <input type="text" id="series_name" name="series_name"
                       placeholder="Ej: Breaking Bad" required value="%s">
            </div>
            <div class="form-group">
                <label for="current_episode">Episodio actual</label>
                <input type="number" id="current_episode" name="current_episode"
                       min="1" value="%s" required>
            </div>
            <div class="form-group">
                <label for="total_episodes">Total de episodios</label>
                <input type="number" id="total_episodes" name="total_episodes"
                       min="1" placeholder="Ej: 62" value="%s" required>
            </div>
            <button type="submit" class="btn-submit">Guardar</button>
        </form>
    </div>
</body>
</html>`, errorHTML, name, current, total)
}
