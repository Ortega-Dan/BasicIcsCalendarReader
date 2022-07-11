import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.File;
import java.io.FileReader;
import java.io.FileWriter;
import java.math.BigDecimal;
import java.math.RoundingMode;
// import java.nio.charset.StandardCharsets;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.time.temporal.ChronoUnit;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Scanner;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * AdvancedReader for iCalendar 2.0 specification
 * 
 * @author Dan Ortega
 */
public class BasicICSreader {

    /**
     * Modify the first 6 lines of code in this script to set the input and output
     * files, the dates range, define if the clients-query is inclusive or
     * exclusive, and set the clients for the query.
     */
    public static void main(String[] args) throws Exception {

        if (args.length == 0) {
            System.out.println("Basic ICS Reader");
            System.out.println(
                    "Usage: [executable] [icsFilePath] [outputFileDir] [dateFrom (yyyy-mm-dd)] [dateTo (yyyy-mm-dd)] [clientsToFilter (- for no filter)] [exclusiveFilter (default: false)]");

            Scanner scan = new Scanner(System.in);
            scan.nextLine();
            scan.close();
            return;
        }

        // ********** Setting up variables

        // Input file-path in .ics format
        String inputFilePathString = args[0];

        // Output file-path in .csv format
        String outputDirPathString = args[1];

        // Dates from and to-inclusive (meaning from fromDate at 00:00 hours, until
        // toDate at 23:59 hours) in YYYY-MM-DD format
        String dateFrom = args[2];

        // Last inclusive non-ikno report ..
        // String inclusiveDateTo = "2020-06-03";
        // Last inclusive IKNO IknoPlus report !!!!!!!
        String inclusiveDateTo = args[3];

        // ClientsToQuery: Case insensitive String of clients separated by pipes "|"
        // OR AN EMPTY STRING IF LOOKING FOR ALL ACTIVITIES in a time range or span.
        // For example to look for MTI and DATAFILE activities, set it to "mti|DataFile"
        // String clientsToQuery = "ikno|kno|noclient|rv";
        String clientsToQuery = args[4].equals("-") ? "" : args[4];

        // Set to true if you want to print all clients but the ones in the query ...
        // ... set to false if you want to print the clients in the query
        boolean excludeClientsInQuery = false;

        try {
            excludeClientsInQuery = Boolean.parseBoolean(args[5]);
        } catch (Exception e) {
        }

        //
        // ********** end of variables setting up
        //

        //
        // Logic goes from here ...
        Pattern clientFinderPattern = Pattern.compile("([^:]+):.+");

        FileReader fr = new FileReader(inputFilePathString);
        BufferedReader br = new BufferedReader(fr);

        String startTime, endTime, summary, client;
        startTime = endTime = summary = client = "";

        // reference
        // https://docs.oracle.com/en/java/javase/11/docs/api/java.base/java/time/format/DateTimeFormatter.html
        DateTimeFormatter fromZuluFormatter = DateTimeFormatter.ofPattern("yyyyMMdd'T'HHmmssX");
        DateTimeFormatter dateTimeFormatter = DateTimeFormatter.ofPattern("yyyyMMdd'T'HHmmss");

        HashSet<String> allClientsInRange = new HashSet<>();
        HashSet<String> clientsMatchedInQuery = new HashSet<>();

        // Creating missing directories if needed
        File outFile = new File(outputDirPathString + File.separator + inputFilePathString.replaceAll("(.*/|@.+)", "")
                + "_" + dateFrom + "_" + inclusiveDateTo + ".tsv");
        outFile.getParentFile().mkdirs();

        // Writing it for Windows csv:
        // FileWriter fw = new FileWriter(outFile, Charset.forName("Cp1252"));
        // Do this if you want to write it in UTF-8 instead
        // FileWriter fw = new FileWriter(outFile, StandardCharsets.UTF_8); // This
        // constructor works from Java 11
        FileWriter fw = new FileWriter(outFile);

        BufferedWriter writer = new BufferedWriter(fw);

        // Getting from and to RANGE dates ready
        ZonedDateTime startRangeDT = LocalDateTime.parse(dateFrom.replace("-", "") + "T000000", dateTimeFormatter)
                .atZone(ZoneId.systemDefault());

        ZonedDateTime endRangeDT = LocalDateTime.parse(inclusiveDateTo.replace("-", "") + "T235959", dateTimeFormatter)
                .atZone(ZoneId.systemDefault());

        BigDecimal totalHours = new BigDecimal("0");

        writer.write("Fecha\tHoras\tCliente\tDescripción");

        String line = "";
        while ((line = br.readLine()) != null) {

            // Checking iCalendar specification version to ensure logic processing
            if (line.trim().matches("VERSION:.+")) {

                final String supportedVersionString = "2.0";

                BigDecimal version = new BigDecimal(line.replace("VERSION:", ""));
                System.out.println("\niCalendar specification: " + version + "\n");

                if (version.compareTo(new BigDecimal(supportedVersionString)) != 0) {

                    System.out.println("This program was written for iCalendar " + supportedVersionString);
                    System.out.println("Please use an .ics file with " + supportedVersionString
                            + " version or contact the developer of this program.");

                    writer.close();
                    br.close();
                    return;
                }

            }

            // Setting startTime, endTime, client, and summary variables
            if (line.trim().matches("BEGIN:VEVENT")) {

                StringBuilder allEventLines = new StringBuilder();

                while ((line = br.readLine()) != null && !line.trim().equalsIgnoreCase("END:VEVENT")) {

                    allEventLines.append(line + "\n");

                    if (line.trim().matches("DTSTART.*:.+")) {
                        line = line.replaceFirst("DTSTART.*:", "");
                        startTime = line;
                    } else if (line.trim().matches("DTEND.*:.+")) {
                        line = line.replaceFirst("DTEND.*:", "");
                        endTime = line;
                    } else if (line.trim().matches("SUMMARY:.+")) {
                        summary = line;
                        while ((line = br.readLine()) != null && line.matches(" .+")) {
                            summary += line.substring(1);
                        }
                        summary = summary.replaceFirst("SUMMARY:", "").replaceAll("\\\\,", ",");

                        // Getting the client out of the summary
                        Matcher clientMatcher = clientFinderPattern.matcher(summary);
                        if (clientMatcher.find()) {
                            client = clientMatcher.group(1).trim();
                        } else {
                            client = "NoClient";
                        }
                    }
                }

                // Avoiding lunch or ignore hours
                if (summary.trim().matches("Lunch") || summary.contains("cal.ignore")) {
                    continue;
                }

                // Full days activities are only displayed and ignored after being filtered
                // see later usage of this variable
                boolean isFullDayEvent = false;
                if (startTime.length() == 8) {
                    startTime += "T000000";
                    endTime += "T000000";
                    isFullDayEvent = true;
                }

                ZonedDateTime startDT = null;
                if (startTime.contains("Z")) {
                    startDT = ZonedDateTime.parse(startTime, fromZuluFormatter)
                            .withZoneSameInstant(ZoneId.systemDefault());
                } else {
                    startDT = LocalDateTime.parse(startTime, dateTimeFormatter).atZone(ZoneId.systemDefault());
                }

                ZonedDateTime endDT = null;
                if (endTime.contains("Z")) {
                    endDT = ZonedDateTime.parse(endTime, fromZuluFormatter).withZoneSameInstant(ZoneId.systemDefault());
                } else {
                    endDT = LocalDateTime.parse(endTime, dateTimeFormatter).atZone(ZoneId.systemDefault());
                }

                // Checking if date is in the range
                if ((startDT.isEqual(startRangeDT) || startDT.isAfter(startRangeDT))
                        && (startDT.isBefore(endRangeDT) || startDT.isEqual(endRangeDT))) {

                    // Uncomment next line for stdout insight
                    // if (summary.contains("ChangeThisSomething")) {
                    // System.out.println(allEventLines.toString()+"\n************************************************\n");
                    // }

                    // Full days activities are only displayed and ignored after being filtered
                    if (isFullDayEvent) {
                        System.out.println("FullDaysActivity: [" + client + "] " + startTime + "-" + endTime + " ["
                                + summary + "]");
                        continue;
                    }

                    // Adding all clients in the dates-range for later information to the user
                    allClientsInRange.add(client);

                    // CLIENT FILTERING (case insensitive matching)
                    boolean isClientMatching = client.matches("(?i)(" + clientsToQuery + ").*");
                    // if client is match and query is inclusive **** OR **** if client doesn't
                    // match and query is exclusive
                    if ((isClientMatching && !excludeClientsInQuery) || (!isClientMatching && excludeClientsInQuery)) {

                        // Recording filtered clients to show them later
                        clientsMatchedInQuery.add(client);

                        BigDecimal minutesLenght = new BigDecimal(startDT.until(endDT, ChronoUnit.MINUTES) + ".0");
                        // zeroing non-end-time events
                        if (endDT.isBefore(startDT)) {
                            minutesLenght = new BigDecimal(0);
                        }

                        BigDecimal hoursLength = minutesLenght.divide(new BigDecimal("60.0"), 2, RoundingMode.HALF_UP);

                        totalHours = totalHours.add(hoursLength);

                        // Clean client from summary before writing
                        summary = summary.replaceFirst(client + ":", "").trim();
                        writer.write(
                                "\n" + startDT.toLocalDate() + "\t" + hoursLength + "\t" + client + "\t" + summary);
                    }
                }
            }

        }

        // Writing last line with total hours
        writer.write("\n\t" + totalHours + "\tTOTAL\t");

        writer.close();
        br.close();

        System.out.println("\n *** Clients included in report: " + clientsMatchedInQuery.size());
        if (clientsMatchedInQuery.size() != allClientsInRange.size()) {
            printSortedClients(clientsMatchedInQuery);
        } else {
            System.out.println("\nAll found clients were included in report.");
        }

        System.out.println("\n *** Total Clients found in required dates range: " + allClientsInRange.size() + "\n");
        printSortedClients(allClientsInRange);

    }

    private static void printSortedClients(HashSet<String> clients) {

        List<String> clientsList = new ArrayList<>(clients);

        Collections.sort(clientsList);

        clientsList.forEach(System.out::println);

    }

}
