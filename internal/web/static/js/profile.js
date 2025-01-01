document.getElementById("updateAvatarBtn").addEventListener("click", function () {
    const userID = "{{ .userID }}";
    fetch(`/api/user/${userID}/update-avatar`, {
        method: "POST",
    })
        .then(response => response.json())
        .then(data => {
            if (data.message) {
                alert(data.message);
                window.location.reload();
            } else {
                alert('Ошибка обновления фото');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Произошла ошибка при обновлении фото');
        });
});
