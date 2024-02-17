$(document).ready(function () {
    $('#loginBut').click(function () {
        var username = $('#username').val();
        var password = $('#password').val();
        // 校验用户名和密码
        if (username.trim() === '' || password.trim() === '') {
            alert("用户名和密码不能为空！");
            return;
        }
        if (password.length < 8) {
            alert("密码长度不能小于8！");
            return;
        }

        $.post('/login',
            {
                username: username,
                password: password
            },
            function (data, status, xhr) {
                if (status == 'success') {
                    var code = data.code;
                    var success = data.success;
                    var msg = data.message;

                    if (success) {
                        window.location.href = "/index/home"
                    } else {
                        alert(msg);
                    }
                } else {
                    alert("network err, please retry!");
                }
            }
        )
    });
});