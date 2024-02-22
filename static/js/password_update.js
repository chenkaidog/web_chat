document.getElementById('update_but').addEventListener('click', function () {
    var password = document.getElementById('password').value;
    var newPassword = document.getElementById('new_password').value;
    if (newPassword.trim() === '' || password.trim() === '') {
        alert("密码不能为空！");
        return;
    }
    if (newPassword.length < 8) {
        alert("密码长度不能小于8！");
        return;
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
        .then(response => response.json())
        .then(data => {
            var code = data.code;
            var success = data.success;
            var msg = data.message;

            if (success) {
                alert("修改成功，即将重新登陆");
                window.location.reload();
            } else {
                alert(msg);
            }
        })
        .catch(
            error => {
                console.log(error)
            }
        );
});

document.getElementById('logout_but').addEventListener('click', function () {
    fetch(
        '/logout',
        {
            method: 'POST'
        }
    )
        .then(response => response.json())
        .then(data => {
            var code = data.code;
            var success = data.success;
            var msg = data.message;

            if (success) {
                window.location.reload();
            } else {
                alert(msg);
            }
        })
        .catch(
            error => {
                console.log(error)
            }
        );
});