class MyClass:
    """A simple class"""
    class_var = 50
    
    def foo(self, msg=None):
        print "In foo function: " + str(self) + "; i=" + str(self.i)
        if (msg is not None):
            print msg

    def __init__(self, adjustment = 0):
        self.i = 50;
        self.i += adjustment
        self.foo()

my_class1 = MyClass()
my_class1.i = 60
my_class1.foo()

my_class1.i = 51
method_pointer = my_class1.foo
method_pointer("from method ptr")

my_class2 = MyClass(1)
print "class_var=" + str(MyClass.class_var)
print "class #1 class_var=" + str(my_class1.class_var)
print "class #2 class_var=" + str(my_class2.class_var)

my_class2.class_var += 1
MyClass.class_var += 1
print "class_var=" + str(MyClass.class_var)
print "class #1 class_var=" + str(my_class1.class_var)
print "class #2 class_var=" + str(my_class2.class_var)

