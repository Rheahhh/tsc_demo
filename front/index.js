// 切换显示页面
document.getElementById('blacklist-button').addEventListener('click', function() {
    document.getElementById('blacklist').classList.remove('hidden');
    document.getElementById('alerts').classList.add('hidden');
    document.getElementById('home').classList.add('hidden');
    fetchBlacklist();
});
document.getElementById('alerts-button').addEventListener('click', function() {
    document.getElementById('alerts').classList.remove('hidden');
    document.getElementById('blacklist').classList.add('hidden');
    document.getElementById('home').classList.add('hidden');
    fetchAlerts();
});
document.getElementById('home-button').addEventListener('click', function() {
    document.getElementById('home').classList.remove('hidden');
    document.getElementById('blacklist').classList.add('hidden');
    document.getElementById('alerts').classList.add('hidden');
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
                row.classList.add('input-group', 'mb-3');

                var urlInput = document.createElement('input');
                urlInput.type = 'text';
                urlInput.value = item;
                urlInput.classList.add('form-control');
                urlInput.readOnly = true;
                row.appendChild(urlInput);

                var deleteButtonContainer = document.createElement('div');
                deleteButtonContainer.classList.add('input-group-append');

                var deleteButton = document.createElement('button');
                deleteButton.textContent = '删除';
                deleteButton.id = 'delete-button';
                deleteButton.classList.add('btn', 'btn-outline-secondary');
                deleteButton.addEventListener('click', function() {
                    manageBlacklist('delete', item);
                });
                deleteButtonContainer.appendChild(deleteButton);

                row.appendChild(deleteButtonContainer);
                table.appendChild(row);
            });
        });
}




// 获取告警历史
function fetchAlerts() {
    fetch('http://localhost:8080/alerts')
        .then(response => response.json())
        .then(data => {
            var table = document.getElementById('alerts-table');
            table.innerHTML = '';
            data.forEach(item => {
                var row = document.createElement('tr');

                var clientIdCell = document.createElement('td');
                clientIdCell.textContent = item.client_id;
                row.appendChild(clientIdCell);

                var nameCell = document.createElement('td');
                nameCell.textContent = item.name;
                row.appendChild(nameCell);

                var urlCell = document.createElement('td');
                urlCell.textContent = item.url;
                row.appendChild(urlCell);

                var viewCountCell = document.createElement('td');
                viewCountCell.textContent = item.view_count;
                row.appendChild(viewCountCell);

                var timeCell = document.createElement('td');
                timeCell.textContent = new Date(item.view_timestamp).toLocaleString();
                row.appendChild(timeCell);

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
