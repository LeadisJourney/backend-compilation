#!/usr/bin/env python2

import enum
import os
import socket
import sys
import time

from docker import Client

gClient         = None
gClient_version = "1.18"
gUsers          = []

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
    TIMEOUT =             600 # 10 minutes.
    name =                None
    client =              None
    time =                0
    container =           None
    volume =              None
    unixSocket =          None
    containerSocket =     None
    new_requests_queue =  []
    busy =                False

    def __init__(self, name, client):
        self.name = name
        self.client = client
        self.time  = time.time()
        res = client.create_container(image="leadis_image", detach=True, volumes="/root/host", command="/root/server.py")
        self.container = res["Id"]
        self.volume = client.inspect_container(container=self.container)["Volumes"]["/root/host"]
        self.unixSocket = UnixSocket(self.volume + "/host.sock", 1)
        self.containerSocket = None
        self.new_requests_queue = []
        self.busy = False
        client.start(container=self.container)

    def __del__(self):
        if self.container:
            self.client.stop(container=self.container, timeout=10)
            self.client.remove_container(container=self.container, v=True)
            self.container = None

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
        if self.containerSocket and not self.busy and len(self.new_requests_queue) > 0 and self.containerSocket not in wlist:
            # Need to write to container
            return (None, self.containerSocket, None)
        elif self.containerSocket and not self.busy and len(self.new_requests_queue) > 0 and self.containerSocket in wlist:
            # Can write to container
            data, language, request, request_id = self.new_requests_queue[0]
            if language == Language.C:
                with open(self.volume + "/main.c", "w") as file:
                    file.write(data)
            elif language == Language.CPP:
                with open(self.volume + "/main.cpp", "w") as file:
                    file.write(data)
            else:
                raise NotImplementedError("Not implemented language!")
            self.containerSocket.send(b"{} {} {}".format(language.value, request.value, request_id))
            self.new_requests_queue = self.new_requests_queue[1:]
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
                if len(self.new_requests_queue) > 0:
                    return (None, self.containerSocket, None)
                else:
                    return (None, None, None)
            else:
                # Error with the container
                return None
        else:
            # Nothing to do
            return (None, None, None)

    def new_request(self, code, language, request_type, request_id):
        self.time = time.time()
        self.new_requests_queue.append((code, language, request_type, request_id))


def main():
    global gClient
    global gClient_version

    # Docker client initialization.
    gClient = Client(version=gClient_version)
    for line in gClient.build(path=os.getcwd(), tag="leadis_image", rm=True):
        print line
def exit():
    global gUsers

    # Collect all the users
    gUsers = []

def new_user(name):
    global gClient
    global gUsers

    gUsers.append(User(name, gClient))

def new_request(user_name, code, language, request_type, request_id):
    global gUsers

    for user in gUsers:
        if user.name == user_name:
            user.new_request(code, language, request_type, request_id)

def update(rlist, wlist, xlist):
    global gUsers

    new_rlist, new_wlist, new_xlist = [], [], []
    for user in gUsers:
        # Update every user, and get their socket to listen/write to
        res = user.update(rlist, wlist, xlist)
        if not res:
            # Remove user (warning: user is not collect here!)
            gUsers.remove(user)
        else:
            r, w, x = res
            if r:
                new_rlist.append(r)
            if w:
                new_wlist.append(w)
            if x:
                new_xlist.append(x)
    return (new_rlist, new_wlist, new_xlist)

def info(command):
    global gClient
    global gClient_version
    global gUsers
    def users():
        res = ""
        for user in gUsers:
            res += user.name + "\n"
        return res
    def user(name):
        found = False

        for user in gUsers:
            if user.name == name:
                sys.stdout.write(user.name + ": \n" + \
                                 "container: " + user.container + "\n" + \
                                 "volume: " + user.volume + "\n" + \
                                 "pending requests: " + str(len(user.new_requests_queue)) + "\n")
                found = True
        if not found:
            sys.stderr.write("Unknown user!\n")

    switch = {
        "help" : lambda: sys.stdout.write("help         Print this help\n" + \
                                          "users        Users names\n" + \
                                          "user NAME    User Infos\n"),
        "users": lambda: sys.stdout.write(users()),
        "user" : lambda: user(command[1]) if len(command) == 2 else sys.stderr.write("Invalid syntax!\n")
    }
    switch.get(command[0], lambda: sys.stderr.write("Unknown command! Try 'help'.\n"))()
