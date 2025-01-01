const userAuthToken = authToken;
if (!userAuthToken) {
    console.error("Ошибка: токен авторизации пустой.");
    alert("Ошибка: токен не передан или некорректен.");
} else {
    let tokenExpired = false;

    function checkStatus() {
        if (tokenExpired) return;

        fetch(`/check-status/${userAuthToken}`)
            .then(response => {
                if (response.status === 404) {
                    console.error("Token not found.");
                    tokenExpired = true;
                    alert("Токен не найден. Пожалуйста, попробуйте заново.");
                    return;
                }
                return response.json();
            })
            .then(data => {
                if (data) {
                    console.log('Received data:', data);
                    if (data.status === "confirmed") {
                        window.location.href = "/profile";
                    } else if (data.error === "token expired") {
                        alert("Токен истёк. Пожалуйста, авторизуйтесь заново.");
                        window.location.href = "/";
                    }
                }
            })
            .catch(err => console.error("Ошибка проверки статуса:", err));
    }

    setInterval(checkStatus, 5000);
}

