let blacklistCounter = 1;

document.getElementById('add-button-0').addEventListener('click', function(e) {
    e.preventDefault();

    // Find the form
    var form = document.getElementById('blacklist-form');

    // Create the new div
    var newDiv = document.createElement('div');
    newDiv.id = 'blacklist-group-' + blacklistCounter;
    newDiv.className = 'input-group';

    // Create the new label
    var newLabel = document.createElement('label');
    newLabel.htmlFor = 'blacklist-' + blacklistCounter;
    newLabel.textContent = '黑名单:';

    // Create the new input box
    var newInput = document.createElement('input');
    newInput.type = 'text';
    newInput.name = 'blacklist-' + blacklistCounter;
    newInput.id = 'blacklist-' + blacklistCounter;
    newInput.className = 'new-blacklist-input';

    // Create the new button
    var newButton = document.createElement('button');
    newButton.textContent = '删除';
    newButton.id = 'add-button-' + blacklistCounter;
    newButton.addEventListener('click', function(e) {
        e.preventDefault();
        newDiv.remove();
    });

    // Add the elements to the new div
    newDiv.appendChild(newLabel);
    newDiv.appendChild(newInput);
    newDiv.appendChild(newButton);

    // Add the new div to the form
    form.insertBefore(newDiv, document.getElementById('monitor-button'));

    // Increment the counter
    blacklistCounter++;
});

document.getElementById('back-button').addEventListener('click', function(e) {
    // 显示表单并隐藏浏览器历史和错误消息

    document.getElementById('back-button').style.display = 'none';  // Hide back button
    document.getElementById('blacklist-div').style.display = 'block';
    document.getElementById('homepage-title').style.display = 'block';
    document.getElementById('history-div').style.display = 'none';
    document.getElementById('error-message').style.display = 'none'; // hide the error message

});

document.getElementById('monitor-button').addEventListener('click', function (e) {
    e.preventDefault();

    // 收集所有黑名单输入框中的数据
    const blacklist = [];
    for (let i = 0; i < blacklistCounter; i++) {
        const input = document.getElementById('blacklist-' + i);
        if (input && input.value) {
            blacklist.push(input.value);
        }
    }

    // 发送请求
    sendMonitorRequest(blacklist);

    // 隐藏表单并显示浏览器历史
    document.getElementById('back-button').style.display = 'block';
    document.getElementById('blacklist-div').style.display = 'none';
    document.getElementById('homepage-title').style.display = 'none';
    document.getElementById('history-div').style.display = 'block';
});



function sendMonitorRequest(blacklist) {
    // 创建请求对象
    const xhr = new XMLHttpRequest();
    xhr.open('POST', 'http://localhost:8080/monitor', true);
    xhr.setRequestHeader('Content-Type', 'application/json;charset=UTF-8');

    // 发送请求
    xhr.send(JSON.stringify({ blacklist: blacklist }));

    // 监听响应
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            console.log('请求成功');
            // 解析服务器返回的 JSON 数据
            const response = JSON.parse(xhr.responseText);
            const historyData = response.result; // 注意这里从返回的数据中提取 result 字段

            // 将数据插入表格
            const historyTable = document.getElementById('history-result');
            historyTable.innerHTML = ''; // 清空表格

            // 这里显示表头和表格内容
            document.querySelector('#history-table thead').style.display = 'table-header-group';
            document.querySelector('#history-table tbody').style.display = 'table-row-group';

            for (const item of historyData) {
                const newRow = historyTable.insertRow();

                newRow.insertCell().textContent = item.title; // 添加 title 列
                newRow.insertCell().textContent = item.url; // 添加 url 列
                newRow.insertCell().textContent = item.visit_count; // 添加 visit_count 列
                newRow.insertCell().textContent = convertTime(item.last_visit_time); // 添加 last_visit_time 列

                // 如果 URL 在黑名单中，改变行的背景色
                // 添加 isInBlacklist 列
                const blacklistCell = newRow.insertCell();
                blacklistCell.textContent = item.is_in_blacklist ? '是' : '否';

                // 如果 URL 在黑名单中，改变行的背景色
                if (item.is_in_blacklist) {
                    newRow.style.backgroundColor = '#F8CECC';
                }
            }

            // 隐藏黑名单表单并显示浏览历史
            document.getElementById('blacklist-div').style.display = 'none';
            document.getElementById('homepage-title').style.display = 'none';
            document.getElementById('history-div').style.display = 'block';
            document.getElementById('error-message').style.display = 'none';
        } else if (xhr.readyState === 4) {
            console.log('请求失败');

            // 解析服务器返回的错误JSON
            const response = JSON.parse(xhr.responseText);
            const errorMessage = response.error; // 注意这里从返回的数据中提取 error 字段

            // 将错误信息设置为错误提示<div>的文本，并将<div>显示出来
            const errorMessageDiv = document.getElementById('error-message');
            errorMessageDiv.textContent = errorMessage;
            errorMessageDiv.style.display = 'block';

            // 这里隐藏表头和表格内容
            document.querySelector('#history-table thead').style.display = 'none';
            document.querySelector('#history-table tbody').style.display = 'none';


            document.getElementById('error-message').style.display = 'block';
            document.getElementById('history-div').style.display = 'none';
        }
    };
}

// 实现 convertTime 函数，将微秒时间戳转换为易读的日期时间格式
function convertTime(microseconds) {
    // Edge 浏览器的时间戳是以微秒为单位的，从 1601 年 1 月 1 日开始
    // 我们需要将这个时间戳转换为以毫秒为单位的，从 1970 年 1 月 1 日开始
    const timestamp = (microseconds - 11644473600000000) / 1000;
    const date = new Date(timestamp);
    return date.toLocaleString(); // 返回本地格式的日期时间字符串
}
