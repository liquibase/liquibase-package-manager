package dependencies

import "testing"

func TestDependency_GetName(t *testing.T) {
	tests := []struct {
		name string
		d    Dependency
		want string
	}{
		{
			name: "Can Get Name",
			d: Dependency{"package":"tag"},
			want: "package",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDependency_GetVersion(t *testing.T) {
	tests := []struct {
		name string
		d    Dependency
		want string
	}{
		{
			name: "Can Get Version",
			d: Dependency{"package":"tag"},
			want: "tag",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetVersion(); got != tt.want {
				t.Errorf("GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
