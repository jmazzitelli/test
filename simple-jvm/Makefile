JAR_FILE=simple-jvm.jar

JAVAAGENT_FILE=$(shell ls -1 /home/mazz/source/hawkular-agent/hawkular-javaagent/target/hawkular-javaagent-*.jar)

JAVAAGENT_CONFIG_FILE=config-jmx.yaml

JAR_OPT=-jar lib/${JAR_FILE}

JAVAAGENT_OPT=-javaagent:${JAVAAGENT_FILE}=config=${JAVAAGENT_CONFIG_FILE},delay=2

DEBUG_OPT=-agentlib:jdwp=transport=dt_socket,address=8787,server=y,suspend=y

JOLOKIA_OPT=-javaagent:lib/jolokia.jar=port=9191,host=127.0.0.1

JAVA_OPTS=-Dsimple-flag

all:
	rm -f *.class
	javac *.java
	rm -f ${JAR_FILE}
	jar cfe ${JAR_FILE} SimpleJVM *.class
	mv ${JAR_FILE} ./lib

run:
	java ${JAVA_OPTS} ${JAR_OPT}

run-with-agent:
	java ${JAVA_OPTS} ${JAVAAGENT_OPT} ${JAR_OPT}

debug-with-agent:
	java ${JAVA_OPTS} ${DEBUG_OPT} ${JAVAAGENT_OPT} ${JAR_OPT}

run-with-agent-and-jolokia:
	java ${JAVA_OPTS} ${JAVAAGENT_OPT} ${JOLOKIA_OPT} ${JAR_OPT}

debug-with-agent-and-jolokia:
	java ${JAVA_OPTS} ${DEBUG_OPT} ${JAVAAGENT_OPT} ${JOLOKIA_OPT} ${JAR_OPT}

download-jolokia:
	mkdir -p ./lib
	wget --output-document ./lib/jolokia.jar 'http://search.maven.org/remotecontent?filepath=org/jolokia/jolokia-jvm/1.3.5/jolokia-jvm-1.3.5-agent.jar'
