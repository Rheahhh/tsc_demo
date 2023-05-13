window.onload = function() {
    // 创建请求对象
    const xhr = new XMLHttpRequest();
    xhr.open('POST', 'http://localhost:8080/monitor', true);
    xhr.setRequestHeader('Content-Type', 'application/json;charset=UTF-8');

    // 读取存储在本地的黑名单数据
    const blacklist = JSON.parse(localStorage.getItem('blacklist'));

    // 发送请求
    xhr.send(JSON.stringify({ blacklist: blacklist }));

    // 监听响应
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            const data = JSON.parse(xhr.responseText);

            // 为每条历史记录添加一行
            for (const record of data.result) {
                const row = document.createElement('tr');

                const urlCell = document.createElement('td');
                urlCell.textContent = record.url;
                row.appendChild(urlCell);

                const visitTimeCell = document.createElement('td');
                visitTimeCell.textContent = new Date(record.visitTime).toLocaleString();
                row.appendChild(visitTimeCell);

                const isInBlacklistCell = document.createElement('td');
                isInBlacklistCell.textContent = record.isInBlacklist ? 'Yes' : 'No';
                row.appendChild(isInBlacklistCell);

                document.getElementById('history-table').appendChild(row);
            }
        } else if (xhr.readyState === 4) {
            console.log('请求失败');
        }
    };
}
