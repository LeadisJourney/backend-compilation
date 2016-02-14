#!/usr/bin/env python3

import enum
import os
import socket
import subprocess
import time

class Language(enum.Enum):
    C =     "C"
    CPP =   "CPP"

class Request(enum.Enum):
    EXEC =      "EXEC"
    SYNTAX =    "SYNTAX"

socket = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
socket.connect("/root/host/host.sock")

while True:
    request = socket.recv(512).decode().split()

    if request[1] == Request.EXEC.value:
        # Execution
        if request[0] == Language.C.value:
            with open("/root/host/stderr", "w") as stderr:
                subprocess.call(["gcc", "/root/host/main.c", "-o", "/root/a.out"], stderr=stderr)
        elif request[0] == Language.CPP.value:
            with open("/root/host/stderr", "w") as stderr:
                subprocess.call(["g++", "/root/host/main.cpp", "-o", "/root/a.out"], stderr=stderr)
        else:    
            raise NotImplementedError()
        if os.path.exists("/root/a.out"):
            with open("/root/host/stdout", "w") as stdout:
                with open("/root/host/stderr", "w") as stderr:
                    subprocess.call(["/root/a.out"], stdout=stdout, stderr=stderr)
        socket.send(request[2].encode())
    elif request[1] == Request.SYNTAX.value:
        # Syntax analysis
        if request[0] == Language.C.value:
            with open("/root/host/stderr", "w") as stderr:
                subprocess.call(["gcc", "-fdiagnostics-color=never", "-fsyntax-only", "/root/host/main.c"], stderr=stderr)
        elif request[0] == Language.CPP.value:
            with open("/root/host/stderr", "w") as stderr:
                subprocess.call(["g++", "-fdiagnostics-color=never", "-fsyntax-only", "/root/host/main.cpp"], stderr=stderr)
        else:    
            raise NotImplementedError()
        socket.send(request[2].encode())
    else:
        raise NotImplementedError()
