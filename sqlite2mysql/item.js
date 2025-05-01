const mysql = require('mysql2');
const fs = require('fs').promises;
const path = require('path');
(async () => {
    const mysqlDb = mysql.createConnection({
        host: '192.168.1.30',
        user: 'sql_192_168_1_30',
        password: '3839b7857eb1f',
        database: 'sql_192_168_1_30',
        charset: 'utf8mb4',
    }).promise();

    let [results, fields] = await mysqlDb.query('SELECT id, content FROM product_descriptions')
    for (const row of results) {
        console.log('Row:', row.id);
        try {
            row.content = row.content.replace(/\\/g, '/')
            if (row.content.startsWith('/item/')) {
                const filePath = path.join('/Users/forke/fsdownload/', row.content);
                let content = await fs.readFile(filePath, 'utf8')
                await mysqlDb.query('UPDATE product_descriptions SET content = ? WHERE id = ?', [content, row.id]);
            }
        } catch (e) { 
            continue
        }
    };
})();
