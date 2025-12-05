package org.liquibase.lpm.model;

/**
 * Represents a single dependency entry in liquibase.json.
 * <p>
 * In the Go implementation, this is a map with a single key-value pair.
 * Here we use a record for type safety and clarity.
 *
 * @param name    the package name
 * @param version the installed version tag
 */
public record Dependency(String name, String version) {

    /**
     * Creates a dependency from a package specification string.
     * <p>
     * Supports formats: "package-name" or "package-name@version"
     *
     * @param spec the package specification
     * @return the parsed dependency
     */
    public static Dependency fromSpec(String spec) {
        if (spec == null || spec.isBlank()) {
            throw new IllegalArgumentException("Package specification cannot be empty");
        }

        int atIndex = spec.lastIndexOf('@');
        if (atIndex > 0 && atIndex < spec.length() - 1) {
            return new Dependency(spec.substring(0, atIndex), spec.substring(atIndex + 1));
        }

        return new Dependency(spec, null);
    }

    /**
     * Checks if a specific version was specified.
     *
     * @return true if a version was specified
     */
    public boolean hasVersion() {
        return version != null && !version.isBlank();
    }

    /**
     * Returns the specification string (name@version or just name).
     *
     * @return the specification string
     */
    public String toSpec() {
        return hasVersion() ? name + "@" + version : name;
    }
}
