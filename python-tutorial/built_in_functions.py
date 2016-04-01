print "===right justifies string, with ~ as the fill character"
print format("gorilla?", "~>15")

print "===getattr"
class MyClassWithAttrib:
    datum = "datum value"
print getattr(MyClassWithAttrib(), "datum")
print hasattr(MyClassWithAttrib(), "datum")
print hasattr(MyClassWithAttrib(), "nothing")

print "===Get all the symbols from the global symbol table"
all_globals = dict(globals())
for k, v in enumerate(all_globals):
    print "Global symbol {:>20} = {}".format(v, all_globals[v])
    
print "===hash"
print hash(1.0)
print hash(1)
print hash("hello")

print "===identity"
a = "hello world"
b = "hello" + " " + "world"
c = a
print "a={}, b={}, c={}".format(a, b, c)
print "ids: a={}, b={}, c={}".format(id(a), id(b), id(c))
print "equality: a==b{}, b==c{}, a==c{}".format(a == b, b == c, a == c)
print "identity: a is b={}, b is c={}, a is c={}".format(a is b, b is c, a is c)

print "===string to int"
print int("123") + int("27")

print "=== get length of things"
print "Length of string = ", len("hello")
print "Lengh of array = ", len([1, 2, 3])
print "Length of dict = ", len({"one":"1", "two": "2"})
print "Length of set = ", len(set(["a", "a", "b", "b"])) # set really has two items
print "Length of range = ", len(range(10))


print "===Get all the symbols from the local symbol table"
all_locals = dict(locals())
for k, v in enumerate(all_locals):
    print "Local symbol {:>20} = {}".format(v, all_locals[v])

print "===map function"
def multiply_numbers(*x):
    if (len(x) == 0):
        return 0;
    result = 1;
    for n in x:
        result *= n;
    return result

for i in map(multiply_numbers, [1, 2, 3, 4]):
    print i

squares = map(multiply_numbers, [1,2,3], [1,2,3]) # essentially returns the squares
print str(squares)

cubes = map(multiply_numbers, [1,2,3], [1,2,3], [1,2,3]) # essentially returns the cubes
print str(cubes)

print "===reduce"
print reduce(lambda x, y: x+y, [1, 2, 3, 4, 5])

