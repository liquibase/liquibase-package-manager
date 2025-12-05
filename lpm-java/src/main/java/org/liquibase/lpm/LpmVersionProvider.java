package org.liquibase.lpm;

import picocli.CommandLine.IVersionProvider;

import java.io.InputStream;
import java.nio.charset.StandardCharsets;

/**
 * Provides version information for the CLI from embedded VERSION resource.
 */
public class LpmVersionProvider implements IVersionProvider {

    @Override
    public String[] getVersion() throws Exception {
        try (InputStream is = getClass().getResourceAsStream("/VERSION")) {
            if (is == null) {
                return new String[]{"unknown"};
            }
            String version = new String(is.readAllBytes(), StandardCharsets.UTF_8).trim();
            return new String[]{"lpm (Liquibase Package Manager) version " + version};
        }
    }
}
