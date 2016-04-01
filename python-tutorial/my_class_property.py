class MyClassWithProperty(object):
    def __init__(self):
        self._x = "<not set>"

    @property #this is the getter
    def x(self):
        """I'm the 'x' property."""
        print "Getting x"
        return self._x

    @x.setter
    def x(self, value):
        print "Setting x"
        self._x = value

    @x.deleter
    def x(self):
        print "Deleting x"
        self._x = "<was deleted>"

myc = MyClassWithProperty()
myc.x = "foo" # calls setter
print "x=" + myc.x # calls getter
del myc.x # calls deleter
print "x=" + myc.x # calls getter
