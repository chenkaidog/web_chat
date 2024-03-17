document.getElementById("account_btn").addEventListener('click', function () {
    window.location.href = '/index/password_update';
})

document.getElementById("readme_but").addEventListener('click', function () {
    markMenuButton();
    this.classList.add('active');
    document.getElementById("content_frame").src = '/index/readme';
})

document.getElementById("chat_record_but").addEventListener('click', function () {
    markMenuButton();
    this.classList.add('active');
    document.getElementById("content_frame").src = '/index/chat';
})

document.getElementById("create_chat_but").addEventListener('click', function () {
    markMenuButton();
    this.classList.add('active');
    alert("功能开发中");
})

function markMenuButton() {
    var elems = document.getElementsByClassName('menu_btn');

    Array.from(elems).forEach(function (elem) {
        elem.classList.remove('active')
    });
}