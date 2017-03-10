import javax.management.*;
import java.lang.management.*;

public class SimpleJVM implements SimpleMXBean {
    public static final void main(String[] args) {
        try {
            MBeanServer mbs = ManagementFactory.getPlatformMBeanServer();
            mbs.registerMBean(new SimpleJVM(), new ObjectName("org.hawkular.test:type=simple"));
        } catch (Exception e) {
            System.out.println("Failed to register the MBean");
            e.printStackTrace();
        }

        synchronized (SimpleJVM.class) {
            System.out.println("Going in a coma...");
            try {
                SimpleJVM.class.wait();
            } catch (Exception e) {
            }
            System.out.println("Woke up from coma and exiting.");
        }
    }

   // JMX Attributes

    public String getTestString() { return "Test string from SimpleJVM MBean"; }

    public int     getTestIntegerPrimitive() { return 12345; }
    public Integer getTestInteger()          { return 54321; }

    public boolean getTestBooleanPrimitive() { return false; }
    public Boolean getTestBoolean()          { return true; }

    public long getTestLongPrimitive()       { return 123456789L; }
    public Long getTestLong()                { return 987654321L; }

    public double getTestDoublePrimitive()   { return (double) 1.23456789; }
    public Double getTestDouble()            { return (double) 9.87654321; }

    public float getTestFloatPrimitive()     { return (float) 3.14; }
    public Float getTestFloat()              { return (float) 6.28; }

    public short getTestShortPrimitive()     { return (short) 12; }
    public Short getTestShort()              { return (short) 21; }

    public char      getTestCharPrimitive()  { return (char) 'a'; }
    public Character getTestChar()           { return (char) 'z'; }

    public byte getTestBytePrimitive()       { return (byte) 1; }
    public Byte getTestByte()                { return (byte) 2; }

    // JMX Operations

    @Override
    public void testOperationNoParams() {
        System.out.println("JMX operation testOperationNoParams has been invoked.");
    }

    @Override
    public String testOperationPrimitive(String s, int i, boolean b, long l, double d, float f, short h, char c, byte y) {
        System.out.println("JMX operation testOperationPrimitive has been invoked.");
        return String.format("string=%s, int=%s, boolean=%s, long=%s, double=%s, float=%s, short=%s, char=%s, byte=%s",
                             s,
                             String.valueOf(i),
                             String.valueOf(b),
                             String.valueOf(l),
                             String.valueOf(d),
                             String.valueOf(f),
                             String.valueOf(h),
                             String.valueOf(c),
                             String.valueOf(y));
    }

    @Override
    public String testOperation(String s, Integer i, Boolean b, Long l, Double d, Float f, Short h, Character c, Byte y) {
        System.out.println("JMX operation testOperation has been invoked.");
        return String.format("String=%s, Int=%s, Boolean=%s, Long=%s, Double=%s, Float=%s, Short=%s, Char=%s, Byte=%s",
                             s,
                             (i == null) ? "null" : i.toString(),
                             (b == null) ? "null" : b.toString(),
                             (l == null) ? "null" : l.toString(),
                             (d == null) ? "null" : d.toString(),
                             (f == null) ? "null" : f.toString(),
                             (h == null) ? "null" : h.toString(),
                             (c == null) ? "null" : c.toString(),
                             (y == null) ? "null" : y.toString());
    }
}
