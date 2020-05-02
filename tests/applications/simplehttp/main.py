from http.server import HTTPServer, BaseHTTPRequestHandler
"""
 This program will count the number of HTTP requests recieved,
 returning the current count to every request made
  Parameters: none
  Return: none (unreachable during normal execution)
"""

counter = 0

class CounterRequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        global counter
        counter += 1

        self.send_response(200)
        self.end_headers()
        return_str = str(counter) + "\n"
        self.wfile.write(return_str.encode('utf-8'))

    def do_POST(self):
        global counter
        counter += 1

        self.send_response(200)
        self.end_headers()
        return_str = str(counter) + "\n"
        self.wfile.write(return_str.encode('utf-8'))

    do_PUT = do_POST
    do_DELETE = do_GET

def main():
    port = 8080
    print('Listening on 0.0.0.0:%s' % port)
    server = HTTPServer(('', port), CounterRequestHandler)
    server.serve_forever()

if __name__ == "__main__":
    main()
