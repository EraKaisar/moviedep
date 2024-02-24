
function showForm(formId) {
    document.getElementById(formId).style.display = 'block';
}




function createMovie() {
    var form = $("#create-form");
    var data = form.serialize();

    $.post("/createMovie", data, function (response) {
        if (response.status === "success") {
            alert("Статья успешно создана!");
            location.reload();
        } else {
            alert("Ошибка при создании статьи!");
        }
    });
}
function updateMovie(movieId) {
    var form = $("#update-form-" + movieId);
    var data = form.serialize();

    $.post("/updateMovie", data, function (response) {
        if (response.status === "success") {
            alert("Статья успешно обновлена!");
            location.reload(); // Обновить страницу после успешного обновления
        } else {
            alert("Ошибка при обновлении статьи!");
        }
    });
}

