'use strict';

require('./lib').generate()
.then( () => {
    console.log('done')
}).catch( e => {
    console.error(e)
})