package org.liquibase.lpm.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.liquibase.lpm.exception.ManifestException;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.stream.Collectors;

/**
 * Manages the liquibase.json dependency manifest file.
 * <p>
 * This file tracks which packages have been installed in the local project,
 * allowing them to be reinstalled via the "install" command.
 */
public class DependencyManifest {

    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();
    private static final String DEFAULT_FILENAME = "liquibase.json";

    private final Path manifestPath;
    private List<Dependency> dependencies;

    /**
     * Creates a new dependency manifest for the current working directory.
     */
    public DependencyManifest() {
        this(Path.of(System.getProperty("user.dir"), DEFAULT_FILENAME));
    }

    /**
     * Creates a new dependency manifest at the specified path.
     *
     * @param manifestPath the path to the manifest file
     */
    public DependencyManifest(Path manifestPath) {
        this.manifestPath = manifestPath;
        this.dependencies = new ArrayList<>();
    }

    /**
     * Gets the path to the manifest file.
     *
     * @return the manifest path
     */
    public Path getManifestPath() {
        return manifestPath;
    }

    /**
     * Checks if the manifest file exists.
     *
     * @return true if the file exists
     */
    public boolean fileExists() {
        return Files.exists(manifestPath);
    }

    /**
     * Reads the manifest from the file.
     * <p>
     * If the file doesn't exist, the dependencies list is cleared.
     *
     * @return this manifest (for chaining)
     */
    public DependencyManifest read() {
        if (!fileExists()) {
            dependencies = new ArrayList<>();
            return this;
        }

        try {
            String content = Files.readString(manifestPath);
            ManifestJson json = OBJECT_MAPPER.readValue(content, ManifestJson.class);

            dependencies = new ArrayList<>();
            if (json.dependencies != null) {
                for (Map<String, String> dep : json.dependencies) {
                    // Each map has one entry: {name: version}
                    for (Map.Entry<String, String> entry : dep.entrySet()) {
                        dependencies.add(new Dependency(entry.getKey(), entry.getValue()));
                    }
                }
            }
        } catch (IOException e) {
            throw new ManifestException("Failed to read " + manifestPath, e);
        }

        return this;
    }

    /**
     * Writes the manifest to the file.
     *
     * @return this manifest (for chaining)
     */
    public DependencyManifest write() {
        try {
            // Convert to the JSON structure: {"dependencies": [{"name": "version"}, ...]}
            List<Map<String, String>> depList = dependencies.stream()
                    .map(d -> Map.of(d.name(), d.version()))
                    .collect(Collectors.toList());

            ManifestJson json = new ManifestJson(depList);
            String content = OBJECT_MAPPER.writerWithDefaultPrettyPrinter().writeValueAsString(json);
            Files.writeString(manifestPath, content);
        } catch (IOException e) {
            throw new ManifestException("Failed to write " + manifestPath, e);
        }

        return this;
    }

    /**
     * Creates the manifest file if it doesn't exist.
     *
     * @return this manifest (for chaining)
     */
    public DependencyManifest createIfNotExists() {
        if (!fileExists()) {
            dependencies = new ArrayList<>();
            write();
        }
        return this;
    }

    /**
     * Gets all dependencies in the manifest.
     *
     * @return an unmodifiable list of dependencies
     */
    public List<Dependency> getDependencies() {
        return Collections.unmodifiableList(dependencies);
    }

    /**
     * Gets a dependency by name.
     *
     * @param name the package name
     * @return the dependency, or empty if not found
     */
    public Optional<Dependency> getByName(String name) {
        return dependencies.stream()
                .filter(d -> name.equals(d.name()))
                .findFirst();
    }

    /**
     * Adds or updates a dependency.
     *
     * @param dependency the dependency to add/update
     * @return this manifest (for chaining)
     */
    public DependencyManifest add(Dependency dependency) {
        // Remove existing entry for this package if present
        dependencies.removeIf(d -> dependency.name().equals(d.name()));
        dependencies.add(dependency);
        return this;
    }

    /**
     * Removes a dependency by name.
     *
     * @param name the package name to remove
     * @return this manifest (for chaining)
     */
    public DependencyManifest remove(String name) {
        dependencies.removeIf(d -> name.equals(d.name()));
        return this;
    }

    /**
     * Checks if a package is in the manifest.
     *
     * @param name the package name
     * @return true if the package is in the manifest
     */
    public boolean contains(String name) {
        return dependencies.stream().anyMatch(d -> name.equals(d.name()));
    }

    /**
     * Checks if the manifest is empty.
     *
     * @return true if there are no dependencies
     */
    public boolean isEmpty() {
        return dependencies.isEmpty();
    }

    /**
     * Gets the number of dependencies.
     *
     * @return the dependency count
     */
    public int size() {
        return dependencies.size();
    }

    /**
     * Internal JSON structure for serialization.
     * Matches the Go format: {"dependencies": [{"name": "version"}, ...]}
     */
    private record ManifestJson(
            @JsonProperty("dependencies") List<Map<String, String>> dependencies
    ) {
    }
}
