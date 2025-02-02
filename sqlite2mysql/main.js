const mysql = require('mysql2');
const Database = require('better-sqlite3');
const db = new Database('./ruten.sqlite');

(async () => {
    // MySQL connection
    const mysqlDb = mysql.createConnection({
        host: '192.168.1.30',
        user: 'sql_192_168_1_30',
        password: '3839b7857eb1f',
        database: 'sql_192_168_1_30',
        charset: 'utf8mb4',
    }).promise();


    try {
        const stmt = db.prepare('SELECT * FROM itemview');

        // 使用 iterate 方法逐行讀取
        for (const row of stmt.iterate()) {
            console.log('Row:', row);
            let images = row.圖片.replace(/\\/g, '/').replace('C:/lutian/PicBackup', '');
            images = images.split('|');
            const jsonString = JSON.stringify(images);
            let temp = [];
            const skulist = [];
            try {
                const skus = JSON.parse(row.自訂規格)[1];
                for (let key in skus) {
                    const lang = {
                        zh_hk: key,
                        zh_cn: key,
                    };
                    skulist.push({ name: lang, image: '' });
                }
                temp = [{
                    "name": {
                        "zh_hk": "規格",
                        "zh_cn": "規格"
                    },
                    "values": [],
                    "isImage": false
                }];

                temp[0].values = skulist;
                temp = JSON.stringify(temp);
            } catch (e) {

            }

            const [rows, fields] = await mysqlDb.query(
                `INSERT INTO products_1 
                (price, images, variables,brand_id,video) 
                VALUES 
                (?, ?, ?,0,"")`,
                [row.直購價, jsonString, temp]);
            console.log('Inserted:', rows.insertId);
        }

    } catch (e) {
        throw e
    }
    // Close the connections gracefully
    process.on('exit', () => {
        sqliteDb.close();
        mysqlDb.end();
    });

})();