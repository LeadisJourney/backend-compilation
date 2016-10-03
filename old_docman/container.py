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
    EXEC =      "EXECUTION"
    SYNTAX =    "COMPILATION"

socket = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
socket.connect("/root/host/host.sock")

while True:
    request = socket.recv(512).decode().split()

    if request[1] == Request.EXEC.value:
        # Execution
        # subprocess.call(["make", "realclean", "-C", "/root/host/exercise"], stdout=stdout, stderr=stderr)
        if request[0] == Language.C.value:
            with open("/root/host/stdout", "w") as stderr:
                with open("/root/host/stdout", "w") as stdout:
                    subprocess.call(["make", "-C", "/root/host/exercise"], stdout=stdout, stderr=stderr)
        elif request[0] == Language.CPP.value:
            with open("/root/host/stdout", "w") as stderr:
                with open("/root/host/stdout", "w") as stdout:
                    subprocess.call(["make", "-C", "/root/host/exercise"], stdout=stdout, stderr=stderr)
        else:    
            raise NotImplementedError()
        if os.path.exists("/root/host/exercise/res"):
            os.chdir("/root/host/exercise/")
            with open("/root/host/stdout", "w") as stdout:
                with open("/root/host/stdout", "w") as stderr:
                    subprocess.call(["./res"], stdout=stdout, stderr=stderr)
        socket.send(b'1')
    elif request[1] == Request.SYNTAX.value:
        # Syntax analysis
        if request[0] == Language.C.value:
            with open("/root/host/stdout", "w") as stderr:
                subprocess.call(["gcc", "-fdiagnostics-color=never", "-fsyntax-only", "/root/host/main.c"], stderr=stderr)
        elif request[0] == Language.CPP.value:
            with open("/root/host/stdout", "w") as stderr:
                subprocess.call(["g++", "-fdiagnostics-color=never", "-fsyntax-only", "/root/host/main.cpp"], stderr=stderr)
        else:    
            raise NotImplementedError()
        socket.send(b'1')
    else:
        raise NotImplementedError()
