class BaseClass:
    def doit(self):
        print "Base class is being called"

class SubClass(BaseClass):
    def doit(self):
        BaseClass.doit(self)
        print "Now this is subclass"

    def dobaseclass(self):
        BaseClass.doit(self)

class AnotherClass:
    def do_something(self):
        print "doing something"

class MultipleParentsClass(BaseClass, AnotherClass):
    def do_alot(self):
        self.doit()
        self.do_something()

print "Instantiate SubClass..."
SubClass().doit();
print "Instantiate MultipleParentsClass..."
MultipleParentsClass().do_alot()

print "SubClass isinstance BaseClass =",
print isinstance(SubClass(), BaseClass)

print "SubClass isinstance AnotherClass =",
print isinstance(SubClass(), AnotherClass)

print "SubClass issubclass BaseClass =",
print issubclass(SubClass, BaseClass)

print "SubClass issubclass AnotherClass =",
print issubclass(SubClass, AnotherClass)

print "MultipleParentsClass issubclass BaseClass =",
print issubclass(MultipleParentsClass, BaseClass)

print "MultipleParentsClass issubclass AnotherClass =",
print issubclass(MultipleParentsClass, AnotherClass)

print "MultipleParentsClass issubclass SubClass =",
print issubclass(MultipleParentsClass, SubClass)
