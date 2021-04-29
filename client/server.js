/*
 sample from https://gist.github.com/ryanoglesby08/1e1f49d87ae8ab2cabf45623fc36a7fe
*/

const express = require('express');
const path = require('path');
const port = process.env.PORT || 5000;
const app = express();

// serve static assets normally
app.use(express.static('/usr/share/ui-dist'));

// handle every other route with index.html, which will contain
// a script tag to your application's JavaScript file(s).
app.get('*', function (_, response) {
  response.sendFile(path.resolve('/usr/share/ui-dist/index.html'));
});

app.listen(port);
console.log("server started on port " + port);