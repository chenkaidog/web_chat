document.addEventListener("DOMContentLoaded", autoResize);

function autoResize() {
    var textarea = document.getElementById("msg_input");
    var inputDiv = document.getElementById("input");

    if (textarea.value.trim() === "") {
        inputDiv.style.height = "45px";
        return;
    }

    inputDiv.style.height = "auto"
    inputDiv.style.height = textarea.scrollHeight + "px"; // 根据内容高度调整textarea高度
}