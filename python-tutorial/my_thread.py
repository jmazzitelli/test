import threading
from multiprocessing.queues import Queue

q = Queue()
for s in ["First", "Second", "Third"]: q.put(s)

class MyThread(threading.Thread):
    def run(self):
        while True:
            job_text = q.get()
            if (job_text == "quit"):
                return
            print "Job:" + job_text
        
thread = MyThread()
thread.start()

q.put("Another One")
q.put("quit")
thread.join()
print "Done"
