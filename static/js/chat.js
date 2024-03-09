
const chatTypeSingle = "single"
const chatTypeRound = "round"

const platformBaidu = "baidu"
const platformOpenai = "openai"
const modelErine4 = "erine-4"
const modelGPT3 = "gpt-3.5"
const modelGPT4 = "gpt-4"

const roleUser = "user"
const roleAssistant = "assistant"

const pageStageUserInput = "user_inout"
const pageStageAssistantOutput = "ai_output"

var pageStage = pageStageUserInput;
var fixOutputBoardToBottom = true;

var eventsource;

var md = window.markdownit({
    linkify: false,
    highlight: function (str, lang) {
        if (lang && hljs.getLanguage(lang)) {
            try {
                return '<pre class="hljs_code"><code>' +
                    hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
                    '</code></pre>';
            } catch (__) { }
        }

        return '<pre class="hljs_code"><code>' + md.utils.escapeHtml(str) + '</code></pre>';
    }
});

var chatRecord = [];
const recordMaxSize = 10;

function recordUserMsg(msg) {
    chatRecord.push([msg]);
    while (chatRecord.length > recordMaxSize) {
        chatRecord.shift();
    }
}

function recordAssistantMsg(msg) {
    if (msg.length <= 0) {
        return;
    }
    var latestRecord = chatRecord[chatRecord.length - 1];
    if (latestRecord.length == 1) {
        latestRecord.push(msg);
    } else if (latestRecord.length == 2) {
        latestRecord[1] = msg;
    }

    while (chatRecord.length > recordMaxSize) {
        chatRecord.shift();
    }
}

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
        if (pageStage === pageStageUserInput) {
            sendUserMessage();
            return;
        }
    }
});

// 监听发送按钮点击事件
document.getElementById("send_but").addEventListener('click', function () {
    if (pageStage === pageStageUserInput) {
        sendUserMessage();
        return;
    }
    if (pageStage === pageStageAssistantOutput) {
        abortAiOutput();
        return;
    }
});

document.getElementById("output").addEventListener('wheel', function (event) {
    if (event.deltaY < 0) {
        fixOutputBoardToBottom = false;
        return;
    }
});

function abortAiOutput() {
    try {
        eventsource.close();
    } catch (error) {

    } finally {
        finishAssistantResponse();
    }
}

function sendUserMessage() {
    document.getElementById("send_but").innerHTML = '⏸️';
    pageStage = pageStageAssistantOutput;

    var textarea = document.getElementById("msg_input");
    var userMessage = textarea.value
    if (userMessage.trim() === "") {
        finishAssistantResponse();
        return alert("消息不能为空!!!")
    }

    return parseChatReq(userMessage);
}

function getChatType() {
    return document.getElementById("chat_type_label").getAttribute("value")
}

function parseChatReq(msg) {
    var chatID = new Date().getTime().toString();

    startGenerateResp(msg, chatID);

    recordUserMsg(msg);

    if (getChatType() === chatTypeSingle) {
        return sendMessage(
            [{
                role: roleUser,
                content: msg
            }],
            chatID)
    }

    if (getChatType() === chatTypeRound) {
        var messages = [];
        for (var i = 0; i < chatRecord.length; i++) {
            var record = chatRecord[i];
            if (record.length != 2) {
                continue;
            }
            messages.push(
                {
                    role: roleUser,
                    content: record[0]
                },
                {
                    role: roleAssistant,
                    content: record[1]
                },
            )
        }

        messages.push({
            role: roleUser,
            content: msg
        })
        return sendMessage(messages, chatID);
    }

    return;
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

function sendMessage(messages, chatID) {
    var selectedModel = getSelectedModel()
    var chatReq = {
        platform: selectedModel.platform,
        model: selectedModel.model,
        messages: messages
    }

    var contentResp = '';

    fetch(
        '/chat/create',
        {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(chatReq),
        })
        .then(response => {
            if (response.ok) {
                response.json().then(data => {
                    if (!data.success) {
                        alert(data.message);
                        return finishAssistantResponse();
                    }

                    eventsource = new EventSource('/chat/stream');
                    eventsource.onmessage = function (event) {
                        var message = JSON.parse(event.data);
                        if (!message.success) {
                            alert(message.message);
                            this.close();
                            return finishAssistantResponse();
                        }

                        if (message.is_end) {
                            this.close();
                            return finishAssistantResponse();
                        }

                        contentResp += message.content
                        renderResp(contentResp, chatID);
                    }

                    eventsource.onerror = function (error) {
                        alert('请求失败');
                        this.close();
                        return finishAssistantResponse();
                    }

                    eventsource.onclose = function (event) {
                        return finishAssistantResponse();
                    }
                })
            } else {
                alert('请求失败');
                return finishAssistantResponse();
            }
        });
}

function startGenerateResp(textarea, respDivID) {
    // clean up textarea
    document.getElementById("msg_input").value = "";
    autoResize();

    var outputDiv = document.getElementById("output");

    var userContentDiv = document.createElement('div');
    userContentDiv.setAttribute('class', 'user_content')
    userContentDiv.innerHTML = `<label><b>You</b></label><br><p>${renderInput(textarea)}</p><br>`
    outputDiv.appendChild(userContentDiv)

    var respDiv = document.createElement('div');
    respDiv.setAttribute('id', respDivID)
    respDiv.setAttribute('class', "assistant_content")
    outputDiv.appendChild(respDiv);

    fixOutputBoardToBottom = true;
    outputDiv.scrollTop = outputDiv.scrollHeight;
}

function renderInput(text) {
    // 将空格替换为HTML空格字符
    text = text.replace(/ /g, "&nbsp;");

    // 将换行符替换为HTML换行标签
    text = text.replace(/\n/g, "<br>");

    // 返回转换后的文本
    return text;
}

function renderResp(content, respDivID) {
    recordAssistantMsg(content);

    var htmlContent = md.render(content)
    var respContent = `<label><b>AI</b></label><br>${htmlContent}<br>`
    document.getElementById(respDivID).innerHTML = respContent;

    if (fixOutputBoardToBottom) {
        var output = document.getElementById("output");
        output.scrollTop = output.scrollHeight;
    }
}

function finishAssistantResponse() {
    document.getElementById("send_but").innerHTML = '⬆️';
    pageStage = pageStageUserInput;
};