import javax.management.*;

@MXBean
public interface SimpleMXBean {

    // JMX Attributes

    String getTestString();

    int     getTestIntegerPrimitive();
    Integer getTestInteger();

    boolean getTestBooleanPrimitive();
    Boolean getTestBoolean();

    long getTestLongPrimitive();
    Long getTestLong();

    double getTestDoublePrimitive();
    Double getTestDouble();

    float getTestFloatPrimitive();
    Float getTestFloat();

    short getTestShortPrimitive();
    Short getTestShort();

    char      getTestCharPrimitive();
    Character getTestChar();

    byte getTestBytePrimitive();
    Byte getTestByte();

    // JMX Operations

    void testOperationNoParams();
    String testOperationPrimitive(String s, int i, boolean b, long l, double d, float f, short h, char c, byte y);
    String testOperation(String s, Integer i, Boolean b, Long l, Double d, Float f, Short h, Character c, Byte y);

}
