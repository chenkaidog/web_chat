document.getElementById("account_btn").addEventListener('click', function () {
    window.location.href = '/index/password_update';
})

document.getElementById("readme_but").addEventListener('click', function () {
    markMenuButton();
    this.style.backgroundColor = '#999999';
    document.getElementById("content_frame").src = '/index/readme';
})
 
document.getElementById("single_chat_but").addEventListener('click', function () {
    markMenuButton();
    this.style.backgroundColor = '#999999';
    document.getElementById("content_frame").src = '/index/chat?chat_type=single';
})

document.getElementById("new_chat_but").addEventListener('click', function () {
    markMenuButton();
    this.style.backgroundColor = '#999999';
    document.getElementById("content_frame").src = '/index/chat?chat_type=round';
})

function markMenuButton() {
    var elems = document.getElementsByClassName('menu_btn');

    Array.from(elems).forEach(function (elem) {
        elem.style.backgroundColor = '#dddddddd';
    });
}