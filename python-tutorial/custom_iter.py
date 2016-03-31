class CanIterate:
    def __init__(self, *data):
        if (len(data) == 0):
            data = [1, 2, 3]
        self.data = data
    
    def __iter__(self):
        class CanIterateIterator:
            def __init__(self, thedata):
                self.index = len(thedata)
                self.thedata = list(reversed(thedata))
                
            def next(self):
                if (self.index == 0):
                    raise StopIteration
                self.index = self.index - 1;
                return self.thedata[self.index]

        return CanIterateIterator(self.data)

for i in CanIterate():
    print i

print "Do it again"
for i in CanIterate(5,4,3,2,1):
    print i


