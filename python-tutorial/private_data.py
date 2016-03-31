# There is no such thing as private data members/methods in Python.
# Convention is variables with _ are to be considered not part of the API.
# You can name-mangle variables with double-underscore with at most one trailing "_"
# such as __name_mangle_this or __protect_me_

class MyClass:
    def _do_not_call_me(self):
        print "Oh well. You found me."
    
    _private_member = "You found me, too."
    
    __name_mangle_me = "Peek-a-boo"
    
    __hideme_ = "boo!"

myc = MyClass()
myc._do_not_call_me()
print myc._private_member
print myc._MyClass__name_mangle_me
print myc._MyClass__hideme_
