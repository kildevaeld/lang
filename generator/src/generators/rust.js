
const Promise = require('bluebird'),
    request = require('request-promise');



const REPO = "https://static.rust-lang.org/dist/";

module.exports = async function generate (max = 10) {
    
    let [stable, unstable] = await Promise.all([getVersions("stable", max), getVersions("beta", 1)]);
    
    return {
        name: `Rust`,
        description: '',
        stable: stable,
        unstable: unstable,
        export: {
            binary: "bin",
            library: "lib"
        }
    };
    
}

async function getVersions(channel, max = 10) {
    
    let url = `${REPO}channel-rust-${channel}`;
    
    let files = await request(url);
    let sha256 = await request(url + ".sha256")
    
    let split = files.split('\n');
    let output = [];
    for (let i=0, ii = split.length; i<ii;i++) {
        let version = {
            version: split[i].split('-')[1],
            binary: true,
            latest: channel === 'stable' ? true : false,
            source: {
                type: 'URL',
                link: REPO + split[i],
                target: split[i].replace('.tar.gz', '')
            },
            build: [
                {
                    "exec" : "./install.sh --prefix={{.Source}} --verbose"
                }
            ]
        }
        
        if (split[i].indexOf("darwin.pkg") != -1) {
            continue;
        }
        
        if (split[i].indexOf('linux') != -1) {
            version.os = "linux";
        } else if (split[i].indexOf('darwin') != -1) {
            version.os = "darwin";
        } else {
            continue;
        }
        
        if (split[i].indexOf("x86_64") != -1) {
            version.arch = "x64";
        } else if (split[i].indexOf("i686") != -1) {
            version.arch = "x86";
        } else {
            continue;
        }
        
        /*let reg = new RegExp('([0-9-a-z]{64})\\s+' + split[i])
        
        let match = sha256.match(reg)
        
        if (match == null) {
            console.log('No match for %s', split[i])
            continue;
        }*/
        
        //version.source.hash.value = match[1];
        output.push(version)
    }
    
    return output;
    
}
