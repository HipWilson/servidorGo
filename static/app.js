// Send +1 or -1 to the server without reloading the page
async function changeEp(id, action) {
    const url = "/" + action + "?id=" + id;

    const response = await fetch(url, { method: "POST" });

    if (response.ok) {
        location.reload();
    }
}

// Delete a series using the DELETE method
async function deleteSeries(id) {
    if (!confirm("Seguro que quieres eliminar esta serie?")) {
        return;
    }

    const response = await fetch("/delete?id=" + id, { method: "DELETE" });

    if (response.ok) {
        location.reload();
    }
}