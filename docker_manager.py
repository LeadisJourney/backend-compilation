#!/usr/bin/env python2

import enum
import io
import os
import select
import socket
import time
from docker import Client

client_version = "1.18"

class Language(enum.Enum):
    C =     "C"
    CPP =   "CPP"

class Request(enum.Enum):
    EXEC =      "EXEC"
    SYNTAX =    "SYNTAX"

class UnixSocket:
    socket =    None
    filename =  None
    clients =   []

    def __init__(self, filename, listen):
        if os.path.exists(filename):
            os.remove(filename)
        self.socket = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        self.socket.bind(filename)
        self.socket.listen(listen)
        self.filename = filename
        self.clients = []

    def accept(self):
        socket, address = self.socket.accept()
        self.clients.append((socket, address))
        return (socket, address)

    def __del__(self):
        for (socket, _) in self.clients:
            socket.close()
        if self.socket:
            self.socket.close()
        if os.path.exists(self.filename):
            os.remove(self.filename)

class User:
    TIMEOUT =           600 # 10 minutes.
    name =              None
    client =            None
    time =              0
    container =         None
    volume =            None
    unixSocket =        None
    containerSocket =   None
    new_files_queue =   []
    busy =              False

    def __init__(self, name, client):
        self.name = name
        self.client = client
        self.time  = time.time()
        res = client.create_container(image="leadis_image", detach=True, volumes="/root/host", command="/root/server.py")
        self.container = res["Id"]
        self.volume = client.inspect_container(container=self.container)["Volumes"]["/root/host"]
        self.unixSocket = UnixSocket(self.volume + "/host.sock", 1)
        self.containerSocket = None
        self.new_files_queue = []
        self.busy = False
        client.start(container=self.container)

    def __del__(self):
        if self.container:
            self.client.stop(container=self.container, timeout=10)
            self.client.remove_container(container=self.container, v=True)

    def update(self, rlist, wlist, xlist):
        # User activity
        if (time.time() - self.time) >= self.TIMEOUT:
            # Timeout
            return None

        # Container-host connections
        if not self.containerSocket and self.unixSocket.socket not in rlist and self.unixSocket.socket not in xlist:
            # Add unix socket
            return (self.unixSocket.socket, None, None)
        elif not self.containerSocket and (self.unixSocket.socket in rlist or self.unixSocket.socket in xlist):
            # Connect container to the unix socket and continue user requests processing (if exists)
            self.containerSocket, _ = self.unixSocket.accept()

        # User requests and container responses
        if self.containerSocket and not self.busy and len(self.new_files_queue) > 0 and self.containerSocket not in wlist:
            # Need to write to container
            return (None, self.containerSocket, None)
        elif self.containerSocket and not self.busy and len(self.new_files_queue) > 0 and self.containerSocket in wlist:
            # Can write to container
            data, language, request, request_id = self.new_files_queue[0]
            if language == Language.C:
                with open(self.volume + "/main.c", "w") as file:
                    file.write(data)
            elif language == Language.CPP:
                with open(self.volume + "/main.cpp", "w") as file:
                    file.write(data)
            else:
                raise NotImplementedError("Not implemented language!")
            self.containerSocket.send(b"{} {} {}".format(language.value, request.value, request_id))
            self.new_files_queue = self.new_files_queue[1:]
            self.busy = True
            return (self.containerSocket, None, None)
        elif self.containerSocket and self.busy and self.containerSocket not in rlist and self.containerSocket not in xlist: 
            # Still waiting for container responses
            return (self.containerSocket, None, None)
        elif self.containerSocket and self.busy and (self.containerSocket in rlist or self.containerSocket in xlist): 
            # Need to read some thing from container
            response = self.containerSocket.recv(512)
            if len(response) > 0:
                print "Request No.: {}".format(response)
                if os.path.exists(self.volume + "/stderr"):
                    print "stderr:"
                    with open(self.volume + "/stderr") as file:
                        print file.read()
                    os.remove(self.volume + "/stderr")
                if os.path.exists(self.volume + "/stdout"):
                    print "stdout:"
                    with open(self.volume + "/stdout") as file:
                        print file.read()
                    os.remove(self.volume + "/stdout")
                self.busy = False
                if len(self.new_files_queue) > 0:
                    return (None, self.containerSocket, None)
                else:
                    return (None, None, None)
            else:
                # Error with the container
                return None
        else:
            # Nothing to do
            return (None, None, None)

    def new_file(self, code, language, request, request_id):
        self.time = time.time()
        self.new_files_queue.append((code, language, request, request_id))

import sys
import uuid
def main():
    # Docker client initialization.
    client = Client(version=client_version)
    for line in client.build(path=os.getcwd(), tag="leadis_image", rm=True):
        print line
    try:
        users = []
        rlist, wlist, xlist = [sys.stdin], [], []
        r, w, x = [], [], []
        while True:
            # Loop through all users, and update them.
            (r, w, x) = select.select(rlist, wlist, xlist, 60);
            rlist, wlist, xlist = [sys.stdin], [], []
            if sys.stdin in r or sys.stdin in x:
                # New stdin entry (temporary)
                entry = sys.stdin.readline()
                if entry.startswith("NEW"):
                    users.append(User(entry[4:-1], client))
                elif entry.startswith("EXEC"):
                    request_id = uuid.uuid1().hex
                    print "request No.: {}".format(request_id)
                    for user in users:
                        if user.name == entry[5:-1]:
                            with open("main.cpp") as file:
                                code = file.read()
                                user.new_file(code, Language.CPP, Request.EXEC, request_id)
                elif entry.startswith("SYNTAX"):
                    request_id = uuid.uuid1().hex
                    print "request No.: {}".format(request_id)
                    for user in users:
                        if user.name == entry[7:-1]:
                            with open("main.c") as file:
                                code = file.read()
                                user.new_file(code, Language.C, Request.SYNTAX, request_id)
                elif entry.startswith("EXIT"):
                    return
                else:
                    print "Unkown command!"
            for user in users:
                # Update every user, and get these socket to listen/write to
                res = user.update(r, w, x)
                if not res:
                    # Remove user (warning: user is not collect here!)
                    users.remove(user)
                else:
                    new_r, new_w, new_x = res
                    if new_r:
                        rlist.append(new_r)
                    if new_w:
                        wlist.append(new_w)
                    if new_x:
                        xlist.append(new_x)
    except KeyboardInterrupt:
        return

if __name__ == "__main__":
    # Docker need root privileges.
    if os.geteuid() != 0:
        exit("You need to have root privileges to run this script.\nPlease try again, this time using 'sudo'.")
    main()
