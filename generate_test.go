package facto

import (
	_ "embed"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestCammelCase(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"aaa", "Aaa"},
		{"the cow", "TheCow"},
		{"the_cow", "TheCow"},
		{"the-cow", "TheCow"},
	}

	for _, v := range cases {
		t.Run(v.in, func(t *testing.T) {
			out := camelCase(v.in)
			if out != v.want {
				t.Errorf("camelCase(%q) = %q, want %q", v.in, out, v.want)
			}
		})
	}
}

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"camelCase", "camel_case"},
		{"PascalCase", "pascal_case"},
		{"snake_case", "snake_case"},
		{"Pascal_Snake", "pascal_snake"},
		{"SCREAMING_SNAKE", "screaming_snake"},
		{"kebab-case", "kebab_case"},
		{"Pascal-Kebab", "pascal_kebab"},
		{"SCREAMING-KEBAB", "screaming_kebab"},
		{"A", "a"},
		{"AA", "aa"},
		{"AAA", "aaa"},
		{"AAAA", "aaaa"},
		{"AaAa", "aa_aa"},
		{"BatteryLifeValue", "battery_life_value"},
		{"Id0Value", "id0_value"},
		{"ID0Value", "id0_value"},
	}

	for _, v := range tests {
		result := snakeCase(v.input)

		if result != v.expected {
			t.Errorf("snakeCase(%q) = %q, want %q", v.input, result, v.expected)
		}
	}
}

func TestGenerate(t *testing.T) {
	dir := os.TempDir()
	os.Chdir(dir)

	err := Generate(dir, []string{"generate", "user"})
	if err != nil {
		t.Fatalf("Err should be nil, got %v", err)
	}

	data, err := ioutil.ReadFile(filepath.Join(dir, "factories", "user.go"))
	if err != nil {
		t.Fatalf("could not read factory file, got %v", err)
	}

	if !strings.Contains(string(data), `package factories`) {
		t.Errorf("factory file should contain package name, got: \n%s", data)
	}

	if !strings.Contains(string(data), `func UserFactory(h facto.Helper)`) {
		t.Errorf("factory file should contain func definition, got: \n%s", data)
	}
}

func TestGenerateDirExists(t *testing.T) {
	d := os.TempDir()
	os.Chdir(d)

	deleteAll := func() {
		dir, _ := ioutil.ReadDir(d)
		for _, f := range dir {
			os.RemoveAll(path.Join([]string{d, f.Name()}...))
		}
	}

	t.Run("HappyPath", func(t *testing.T) {
		t.Cleanup(deleteAll)

		err := Generate(d, []string{"generate", "user"})
		if err != nil {
			t.Fatalf("Err should be nil, got %v", err)
		}

		data, err := ioutil.ReadFile(filepath.Join(d, "factories", "user.go"))
		if err != nil {
			t.Fatalf("could not read factory file, got %v", err)
		}

		if !strings.Contains(string(data), `package factories`) {
			t.Errorf("factory file should contain package name, got: \n%s", data)
		}

		if !strings.Contains(string(data), `func UserFactory(h facto.Helper)`) {
			t.Errorf("factory file should contain func definition, got: \n%s", data)
		}
	})

	t.Run("FolderExists", func(t *testing.T) {
		t.Cleanup(deleteAll)

		err := os.MkdirAll(filepath.Join(d, "factories"), 0777)
		if err != nil {
			t.Fatal("could not create folder")
		}

		err = Generate(d, []string{"generate", "user"})
		if err != nil {
			t.Fatalf("Err should be nil, got %v", err)
		}

		data, err := ioutil.ReadFile(filepath.Join(d, "factories", "user.go"))
		if err != nil {
			t.Fatalf("could not read factory file, got %v", err)
		}

		if !strings.Contains(string(data), `func UserFactory(h facto.Helper)`) {
			t.Errorf("factory file should contain func definition, got: \n%s", data)
		}
	})

	t.Run("FolderAndFileExists", func(t *testing.T) {
		t.Cleanup(deleteAll)

		err := os.MkdirAll(filepath.Join(d, "factories"), 0777)
		if err != nil {
			t.Fatal("could not create folder")
		}

		err = ioutil.WriteFile(filepath.Join(d, "factories", "user.go"), []byte(`package factories`), 0777)
		if err != nil {
			t.Fatalf("could not write the file")
		}

		err = Generate(d, []string{"generate", "user"})
		if err != nil {
			t.Fatalf("Err should be nil, got %v", err)
		}

		data, err := ioutil.ReadFile(filepath.Join(d, "factories", "user.go"))
		if err != nil {
			t.Fatalf("could not read factory file, got %v", err)
		}

		if !strings.Contains(string(data), `func UserFactory(h facto.Helper)`) {
			t.Errorf("factory file should contain func definition, got: \n%s", data)
		}
	})

}
