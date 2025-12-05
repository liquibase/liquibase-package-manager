package org.liquibase.lpm.model;

import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class DependencyTest {

    @Test
    void fromSpec_parsesNameOnly() {
        Dependency dep = Dependency.fromSpec("package-name");

        assertThat(dep.name()).isEqualTo("package-name");
        assertThat(dep.version()).isNull();
    }

    @Test
    void fromSpec_parsesNameAndVersion() {
        Dependency dep = Dependency.fromSpec("package-name@1.0.0");

        assertThat(dep.name()).isEqualTo("package-name");
        assertThat(dep.version()).isEqualTo("1.0.0");
    }

    @Test
    void fromSpec_handlesVersionWithAtSign() {
        // Version itself contains @ (unlikely but should handle gracefully)
        Dependency dep = Dependency.fromSpec("package@1.0.0@beta");

        assertThat(dep.name()).isEqualTo("package@1.0.0");
        assertThat(dep.version()).isEqualTo("beta");
    }

    @Test
    void fromSpec_throwsForNull() {
        assertThatThrownBy(() -> Dependency.fromSpec(null))
                .isInstanceOf(IllegalArgumentException.class);
    }

    @Test
    void fromSpec_throwsForEmpty() {
        assertThatThrownBy(() -> Dependency.fromSpec(""))
                .isInstanceOf(IllegalArgumentException.class);
    }

    @Test
    void fromSpec_throwsForBlank() {
        assertThatThrownBy(() -> Dependency.fromSpec("   "))
                .isInstanceOf(IllegalArgumentException.class);
    }

    @Test
    void hasVersion_returnsTrueWhenVersionSet() {
        Dependency dep = new Dependency("pkg", "1.0.0");

        assertThat(dep.hasVersion()).isTrue();
    }

    @Test
    void hasVersion_returnsFalseWhenVersionNull() {
        Dependency dep = new Dependency("pkg", null);

        assertThat(dep.hasVersion()).isFalse();
    }

    @Test
    void hasVersion_returnsFalseWhenVersionBlank() {
        Dependency dep = new Dependency("pkg", "  ");

        assertThat(dep.hasVersion()).isFalse();
    }

    @Test
    void toSpec_returnsNameWithVersion() {
        Dependency dep = new Dependency("package", "2.0.0");

        assertThat(dep.toSpec()).isEqualTo("package@2.0.0");
    }

    @Test
    void toSpec_returnsNameOnlyWhenNoVersion() {
        Dependency dep = new Dependency("package", null);

        assertThat(dep.toSpec()).isEqualTo("package");
    }

    @Test
    void roundTrip_preservesData() {
        Dependency original = new Dependency("my-package", "3.2.1");
        Dependency parsed = Dependency.fromSpec(original.toSpec());

        assertThat(parsed).isEqualTo(original);
    }
}
