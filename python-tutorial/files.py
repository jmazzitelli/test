# this closes the file automatically
with open("/etc/hosts") as f:
    print "fileno=", f.fileno()
    for line in f:
        print line,
    print "fpos=", f.tell()

