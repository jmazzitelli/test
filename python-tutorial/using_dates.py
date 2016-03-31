import datetime

now = datetime.date.today()
print "Today Date:" + now.__str__()

now = datetime.datetime.today()
print "Today Time:" + datetime.datetime.isoformat(now)

print now.strftime("%m/%d/%y %H:%M:%S\n%d %b %Y is a %A on the %d day of %B.")
