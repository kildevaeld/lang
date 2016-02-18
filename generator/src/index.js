'use strict';

const fs = require('mz/fs');
const Path = require('path');


export async function generate () {
    let path = Path.join(__dirname,"generators");
    let files = await fs.readdir(path);
    var result = {};
    for (let i=0, ii = files.length; i<ii;i++) {

        let fullp = Path.join(path, files[i]);

        let fn = require(fullp);

        let name = files[i].substr(0, files[i].indexOf("."))

        console.log("Generating %s", name)
        try {
            result[name] = await fn(50);
        } catch (e) {
            console.log('Got error %s', e.stack)
        }

    }

    await fs.writeFile('manifest.json', JSON.stringify(result, null, 2));

}