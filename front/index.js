// 切换显示页面
document.getElementById('blacklist-button').addEventListener('click', function() {
    document.getElementById('blacklist').classList.remove('hidden');
    document.getElementById('alerts').classList.add('hidden');
    fetchBlacklist();
});
document.getElementById('alerts-button').addEventListener('click', function() {
    document.getElementById('alerts').classList.remove('hidden');
    document.getElementById('blacklist').classList.add('hidden');
    fetchAlerts();
});

// 添加新的黑名单项
document.getElementById('add-button').addEventListener('click', function() {
    var url = document.getElementById('blacklist-input').value;
    if (url) {
        manageBlacklist('add', url);
    }
});

// 获取黑名单
function fetchBlacklist() {
    fetch('http://localhost:8080/blacklist')
        .then(response => response.json())
        .then(data => {
            var table = document.getElementById('blacklist-table');
            table.innerHTML = '';
            data.forEach(item => {
                var row = document.createElement('div');
                row.className = 'table-row';
                var urlDiv = document.createElement('div');
                urlDiv.textContent = item;
                var buttonDiv = document.createElement('div');
                var deleteButton = document.createElement('button');
                deleteButton.textContent = '删除';
                deleteButton.addEventListener('click', function() {
                    manageBlacklist('delete', item);
                });
                buttonDiv.appendChild(deleteButton);
                row.appendChild(urlDiv);
                row.appendChild(buttonDiv);
                table.appendChild(row);
            });
        });
}

// 获取告警历史
// 获取告警历史
function fetchAlerts() {
    fetch('http://localhost:8080/alerts')
        .then(response => response.json())
        .then(data => {
            var table = document.getElementById('alerts-table');
            table.innerHTML = '';
            data.forEach(item => {
                var row = document.createElement('div');
                row.className = 'table-row';

                var clientIdDiv = document.createElement('div');
                clientIdDiv.textContent = item.client_id;
                row.appendChild(clientIdDiv);

                var nameDiv = document.createElement('div');
                nameDiv.textContent = item.name;
                row.appendChild(nameDiv);

                var urlDiv = document.createElement('div');
                urlDiv.textContent = item.url;
                row.appendChild(urlDiv);

                var viewCountDiv = document.createElement('div');
                viewCountDiv.textContent = item.view_count;
                row.appendChild(viewCountDiv);

                var timeDiv = document.createElement('div');
                timeDiv.textContent = new Date(item.visit_time).toLocaleString();
                row.appendChild(timeDiv);

                table.appendChild(row);
            });
        });
}


// 管理黑名单
function manageBlacklist(action, url) {
    fetch('http://localhost:8080/blacklist', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ action, url })
    })
        .then(() => {
            fetchBlacklist();
        });
}

// 初始加载黑名单页面
fetchBlacklist();
