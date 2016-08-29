var http = require('http');

var server = http.createServer();

server.on('request', function(request, response) {
 response.writeHead(200, {
        'Content-Type': 'text/plain; charset=utf-8',
        'Transfer-Encoding': 'chunked',
        'X-Content-Type-Options': 'nosniff'});
    response.write('Beginning\n');
    var count = 10;
    var io = setInterval(function() {
        response.write('Doing ' + count.toString() + '\n');
        count--;
        if (count === 0) {
            response.end('Finished\n');
            clearInterval(io);
        }
    }, 200); //Last argument is the time in ms between writes - play around with this and observe what happens
});

server.listen(4545)
