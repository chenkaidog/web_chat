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

document.getElementById("delete_chat_but").addEventListener('click', function () {
    var confirmResponse = confirm("删除后数据不可恢复");
    if (confirmResponse == true) {
        localStorage.removeItem("chat_record");
        document.getElementById("content_frame").src = '/index/chat';
    }
})

function markMenuButton() {
    var elems = document.getElementsByClassName('menu_btn');

    Array.from(elems).forEach(function (elem) {
        elem.classList.remove('active')
    });
}