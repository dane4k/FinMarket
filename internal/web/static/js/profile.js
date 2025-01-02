const usrID = userID;
console.log(usrID)
document.getElementById("updateAvatarBtn").addEventListener("click", function () {
    fetch(`/api/user/${usrID}/update-avatar`, {
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
