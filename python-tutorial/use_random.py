import random

food = ['Apple','Banana','Hamburger','Hotdog','Lobster','Steak','Pasta']
drink = ['Water','Lemonade','Beer','Wine','Soda','Coffee','Tea']

for i in range(10):
    pick_food = random.randint(0, len(food)-1)
    pick_drink = random.randint(0, len(drink)-1)
    print "{}:{}".format(pick_food, pick_drink)
    print "{}:{}".format(food[pick_food], drink[pick_drink])
