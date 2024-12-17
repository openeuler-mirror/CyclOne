var fs = require('fs-extra'),
    path = require('path'),
    exec = require('child_process').exec,
    BUILDENV = 'production',
    currentBuidEnv = '';
module.exports = {
    name: 'publish',
    isDevelopingAddon: function() {
        return true;
    },
    outputReady: function(results) {
        var env = process.env.EMBER_ENV || currentBuidEnv;
        console.log('in ' + env + ' mode');
        if (env === BUILDENV) {
            var bin = path.resolve(__dirname, '../../bin');
            var dist = path.resolve(__dirname, '../../dist');

            console.log('\n Copy dist to bin.... \n');
            try {
                // remove from git
                exec('git rm -r --cached --ignore-unmatch ' + bin, function (error, stdout, stderr) {
                    console.log(' Remove bin from repository \n');
                    if (!!stdout) {
                        console.log('stdout: ' + stdout);
                    }
                    if (!!stderr) {
                        console.log('stderr: ' + stderr);
                    }
                    if (error !== null) {
                        console.log('exec error: ' + error);
                    }
                });
            } catch (err) {
                console.error(err);
            }
            
            fs.emptyDirSync(bin);
            fs.copySync(dist, bin);
        }
    },
    config: function(env, baseConfig) {
        currentBuidEnv = env;
    }
};
