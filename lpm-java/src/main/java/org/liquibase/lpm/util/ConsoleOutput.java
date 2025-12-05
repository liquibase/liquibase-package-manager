package org.liquibase.lpm.util;

import org.liquibase.lpm.model.Package;
import org.liquibase.lpm.model.PackageVersion;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

/**
 * Utility class for formatting CLI output.
 * <p>
 * Provides tree-style formatting consistent with the Go implementation.
 */
public class ConsoleOutput {

    /**
     * Tree branch prefix for non-last items.
     */
    public static final String BRANCH = "├──";

    /**
     * Tree branch prefix for the last item.
     */
    public static final String LAST_BRANCH = "└──";

    /**
     * Formats a list of packages for display.
     *
     * @param packages           the packages to display
     * @param classpathFilenames list of filenames in the classpath (for version display)
     * @return formatted lines for output
     */
    public static List<String> formatPackages(List<Package> packages, List<String> classpathFilenames) {
        List<String> lines = new ArrayList<>();

        // Header
        lines.add(String.format("%-4s %-38s %s", "   ", "Package", "Category"));

        // Packages
        for (int i = 0; i < packages.size(); i++) {
            Package pkg = packages.get(i);
            String prefix = (i == packages.size() - 1) ? LAST_BRANCH : BRANCH;

            // Get installed version suffix
            String versionSuffix = "";
            Optional<PackageVersion> installed = pkg.getInstalledVersion(classpathFilenames);
            if (installed.isPresent()) {
                versionSuffix = "@" + installed.get().tag();
            }

            lines.add(String.format("%-4s %-38s %s", prefix, pkg.name() + versionSuffix, pkg.category()));
        }

        return lines;
    }

    /**
     * Formats outdated packages for upgrade display.
     *
     * @param packages           the outdated packages
     * @param classpathFilenames list of filenames in the classpath
     * @param latestVersions     function to get latest version for each package
     * @return formatted lines for output
     */
    public static List<String> formatOutdatedPackages(
            List<Package> packages,
            List<String> classpathFilenames,
            java.util.function.Function<Package, Optional<PackageVersion>> latestVersionGetter) {

        List<String> lines = new ArrayList<>();

        // Header
        lines.add(String.format("%-4s %-38s %-12s %s", "   ", "Package", "Installed", "Available"));

        // Packages
        for (int i = 0; i < packages.size(); i++) {
            Package pkg = packages.get(i);
            String prefix = (i == packages.size() - 1) ? LAST_BRANCH : BRANCH;

            Optional<PackageVersion> installed = pkg.getInstalledVersion(classpathFilenames);
            Optional<PackageVersion> latest = latestVersionGetter.apply(pkg);

            String installedVersion = installed.map(PackageVersion::tag).orElse("-");
            String latestVersion = latest.map(PackageVersion::tag).orElse("-");

            String packageDisplay = pkg.name();
            if (installed.isPresent()) {
                packageDisplay += "@" + installed.get().tag();
            }

            lines.add(String.format("%-4s %-38s %-12s %s",
                    prefix, packageDisplay, installedVersion, latestVersion));
        }

        return lines;
    }

    /**
     * Formats a list of versions for dedupe display.
     *
     * @param packageName the package name
     * @param versions    the installed versions
     * @return formatted lines for output
     */
    public static List<String> formatVersionsForDedupe(String packageName, List<PackageVersion> versions) {
        List<String> lines = new ArrayList<>();

        lines.add("Package");
        for (int i = 0; i < versions.size(); i++) {
            PackageVersion v = versions.get(i);
            String prefix = (i == versions.size() - 1) ? LAST_BRANCH : BRANCH;
            lines.add(String.format("%-4s %s@%s", prefix, packageName, v.tag()));
        }

        return lines;
    }

    /**
     * Prints lines to standard output.
     *
     * @param lines the lines to print
     */
    public static void printLines(List<String> lines) {
        for (String line : lines) {
            System.out.println(line);
        }
    }

    /**
     * Prints a success message for package installation.
     *
     * @param filename the installed filename
     */
    public static void printInstallSuccess(String filename) {
        System.out.println(filename + " successfully installed in classpath.");
    }

    /**
     * Prints a success message for package removal.
     *
     * @param filename the removed filename
     */
    public static void printRemoveSuccess(String filename) {
        System.out.println(filename + " successfully removed from classpath.");
    }

    /**
     * Prints the classpath location.
     *
     * @param classpathDir the classpath directory
     */
    public static void printClasspath(String classpathDir) {
        System.out.println("Classpath: " + classpathDir);
    }

    /**
     * Prints a message about the number of outdated packages.
     *
     * @param count the number of outdated packages
     */
    public static void printOutdatedCount(int count) {
        if (count == 0) {
            System.out.println("All packages are up to date.");
        } else if (count == 1) {
            System.out.println("1 package can be upgraded.");
        } else {
            System.out.println(count + " packages can be upgraded.");
        }
    }

    /**
     * Prints JAVA_OPTS instructions for older Liquibase versions.
     *
     * @param classpathDir the local classpath directory
     */
    public static void printJavaOptsInstructions(String classpathDir) {
        System.out.println();
        System.out.println("For Liquibase versions prior to 4.6.2, you need to add the following to your JAVA_OPTS:");
        System.out.println("  export JAVA_OPTS=\"-Dliquibase.classpath=" + classpathDir + "\"");
    }
}
