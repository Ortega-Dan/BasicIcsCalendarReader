javac BasicICSreader.java
jar -cvfe icsreader.jar BasicICSreader BasicICSreader.class
native-image -jar icsreader.jar
