const mysql = require('mysql2');
const Database = require('better-sqlite3');
const db = new Database('./ruten.sqlite');
const fs = require('fs').promises;
let category_id = 100018;
let fenzu = '淘幂铺';
(async () => {
    // MySQL connection
    const mysqlDb = mysql.createConnection({
        host: '192.168.1.30',
        user: 'lutian',
        password: 'xdWETicFHZFf3bcT',
        database: 'lutian',
        charset: 'utf8mb4',
    }).promise();


    try {
        const stmt = db.prepare(`SELECT * FROM itemview where 分組='${fenzu}'`);
        await mysqlDb.beginTransaction();
        let jsonString=""
        // 使用 iterate 方法逐行讀取
        for (const row of stmt.iterate()) {
            console.log('Row:', row);
            let images = row.圖片.replace(/\\/g, '/').replace(/C:\/lutian\/PicBackup/g, 'https://img.1pan.me/photos').replace(/\/PicBackup/g, 'https://img.1pan.me/photos');
            images=images.replace(/D:\/露天自動上架精靈/g,'')
            images = images.split('|');
            jsonString = JSON.stringify(images);
            let tempstr = "";
            let skutemp = [];
            let escapedData="";
            try {
                try {
                    if (row.自訂規格 != "") {
                        row.自訂規格=row.自訂規格.replace(/\'/g, '"');
                        const skuType = JSON.parse(row.自訂規格)[0];
                        skuType.dataRows.forEach((item) => {
                            const skulist = [];
                            let temp = {
                                "name": {
                                    "zh_hk": item.name
                                },
                                "values": [],
                                "isImage": false
                            }
                            for (let index = 0; index < item.specValue.length; index++) {
                                const element = item.specValue[index];
                                let obj = {
                                    name: {
                                        zh_hk: element.name,
                                    }
                                };
                                if (element.simage != null) {
                                    obj.image = element.simage;
                                    if (obj.image.startsWith('http') == false) {
                                        obj.image = `https:${obj.image}`;
                                    }
                                    temp.isImage = true;
                                }
                                skulist.push(obj);
                            }
                            temp.values = skulist;
                            skutemp.push(temp);
                        })

                        tempstr = JSON.stringify(skutemp);
                    }else{
                        tempstr = "[]";
                    }
                } catch (e) {
                    throw e;
                }
                try {
                    let deschtml = row.說明.replace(/\\/g, '/')
                    let html = await fs.readFile(`/Users/forke/Downloads/trade${deschtml}`, 'utf8')
                    // escapedData = mysql.escape(html);
                    escapedData = html;
                } catch (e) {
                    throw e;
                }
                let price = row.直購價 || row.起標價;
                price=parseFloat(price)
                let [rows, fields] = await mysqlDb.query(`INSERT INTO products (price, images, variables,brand_id,video) VALUES (?, ?, ?,0,"")`, [price, jsonString, tempstr]);

                console.log('Inserted:', rows.insertId);
                productid = rows.insertId;

                if (skutemp.length != 0) {
                    const skuPrice = JSON.parse(row.自訂規格)[1];
                    let result = combineNames(skutemp.map(item => item.values));
                    let result1 = combineIndexes(skutemp.map(item => item.values));
                    for (let index = 0; index < result.length; index++) {
                        const element = result[index];
                        let priceObj = skuPrice[element];
                        if (priceObj) {
                            let price = priceObj.price
                            let is_default = index == 0 ? 1 : 0;
                            await mysqlDb.query(`INSERT INTO product_skus (product_id, variants, position,sku,price,origin_price,cost_price,quantity,is_default) VALUES (?, ?, ?,?,?,?,?,?,?)`, [productid, JSON.stringify(result1[index]), index, rand(), price, price, price, 1000, is_default]);
                        }
                    }
                } else {
                    await mysqlDb.query(`INSERT INTO product_skus (product_id, variants, position,sku,price,origin_price,cost_price,quantity,is_default) VALUES (?, ?, ?,?,?,?,?,?,?)`, [productid, `""`, 0, rand(), price, price, price, 1000, 1]);
                }

                await mysqlDb.query(`INSERT INTO product_categories (product_id, category_id) VALUES (?, ?)`, [productid, category_id]);

                await mysqlDb.query(`INSERT INTO product_descriptions (product_id, locale, name,content) VALUES (?, ?, ?,?)`, [productid, 'zh_hk', row.標題, escapedData]);
            } catch (e) {
                fs.appendFile('error.txt', `${row.ID}\n`);
                console.log('Row:', e);
                throw e;
            }
        }
        await mysqlDb.commit();

        console.log('事务已提交');
    } catch (e) {
        await mysqlDb.rollback();
        console.error('事务回滚:', e);
        throw e
    }
    process.on('exit', () => {
        sqliteDb.close();
        mysqlDb.end();
    });

})();

function combineNames(arrays) {
    const result = [];

    function backtrack(path, depth) {
        if (depth === arrays.length) {
            result.push(path.join("/"));
            return;
        }

        for (let i = 0; i < arrays[depth].length; i++) {
            let str = arrays[depth][i].name.zh_hk;
            path.push(`${str}`);  // 将索引转为字符串
            backtrack(path, depth + 1);
            path.pop();
        }
    }

    backtrack([], 0);
    return result;
}
function combineIndexes(arrays) {
    const result = [];

    function backtrack(path, depth) {
        if (depth === arrays.length) {
            result.push([...path]);
            return;
        }

        for (let i = 0; i < arrays[depth].length; i++) {
            path.push(`${i}`);  // 将索引转为字符串
            backtrack(path, depth + 1);
            path.pop();
        }
    }

    backtrack([], 0);
    return result;
}
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