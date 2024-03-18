const platformBaidu = "baidu";
const platformOpenai = "openai";
const modelErine4 = "erine-4";
const modelGPT3 = "gpt-3.5";
const modelGPT4 = "gpt-4";

const roleUser = "user";
const roleAssistant = "assistant";

const pageStageUserInput = "user_inout";
const pageStageAssistantOutput = "ai_output";

const chatRecordItem = 'chat_record';

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

var assistantRespContent = '';
var userQuestionContent = '';

var chatRecord = [];
const recordMaxSize = 10;

// 初始化对话记录
function initChatRecord() {
    document.getElementById("msg_input").readOnly = true;
    chatRecord = JSON.parse(localStorage.getItem(chatRecordItem));
    if (chatRecord == null) {
        chatRecord = []
    }
    var chatRecordHtml = '';
    for (var i = 0; i < chatRecord.length; i++) {
        if (chatRecord[i].length != 2) {
            continue;
        }

        chatRecordHtml += `<div class="content"><label><b>You</b></label><br><p>${renderInput(chatRecord[i][0])}</p><br></div>`
        chatRecordHtml += `<div class="content"><label><b>AI</b></label><br>${md.render(chatRecord[i][1])}<br></div>`
    }

    var outputDiv = document.getElementById("output");
    outputDiv.innerHTML = chatRecordHtml;
    outputDiv.scrollTop = outputDiv.scrollHeight;

    document.getElementById("msg_input").readOnly = false;
}

initChatRecord();

function recordMessages(userMsg, assistantMsg) {
    chatRecord.push([userMsg, assistantMsg]);
    localStorage.setItem(chatRecordItem, JSON.stringify(chatRecord));
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
    if (event.ctrlKey && event.key === "Enter") {
        if (pageStage === pageStageUserInput) {
            sendUserMessage();
            return;
        }
    }
});

// 监听发送按钮点击事件
document.getElementById("send_but").addEventListener('click', function () {
    if (pageStage === pageStageUserInput) {
        return sendUserMessage();
    }
    if (pageStage === pageStageAssistantOutput) {
        return abortAiOutput();
    }
});

// 删除对话事件
document.getElementById("delete_but").addEventListener('click', function () {
    localStorage.removeItem("chat_record");
    var outputDiv = document.getElementById("output");
    var sp = document.createElement('div');
    sp.setAttribute("class", "separator");
    var hr = document.createElement('hr');
    var span = document.createElement('span');
    span.innerHTML = "历史对话已清空";
    sp.appendChild(hr);
    sp.appendChild(span);
    outputDiv.appendChild(sp);
    outputDiv.scrollTop = outputDiv.scrollHeight;
})

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
        console.log('event source close  err:', error);
    } finally {
        recordMessages(userQuestionContent, assistantRespContent);
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
        return alert("消息不能为空!!!");
    }

    return parseChatReq(userMessage);
}

function getChatType() {
    return document.getElementById("chat_type_label").getAttribute("value")
}

function parseChatReq(msg) {
    var chatID = new Date().getTime().toString();

    // 将用户问题添加到页面上
    startGenerateResp(msg, chatID);

    var messages = [];
    for (var i = chatRecord.length - 1; i >= 0; i--) {
        var record = chatRecord[i];
        if (record.length != 2) {
            continue;
        }

        // 从后往前添加对话记录
        messages.unshift(
            {
                role: roleUser,
                content: record[0]
            },
            {
                role: roleAssistant,
                content: record[1]
            },
        )

        if (messages.length >= recordMaxSize) {
            break;
        }
    }

    var userMsg = {
        role: roleUser,
        content: msg
    }

    return sendMessage(messages, userMsg, chatID);
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

function sendMessage(messagesContext, userMsg, chatID) {
    messagesContext.push(userMsg);
    var selectedModel = getSelectedModel()
    var chatReq = {
        platform: selectedModel.platform,
        model: selectedModel.model,
        messages: messagesContext
    }

    userQuestionContent = userMsg.content;
    assistantRespContent = '';

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
                            recordMessages(userQuestionContent, assistantRespContent);
                            return finishAssistantResponse();
                        }

                        assistantRespContent += message.content
                        renderResp(assistantRespContent, chatID);
                    }

                    eventsource.onerror = function (error) {
                        alert('请求失败: ', error);
                        this.close();
                        return finishAssistantResponse();
                    }

                    eventsource.onclose = function (event) {
                        return finishAssistantResponse();
                    }
                })
            } else {
                response.text().then(errorText => {
                    alert(`请求失败, ${response.status}: ${errorText}`);
                    return finishAssistantResponse();
                });
            }
        });
}

function startGenerateResp(textarea, respDivID) {
    // clean up textarea
    document.getElementById("msg_input").value = "";
    autoResize();

    var outputDiv = document.getElementById("output");

    var userContentDiv = document.createElement('div');
    userContentDiv.setAttribute('class', 'content')
    userContentDiv.innerHTML = `<label><b>You</b></label><br><p>${renderInput(textarea)}</p><br>`
    outputDiv.appendChild(userContentDiv);

    var respDiv = document.createElement('div');
    respDiv.setAttribute('id', respDivID)
    respDiv.setAttribute('class', "content")
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