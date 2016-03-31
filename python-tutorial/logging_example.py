import logging

logging.debug("a debug msg")
logging.info("an info msg")
logging.warn("a warn msg")
logging.warning("a warning msg")
logging.error("an error msg")
logging.critical("a critical msg")
logging.fatal("a fatal msg")

try:
    print 2/0
except:
    logging.warn("In except")
    logging.exception("Got exception!: ")
    logging.warn("Will now raise the exception again")
    raise
