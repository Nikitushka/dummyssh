# dummysshbruter
A little project of mine for the intensive Go-programming course taught by Tero Karvinen
http://terokarvinen.com/2020/go-programming-course-2020-w22/

# whatis dummysshbruter
This little project in essence is a very basic ssh shell that establishes a connection to host and can run simple commands,
like pwd, ls & whoami by brute-forcing the remote host with username and password lists

# Usage
If you have go installed:

`go build`

Alternatively, you can run the pre-built binary in this repo and see the list of flags:

```
./dummyssh --help

-h string
    remote host (default "127.0.01")
-P string
    password file
-U string
    username file
```

Example:

```
./dummyssh -U users.txt -P passwords.txt -h xxx.xxx.xxx.xxx
Connection established!
pwd
/home

exit

Exiting...
```


# Sources

http://blog.ralch.com/tutorial/golang-ssh-connection/

https://golangcode.com/get-the-current-username-name-and-home-dir-cross-platform/

https://github.com/Nikitushka/HelloGo/blob/master/helloflags/hello.go

https://golang.org/pkg/bytes/

https://gist.github.com/svett/b7f56afc966a6b6ac2fc#gistcomment-2823834

Answer by egorka:

https://stackoverflow.com/questions/24440193/golang-ssh-how-to-run-multiple-commands-on-the-same-session

https://stackoverflow.com/questions/35110610/whats-the-right-way-to-clear-a-bytes-buffer-in-golang

https://gobyexample.com/reading-files

https://stackoverflow.com/questions/8757389/reading-a-file-line-by-line-in-go

flag handling:
https://github.com/ffuf/ffuf/blob/master/main.go

https://golang.org/pkg/net/

https://medium.com/@0x766c6164/writing-a-simple-ssh-brute-forcer-in-go-19c4f928cd3b

form the article above:
https://github.com/vlad-s/gofindssh/blob/master/main.go#L86
