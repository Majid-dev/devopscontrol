const http = require('http');
const port = process.env.PORT || 3000;
http.createServer((req, res) => {
  res.end("Hello from DevOpsControl!");
}).listen(port, () => console.log(`Listening on ${port}`));
