# -*- coding: utf-8 -*-
# simple utilities

import os

from subprocess import Popen, PIPE

def ReadFile(filename):
    data = None
    with open(filename, 'r') as f:
        data = f.read()
    f.close()
    return data

def WriteFile(filename, data):
    with open(filename, 'w+') as f:
        f.write(str(data))
    f.close()

def CheckPrivilege():
    if os.getuid() != 0:
        print('''Warning: Running TiDB Insight with non-superuser privilege may result
         in lack of some information or data in the final output, if
         you find certain data missing or empty in result, please try
         to run this script again with root.''')

def CheckDir(path):
    try:
        os.mkdir(path)
        return path
    except OSError:
        if os.path.isdir(path):
            return path
    return None

# full directory path of this script
def pwd():
    return os.path.dirname(os.path.realpath(__file__))

# full path of current working directory
def cwd():
    return os.getcwd()

def run_cmd(cmd):
    p = Popen(cmd, stdout=PIPE, stderr=PIPE)
    return p.communicate()

def ParseCmdLine(cmdline):
    result = {}
    cmd = cmdline.split()
    for arg in cmd:
        # parse args that start with '--something'
        if arg.startswith("--"):
            argkv = arg.split("=")
            try:
                result[argkv[0][2:]] = argkv[1]
            except IndexError:
                pass
    return result