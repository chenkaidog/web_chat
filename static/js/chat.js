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

var md = window.markdownit({
    linkify: true,
    highlight: function (str, lang) {
        if (lang && hljs.getLanguage(lang)) {
            try {
                return '<pre><code class="hljs">' +
                    hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
                    '</code></pre>';
            } catch (__) { }
        }

        return '<pre><code class="hljs">' + md.utils.escapeHtml(str) + '</code></pre>';
    }
});

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
    sendMessage([{
        role: roleUser,
        content: msg
    }])
}

function getSelectedModel() {
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

function sendMessage(messages) {
    var selectedModel = getSelectedModel()
    var chatReq = {
        platform: selectedModel.platform,
        model: selectedModel.model,
        messages: messages
    }

    // 生成当前时间的时间戳
    var timestamp = new Date().getTime();

    // 将时间戳转换为字符串
    var timestampStr = timestamp.toString();
    var contentResp = ''

    fetch('/chat/stream', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(chatReq)
    })
        .then(function (response) {
            if (response.ok) {
                startGenerateResp(chatReq.messages[chatReq.messages.length - 1].content, timestampStr);

               // 处理接收到的事件流数据
        var reader = response.body.getReader();
        var decoder = new TextDecoder();
        var buffer = '';

        function read() {
            return reader.read().then(function (result) {
                if (result.done) {
                    // 读取完成
                    console.log('Event stream reading completed');
                    return;
                }

                // 将接收到的数据添加到缓冲区
                buffer += decoder.decode(result.value, { stream: true });

                // 按行分割处理数据
                var lines = buffer.split('\n');
                buffer = lines.pop(); // 保留最后一行，可能是不完整的数据

                lines.forEach(function (line) {
                    if (line.startsWith("data:")) {
                        var eventData = line.replace("data:", "").trim(); // 删除前缀并去除空格
                        var msgResp = JSON.parse(eventData);

                        if (!msgResp.success) {
                            alert(msgResp.message);
                            return;
                        }
                        if (msgResp.is_end) {
                            return;
                        }

                        contentResp += msgResp.content;
                        renderResp(contentResp, timestampStr);
                    }
                });

                // 继续读取下一个事件
                return read();
                    });
                }

                // 开始读取事件流数据
                return read();
            } else {
                alert('请求失败');
            }
        })
        .catch(function (error) {
            console.log(error);
        });
}

function startGenerateResp(textarea, respDivID) {
    // clean up textarea
    document.getElementById("msg_input").value = "";

    var outputDiv = document.getElementById("output");

    var userContentDiv = document.createElement('div');
    userContentDiv.setAttribute('class', 'user_content')
    userContentDiv.innerHTML = `<label><b>You</b></label><br><p>${textarea}</p><br>`
    outputDiv.appendChild(userContentDiv)

    var respDiv = document.createElement('div');
    respDiv.setAttribute('id', respDivID)
    respDiv.setAttribute('class', "assistant_content")
    outputDiv.appendChild(respDiv)
}

function renderResp(content, respDivID) {
    var htmlContent = md.render(content)
    var respContent = `<label><b>AI</b></label><br>${htmlContent}<br>`
    document.getElementById(respDivID).innerHTML = respContent;
} 