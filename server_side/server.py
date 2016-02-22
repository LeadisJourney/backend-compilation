#!/usr/bin/env python2
import os
import select
import sys

import docker_server

import uuid
def main():
    print "Type 'quit' to quit, or 'help' for docker command."
    sys.stdout.write(">")
    sys.stdout.flush()
    try:
        new_rlist, new_wlist, new_xlist = [sys.stdin], [], []
        while True:
            # Loop through all users, and update them.
            (rlist, wlist, xlist) = select.select(new_rlist, new_wlist, new_xlist, 60);
            new_rlist, new_wlist, new_xlist = [sys.stdin], [], []
            if sys.stdin in rlist or sys.stdin in xlist:
                # New stdin entry (temporary)
                entry = sys.stdin.readline()
                if entry.startswith("NEW"):
                    docker_server.new_user(entry[4:-1])
                elif entry.startswith("EXEC"):
                    request_id = uuid.uuid1().hex
                    print "Request No.: {}".format(request_id)
                    with open("main.cpp") as file:
                        code = file.read()
                        docker_server.new_request(entry[5:-1], code, docker_server.Language.CPP, docker_server.Request.EXEC, request_id)
                elif entry.startswith("SYNTAX"):
                    request_id = uuid.uuid1().hex
                    print "Request No.: {}".format(request_id)
                    with open("main.c") as file:
                        code = file.read()
                        docker_server.new_request(entry[7:-1], code, docker_server.Language.C, docker_server.Request.SYNTAX, request_id)
                elif entry.startswith("quit"):
                    return
                else:
                    docker_server.info(entry.split())
                    sys.stdout.write(">")
                    sys.stdout.flush()
            (r, w, x) = docker_server.update(rlist, wlist, xlist)
            new_rlist += r
            new_wlist += w
            new_xlist += x
    except KeyboardInterrupt:
        return

if __name__ == "__main__":
    # Docker need root privilege.
    if os.geteuid() != 0:
        exit("You need to have root privileges to run this script.\nPlease try again, this time using 'sudo'.")
    docker_server.main()
    main()
    docker_server.exit()
