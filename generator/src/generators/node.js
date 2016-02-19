
const request = require('request-promise');
const Promise = require('bluebird')
const repo_stable = "https://nodejs.org/download/release"
const repo_unstable = "https://nodejs.org/download/nightly"




module.exports = async function (max = 10) {
    
    let [stable, unstable] = await Promise.all([getVersions(repo_stable, max), getVersions(repo_unstable, 1)]);
    
    return {
        name: `Node`,
        description: `Node.js® is a JavaScript runtime built on Chrome's V8 JavaScript engine. Node.js uses an event-driven, non-blocking I/O model that makes it lightweight and efficient. Node.js' package ecosystem, npm, is the largest ecosystem of open source libraries in the world.`,
        stable: stable,
        unstable: unstable,
        export: {
            binary: "bin"
        }
    };
    
}

async function getVersions(repo, max = 10) {
    let manifest =  repo + "/index.json"
    let results = await request(manifest);
    
    results = JSON.parse(results);
    let output = []
    var found = 0;
    for (let i=0,ii=results.length;i<ii;i++) {
        let result = await getVersion(repo, results[i]);
        if (result && result.length) {
            output = output.concat(result);
            if (found == 0) {
                output.forEach( m => {
                    m.latest = true;
                });
                
            }
            found++;
        }
        
        if (found >= max) {
            break;
        }
        
    }
    
    return output;
}

async function getVersion(repo, v) {
    let output = [];
    
    let sha256 = await request(repo + "/" + v.version + "/SHASUMS256.txt")
    
    let macFound = false
    
    for (let i=0,ii=v.files.length;i<ii;i++) {
        
        if (v.files[i].match(/osx-(x64|x86)-pkg/) || v.files[i].indexOf('win') != -1 || v.files[i] == 'headers') {
            continue;
        }
        
        
        let tmpFile = v.files[i].replace('-tar', '').replace('osx', 'darwin')
        
        let file = tmpFile == "src" ? `node-${v.version}.tar.gz` : `node-${v.version}-${tmpFile}.tar.gz`;
        
        let reg = new RegExp('([0-9-a-z]{64})\\s{2}' + file)
        
        let match = sha256.match(reg)
        
        if (match == null) {
            console.log('No match for %s', file)
            continue;
        }
        
        let [os, arch] = tmpFile.split('-')
        
        let version = {
            version: v.version,
            os: tmpFile == 'src' ? void 0 : os ,
            arch: tmpFile == 'src' ? void 0 : arch,
            binary: v.files[i] !== 'src',
            latest: false,
            source: {
                type: 'URL',
                link: `${repo}/${v.version}/${file}`,
                hash: {
                    type: 'sha256',
                    value: match[1]
                },
                target: file.replace('.tar.gz', '')
            }
        }
        
        if (v.files[i] === 'src') {
           
            version.build = [
                {
                    "interpreter": "python",
                    "cmd": "configure --prefix={{.Source}}"
                }, {
                    "cmd": "make"
                }, {
                    "cmd": "make install"
                }
            ];
        }
        
        
        output.push(version)
        
    }
    
    return output;
    
}