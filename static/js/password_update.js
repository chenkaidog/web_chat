document.getElementById('update_but').addEventListener('click', function () {
    var password = document.getElementById('password').value;
    var newPassword = document.getElementById('new_password').value;
    var confirmPassword = document.getElementById('confirm_password').value;

    if (newPassword.trim() === '' || password.trim() === '') {
        alert("密码不能为空！");
        return;
    }
    if (newPassword.length < 8) {
        alert("密码长度不能小于8！");
        return;
    }

    if (newPassword !== confirmPassword) {
        alert("密码不一致！");
        return
    }

    fetch(
        '/password/update',
        {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                new_password: newPassword,
                password: password
            })
        }
    )

        .then(response => {
            if (response.ok) {
                response.json().then(data => {
                    var success = data.success;
                    var msg = data.message;

                    if (success) {
                        alert("修改成功，即将重新登陆");
                        window.location.reload();
                    } else {
                        alert(msg);
                    }
                })
            } else {
                response.text().then(errorText => {
                    alert(`请求失败, ${response.status}: ${errorText}`);
                    return finishAssistantResponse();
                });
            }
        })
});

document.getElementById('logout_but').addEventListener('click', function () {
    fetch(
        '/logout',
        {
            method: 'POST'
        }
    )
        .then(response => {
            if (response.ok) {
                response.json().then(data => {
                    var success = data.success;
                    var msg = data.message;

                    if (success) {
                        window.location.reload();
                    } else {
                        alert(msg);
                    }
                })
            } else {
                response.text().then(errorText => {
                    alert(`请求失败, ${response.status}: ${errorText}`);
                    return finishAssistantResponse();
                });
            }
        })
});
