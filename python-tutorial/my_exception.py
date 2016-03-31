class MyException(Exception):
    def __init__(self, value):
        self.value = value
    def __str__(self):
        return repr(self.value)

try:
    raise MyException(2*2)
except MyException, e:
    print 'My exception occurred, value:', e.value

raise MyException('oops!')
