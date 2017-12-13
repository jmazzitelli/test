BUILD_DIR=$(shell readlink -f .)/build
JAR_FILE=${BUILD_DIR}/testjavaagent.jar
WILDFLY_HOME=${BUILD_DIR}/wildfly-home
JAVAAGENT_OPT=-javaagent:${JAR_FILE}
DEBUG_OPT=-agentlib:jdwp=transport=dt_socket,address=8787,server=y,suspend=y

compile:
	echo "Premain-Class: TestJavaAgent" > ${BUILD_DIR}/MANIFEST.MF
	rm -f *.class
	rm -f ${JAR_FILE}
	javac *.java
	jar cfm ${JAR_FILE} ${BUILD_DIR}/MANIFEST.MF *.class
	rm -f *.class

run:
	${WILDFLY_HOME}/bin/standalone.sh

kill:
	pkill -f jboss-modules.jar

put-server-in-reload-required:
	${WILDFLY_HOME}/bin/jboss-cli.sh --connect '/subsystem=infinispan/cache-container=ejb/:write-attribute(name=statistics-enabled,value=true)'

download-wildfly:
	mkdir -p ${BUILD_DIR}
	wget --output-document ${BUILD_DIR}/wildfly.zip 'http://download.jboss.org/wildfly/11.0.0.Final/wildfly-11.0.0.Final.zip'

unzip-wildfly:
	rm -rf ${BUILD_DIR}/wildfly*/
	unzip -d ${BUILD_DIR} ${BUILD_DIR}/wildfly.zip
	mv ${BUILD_DIR}/wildfly*/ ${WILDFLY_HOME}
	echo JAVA_OPTS=\"\$$JAVA_OPTS ${JAVAAGENT_OPT}\" >> ${WILDFLY_HOME}/bin/standalone.conf
