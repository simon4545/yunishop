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
            let images = row.圖片.replace(/\\/g, '/').replace(/C:\/lutian\/PicBackup/g, '');
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

            } catch (e) {

            }
            tempstr = JSON.stringify(temp);
            let [rows, fields] = await mysqlDb.query(
                `INSERT INTO products_1 
                (price, images, variables,brand_id,video) 
                VALUES 
                (?, ?, ?,0,"")`,
                [row.直購價, jsonString, tempstr]);
            console.log('Inserted:', rows.insertId);
            productid = rows.insertId;
            await mysqlDb.query(
                `INSERT INTO product_categories_1
                (product_id, category_id) 
                VALUES 
                (?, ?)`,
                [productid, 100003]);
            await mysqlDb.query(
                `INSERT INTO product_descriptions_1
                (product_id, locale, name,content) 
                VALUES 
                (?, ?, ?,?)`,
                [productid, 'zh_hk', row.標題, row.說明]);


            let skupos = null
            try {
                skupos = JSON.parse(row.自訂規格)[1];
            } catch (e) {

            }
            if (skupos != null) {
                let index = 0;
                for (let key in skupos) {
                    const value = skupos[key];
                    let is_default = index == 0 ? 1 : 0;
                    await mysqlDb.query(
                        `INSERT INTO product_skus_1
                        (product_id, variants, position,sku,price,origin_price,cost_price,quantity,is_default) 
                        VALUES 
                        (?, ?, ?,?,?,?,?,?,?)`,
                        [productid, `["${index}"]`, index, rand(), value.price, value.price, value.price, 1000, is_default]);
                    index++;
                }
            } else {
                await mysqlDb.query(
                    `INSERT INTO product_skus_1
                    (product_id, variants, position,sku,price,origin_price,cost_price,quantity,is_default) 
                    VALUES 
                    (?, ?, ?,?,?,?,?,?,?)`,
                    [productid, `""`, 0, rand(), row.直購價, row.直購價, row.直購價, 1000, 1]);

            }

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


function rand() {
    const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
    const numbers = "0123456789";

    let result = "";

    // 生成5个随机字母
    for (let i = 0; i < 5; i++) {
        result += letters.charAt(Math.floor(Math.random() * letters.length));
    }

    // 生成10个随机数字
    for (let i = 0; i < 10; i++) {
        result += numbers.charAt(Math.floor(Math.random() * numbers.length));
    }

    return result;
}