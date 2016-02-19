'use strict';

const Promise = require('bluebird'),
    GitHub = require('github')

const REPO = "https://www.dartlang.org/downloads/archive/";
function dart(max = 10) {
    var github = new GitHub({
        // required
        version: "3.0.0",
        // optional
        debug: false,
        protocol: "https",
        //host: "github.my-GHE-enabled-company.com", // should be api.github.com for GitHub
        //pathPrefix: "/api/v3", // for some GHEs; none for GitHub
        timeout: 5000,
        headers: {
            "user-agent": "My-Cool-GitHub-App" // GitHub is happy with a unique user agent
        }
    });
    
    github.releases.listReleases({
        owner: 'dart-lang',
        repo: 'sdk'
    }, (e, r) => {
        console.log(e, r)
    })

}


dart()