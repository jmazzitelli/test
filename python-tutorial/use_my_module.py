# notice you have to prefix symbols from the imported module with the module name
import my_module

my_module.fib(100)

print "\n" + str(dir(my_module))

# using from...import means you don't have to prefix the name
from my_module import fib2

print str(fib2(50))
