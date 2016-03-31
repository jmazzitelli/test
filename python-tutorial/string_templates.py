from string import Template
t = Template('${village}folk send $$10 to $cause.')
print t.substitute(village='Nottingham', cause='the ditch fund')

print t.safe_substitute(village="Nowhere")

try:
    print t.substitute(village="Nowhere")
except Exception, e:
    print "not safe: " + str(e)
