let blacklistCounter = 1;

document.querySelector('.add-button').addEventListener('click', function(e) {
    e.preventDefault();

    // Find the form
    var form = document.getElementById('blacklist-form');

    // Create the new div
    var newDiv = document.createElement('div');
    newDiv.id = 'blacklist-' + blacklistCounter;
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

    // 将黑名单数据存储到本地存储
    localStorage.setItem('blacklist', JSON.stringify(blacklist));

    // 打开新的页面
    window.open('history.html', '_blank');
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
        } else if (xhr.readyState === 4) {
            console.log('请求失败');
        }
    };
}