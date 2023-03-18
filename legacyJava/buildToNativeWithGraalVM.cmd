REM beware this only works successfully from an x64 Native Tools Command Prompt
REM see https://www.graalvm.org/22.1/docs/getting-started/windows/#prerequisites-for-using-native-image-on-windows


javac BasicICSreader.java
jar -cvfe icsreader.jar BasicICSreader BasicICSreader.class
native-image -jar icsreader.jar
