$(document).ready(function () {
    $('#loginBut').click(function () {
        var newPassword = $('#new_password').val();
        var password = $('#password').val();
        if (newPassword.trim() === '' || password.trim() === '') {
            alert("密码不能为空！");
            return;
        }
        if (newPassword.length < 8) {
            alert("密码长度不能小于8！");
            return;
        }

        $.post('/password/update',
            {
                new_password: newPassword,
                password: password
            },
            function (data, status, xhr) {
                if (status == 'success') {
                    var code = data.code;
                    var success = data.success;
                    var msg = data.msg;

                    if (success) {
                        window.location.href = "/login"
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