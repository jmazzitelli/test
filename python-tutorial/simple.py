#/usr/bin/python

print "--if statement"
x = int("1") #raw_input("Please enter an integer: "))
if x < 0:
    x = 0
    print 'Negative changed to zero'
elif x == 0:
    print 'Zero'
elif x == 1:
    print 'Single'
else:
    print 'More'

print "--for loop"
words = ['cat', 'window', 'defenestrate']
for w in words[:]:
    if len(w) > 5:
        words.insert(0, w)
print words

print "--range"
a = ['Mary', 'had', 'a', 'little', 'lamb']
for i in range(len(a)):
    print i, a[i]
for k,v in enumerate(a):
    print("key={},value={}").format(k, v)

print "--for...else"
for i in range(5):
    print i
else:
    print "done at " + str(i)

print "--function"
def fib(n):
    """Print a Fibonacci series up to n."""
    a, b = 0, 1
    while a < n:
        print a,
        a, b = b, a+b
fib(2000)
print "\npass argument with key=value"
fib(n=100)

print "--varargs"
def dovar(*args, **diction):
    for i in sorted(args):
        print "argument: " + i
    for key in sorted(diction.keys()):
        print "{}={}".format(key, diction[key])
dovar("B.one", "A.two", B_first=111, A_second=222)

print "--lamba"
def make_incrementor(n):
    """Create function that sums two numbers.
    Returns a lamba that takes a single number and
    adds it to a constant.
    """
    return lambda x: x + n
f = make_incrementor(42)
print f(0)
print f(1)

print "--matrix"
matrix = [ [1,2,3], [4,5,6] ]
for i in matrix:
    print i

print "--tuples"
arr = [1, 2, 3]
print arr
del arr[1]
print arr
tuples = 1, 2, 3, "hello"
print tuples
print tuples[3]
try:
    del tuples[1] # fails - tuples are immutable
except:
    print "tuples are immutable!"

print "--sets"
arr = [ "apple", "apple", "orange", "orange" ]
s = set(arr)
print s
print "apple" in s
print "crabgrass" in s

print "--dictionary"
diction = {
  "one": 111,
  'two': 222,
  "three": "three three three"
}
print diction
print diction['two']
print diction
del diction['two']
print diction
try:
    print diction['two']
except KeyError, ke:
    print "we deleted 'two': " + str(   ke.__class__)

print "--zip"
questions = ['name', 'quest', 'favorite color']
answers = ['lancelot', 'the holy grail', 'blue']
for q, a in zip(questions, answers):
    print 'What is your {0}?  It is {1}.'.format(q, a)



