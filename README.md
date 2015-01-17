servelast
=========
Serves the newest file in `basepath` or any subdirectory, updated every 5 seconds.

Might not be the best idea to use on huge directories, I have not tested it.

`basepath` specified via `-basepath` argument
Filename is passed as the HTTP header `X-Filename`
