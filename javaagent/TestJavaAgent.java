import javax.management.*;
import java.lang.management.*;
import java.util.*;

public class TestJavaAgent implements TestJavaAgentMXBean {
    public static void premain(String args) {
        Thread agentThread = new Thread(new Runnable() {
            public void run() {
                TestJavaAgent bean = null;
                try {
                    Thread.sleep(5000);
                    String name = "zztop:type=testjavaagent";
                    System.out.println("TestJavaAgent: Registering the MBean: " + name);
                    MBeanServer mbs = ManagementFactory.getPlatformMBeanServer();
                    bean = new TestJavaAgent();
                    mbs.registerMBean(bean, new ObjectName(name));
                } catch (Exception e) {
                    System.out.println("TestJavaAgent: Failed to register the MBean");
                    e.printStackTrace();
                }

                while (bean != null) {
                    try {
                        Thread.sleep(5000);
                        bean.doit();
                    } catch (Exception e) {
                        e.printStackTrace();
                    }
                }
            }
        }, "Test Java Agent Thread");
        agentThread.setDaemon(true);
        agentThread.start();
    }

    @Override
    public String doit() {
        StringBuilder str = new StringBuilder();
        try {
            ObjectName name = new ObjectName("jboss.as:management-root=server");
            MBeanServer mbs = ManagementFactory.getPlatformMBeanServer();

            str.append("=============================================================\n");
            str.append("FIND INFORMATION ABOUT MBEAN: " + name + "\n");
            str.append("=============================================================\n");

            str.append("isRegistered:\n");
            str.append(mbs.isRegistered(name) + "\n");
            str.append("getMBeanInfo:\n");
            try {
                MBeanInfo mbi = mbs.getMBeanInfo(name);
                str.append("  description: " + mbi.getDescription() + "\n");
                str.append("  #attributes: " + mbi.getAttributes().length + "\n");
            } catch (Exception e) {
                e.printStackTrace();
            }
            str.append("getAttribute:\n");
            try {
                Object o = mbs.getAttribute(name, "serverState");
                str.append("serverState=" + o + "\n");
            } catch (Exception e) {
                e.printStackTrace();
            }
            str.append("queryNames:\n");
            str.append(mbs.queryNames(name, null) + "\n");
            str.append("queryMBeans:\n");
            str.append(mbs.queryMBeans(name, null) + "\n");
            str.append("queryNames(null, null):\n");
            Set<ObjectName> beanNames = mbs.queryNames(null, null);
            for (ObjectName beanName : beanNames) {
                if (beanName.equals(name)) {
                    str.append("FOUND IT: " + beanName + "\n");
                }
            }
            str.append("=============================================================\n");
        } catch (Exception e) {
            e.printStackTrace();
        }
        System.out.println(str);
        return str.toString();
    }
}
