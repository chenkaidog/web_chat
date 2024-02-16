// auto height of textarea
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

// 监听Ctrl+回车事件
document.getElementById("msg_input").addEventListener('keydown', function (event) {
    if (event.ctrlKey && event.keyCode === 13) {
        sendUserMessage();
    }
});

// 监听发送按钮点击事件
document.getElementById("send_but").addEventListener('click', function () {
    sendUserMessage();
});

function sendUserMessage() {
    var textarea = document.getElementById("msg_input");
    var userMessage = textarea.value
    if (userMessage.trim() === "") {
        alert("消息不能为空!!!")
        return
    }
    if (userMessage.length > 500) {
        alert("消息长度不能大于500个字符!!!")
        return
    }

    if (getChatType() === chatTypeSingle) {
        singleChat(userMessage)
        return
    }

    return
}

function getChatType() {
    return document.getElementById("chat_type_label").getAttribute("value")
}

const chatTypeSingle = "single"
const chatTypeRound = "round"

const platformBaidu = "baidu"
const platformOpenai = "openai"
const modelErine4 = "erine-4"
const modelGPT3 = "gpt-3.5"
const modelGPT4 = "gpt-4"

const roleUser = "user"
const roleAssistant = "assistant"

// single chat
function singleChat(msg) {
    var selectedModel = getselectedModel()
    var chatReq = {
        platform: selectedModel.platform,
        model: selectedModel.model,
        messages: [{
            role: roleUser,
            content: msg
        }]
    } 

    sendMessage(chatReq)
}

function getselectedModel() {
    var selectElement = document.getElementById("model_opt");
    var selectedModel = selectElement.options[selectElement.selectedIndex];
    var modelValue = selectedModel.getAttribute("value");
    switch (modelValue) {
        case "ernie4":
            return {
                platform: platformBaidu,
                model: modelErine4
            }
        case "gpt3":
            return {
                platform: platformOpenai,
                model: modelGPT3
            }
        case "gpt4":
            return {
                platform: platformOpenai,
                model: modelGPT4
            }
    }
}

function sendMessage(chatCreateReq) {
    
}