var mqtt = require('mqtt');
var client = mqtt.connect('mqtt://broker');

client.on('connect', function () {
    console.log('connected!!');
});
