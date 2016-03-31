import golf.woods.types
import golf.irons.types
from golf.putters import types

name_types ={
              "Woods": golf.woods.types.kinds,
              "Irons": golf.irons.types.kinds,
              "Putters": types.kinds
             }

for n, t in name_types.iteritems():
    print "{name:8}: {types}".format(name=n, types=t)
