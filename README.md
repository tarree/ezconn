# ezconn
Connection manager with proxy and user-agent written in go

## Installation

You need go 1.7+

```bash
go get github.com/tarree/ezconn/...
```

### For developers

You need golang installed (it's been written with the latest version). You also need glide to build this app. Simply do the following:
```bash
go get -u github.com/Masterminds/glide
glide install
```

## Usage

The last argument (without any flags) is the url.

You can see command line flags by executing:
```bash
./ezconn -h
  -logLevel string
        Log level. Values are: debug, info, warn, error. Default level is error (default "error")
  -o string
        Saves the output to a file.
  -proxy string
        Sets the proxy. Direct connection if you don't set it.
  -timeout duration
        Sets the timeout on request. (default 5s)
  -user-agent string
        Sets the user-agent on client. Default is golang's default user-agent
```

### Examples

```bash
# Setting proxy:
./ezconn -proxy PROXY_URL:PORT http://example.com

# Setting user-agent
./ezconn -proxy PROXY_URL:PORT http://example.com

# Output to a file, 3 seconds timeout:
./ezconn -proxy PROXY_URL:PORT -timeout 3s -o filename.html http://example.com

# Show debug information:
./ezconn -proxy PROXY_URL:PORT -timeout 3s -o filename.html -logLevel debug http://example.com
```

#### Note

All arguments except for the target url can be set in environment variable as well:

```bash
export PROXY=PROXY_URL:PORT
export TIMEOUT=3s
export LOGLEVEL=info

./ezconn -o filename.html http://example.com
./ezconn -o filename2.html http://example.co.uk
```

# License
MIT licence
