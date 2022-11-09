package hashcattypes

type HashType struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Category          string   `json:"category"`
	SlowHash          bool     `json:"slow_hash"`
	PasswordLenMin    int      `json:"password_len_min"`
	PasswordLenMax    int      `json:"password_len_max"`
	IsSalted          bool     `json:"is_salted"`
	KernelType        []string `json:"kernel_types"`
	ExampleHashFormat string   `json:"example_hash_format"`
	ExampleHash       string   `json:"example_hash"`
	ExamplePass       string   `json:"example_pass"`
	BenchmarkMask     string   `json:"benchmark_mask"`
	BenchmarkCharset1 string   `json:"benchmark_charset1"`
	AutodetectEnabled bool     `json:"autodetect_enabled"`
	SelfTestEnabled   bool     `json:"self_test_enabled"`
	PotfileEnabled    bool     `json:"potfile_enabled"`
	CustomPlugin      bool     `json:"custom_plugin"`
	PlaintextEncoding []string `json:"plaintext_encoding"`
}

type HashTypeMap map[int]*HashType
