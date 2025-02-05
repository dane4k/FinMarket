const userAuthToken = authToken;

if (!userAuthToken) {
    alert("Ошибка: токен не передан или некорректен.");
    window.location.href = "/auth";
} else {
    let tokenExpired = false;

    function checkStatus() {
        if (tokenExpired) return;

        fetch(`/check-status/${userAuthToken}`)
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    console.error("Ошибка:", data.error);
                    if (data.error === "auth token not found") {
                        alert("Токен не найден. Авторизуйтесь заново");
                        window.location.href = "/auth";
                    } else if (data.error === "auth token expired") {
                        alert("Токен авторизации истек. Авторизуйтесь заново");
                        window.location.href = "/auth";
                    }
                    tokenExpired = true;
                } else if (data.status === "confirmed") {
                    window.location.href = "/profile";
                }
            })
            .catch(err => {
                console.error("Ошибка при проверке токена:", err);
                alert("Ошибка сервера. Пожалуйста, попробуйте позже.");
            });
    }

    setInterval(checkStatus, 5000);
}
