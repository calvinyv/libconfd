// Copyright confd. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE-confd file.

package libconfd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	pathpkg "path"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type TemplateFunc struct {
	FuncMap       map[string]interface{}
	Store         *KVStore
	PGPPrivateKey []byte
}

func NewTemplateFunc(store *KVStore, pgpPrivateKey []byte) TemplateFunc {
	return TemplateFunc{
		Store:         store,
		PGPPrivateKey: pgpPrivateKey,
	}
}

func NewTemplateFuncMap(store *KVStore, pgpPrivateKey []byte) (TemplateFunc, template.FuncMap) {
	p := TemplateFunc{
		Store:         store,
		PGPPrivateKey: pgpPrivateKey,
	}

	p.FuncMap = template.FuncMap{
		// KVStore
		"exists": p.Exists,
		"ls":     p.Ls,
		"lsdir":  p.Lsdir,
		"get":    p.Get,
		"gets":   p.Gets,
		"getv":   p.Getv,
		"getvs":  p.Getvs,

		// more tmpl func
		"base":           p.Base,
		"split":          p.Split,
		"json":           p.Json,
		"jsonArray":      p.JsonArray,
		"dir":            p.Dir,
		"map":            p.Map,
		"getenv":         p.Getenv,
		"join":           p.Join,
		"datetime":       p.Datetime,
		"toUpper":        p.ToUpper,
		"toLower":        p.ToLower,
		"contains":       p.Contains,
		"replace":        p.Replace,
		"trimSuffix":     p.TrimSuffix,
		"lookupIP":       p.LookupIP,
		"lookupSRV":      p.LookupSRV,
		"fileExists":     p.FileExists,
		"base64Encode":   p.Base64Encode,
		"base64Decode":   p.Base64Decode,
		"parseBool":      p.ParseBool,
		"reverse":        p.Reverse,
		"sortByLength":   p.SortByLength,
		"sortKVByLength": p.SortKVByLength,
		"add":            p.Add,
		"sub":            p.Sub,
		"div":            p.Div,
		"mod":            p.Mod,
		"mul":            p.Mul,
		"seq":            p.Seq,
		"atoi":           p.Atoi,
	}

	// crypt func
	if len(pgpPrivateKey) > 0 {
		p.FuncMap["cget"] = p.Cget
		p.FuncMap["cgets"] = p.Cgets
		p.FuncMap["cgetv"] = p.Cgetv
		p.FuncMap["cgetvs"] = p.Cgetvs
	}

	return p, p.FuncMap
}

// ----------------------------------------------------------------------------
// KVStore
// ----------------------------------------------------------------------------

func (p TemplateFunc) Exists(key string) bool {
	return p.Store.Exists(key)
}

func (p TemplateFunc) Ls(filepath string) []string {
	return p.Store.List(filepath)
}

func (p TemplateFunc) Lsdir(filepath string) []string {
	return p.Store.ListDir(filepath)
}

func (p TemplateFunc) Get(key string) (KVPair, error) {
	return p.Store.Get(key)
}

func (p TemplateFunc) Gets(pattern string) ([]KVPair, error) {
	return p.Store.GetAll(pattern)
}

func (p TemplateFunc) Getv(key string, v ...string) (string, error) {
	return p.Store.GetValue(key, v...)
}

func (p TemplateFunc) Getvs(pattern string) ([]string, error) {
	return p.Store.GetAllValues(pattern)
}

// ----------------------------------------------------------------------------
// Crypt func
// ----------------------------------------------------------------------------

func (p TemplateFunc) Cget(key string) (KVPair, error) {
	kv, err := p.FuncMap["get"].(func(string) (KVPair, error))(key)
	if err == nil {
		var b []byte
		b, err = secconfDecode([]byte(kv.Value), bytes.NewBuffer(p.PGPPrivateKey))
		if err == nil {
			kv.Value = string(b)
		}
	}
	return kv, err
}

func (p TemplateFunc) Cgets(pattern string) ([]KVPair, error) {
	kvs, err := p.FuncMap["gets"].(func(string) ([]KVPair, error))(pattern)
	if err == nil {
		for i := range kvs {
			b, err := secconfDecode([]byte(kvs[i].Value), bytes.NewBuffer(p.PGPPrivateKey))
			if err != nil {
				return []KVPair(nil), err
			}
			kvs[i].Value = string(b)
		}
	}
	return kvs, err
}

func (p TemplateFunc) Cgetv(key string) (string, error) {
	v, err := p.FuncMap["getv"].(func(string, ...string) (string, error))(key)
	if err == nil {
		var b []byte
		b, err = secconfDecode([]byte(v), bytes.NewBuffer(p.PGPPrivateKey))
		if err == nil {
			return string(b), nil
		}
	}
	return v, err
}

func (p TemplateFunc) Cgetvs(pattern string) ([]string, error) {
	vs, err := p.FuncMap["getvs"].(func(string) ([]string, error))(pattern)
	if err == nil {
		for i := range vs {
			b, err := secconfDecode([]byte(vs[i]), bytes.NewBuffer(p.PGPPrivateKey))
			if err != nil {
				return []string(nil), err
			}
			vs[i] = string(b)
		}
	}
	return vs, err
}

// ----------------------------------------------------------------------------
// util func
// ----------------------------------------------------------------------------

func (_ TemplateFunc) Base(path string) string {
	return pathpkg.Base(path)
}

func (_ TemplateFunc) Split(s, sep string) []string {
	return strings.Split(s, sep)
}

func (_ TemplateFunc) Json(data string) (map[string]interface{}, error) {
	var ret map[string]interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func (_ TemplateFunc) JsonArray(data string) ([]interface{}, error) {
	var ret []interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func (_ TemplateFunc) Dir(path string) string {
	return pathpkg.Dir(path)
}

// Map creates a key-value map of string -> interface{}
// The i'th is the key and the i+1 is the value
func (_ TemplateFunc) Map(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid map call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("map keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// getenv retrieves the value of the environment variable named by the key.
// It returns the value, which will the default value if the variable is not present.
// If no default value was given - returns "".
func (_ TemplateFunc) Getenv(key string, defaultValue ...string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (_ TemplateFunc) Join(a []string, sep string) string {
	return strings.Join(a, sep)
}

func (_ TemplateFunc) Datetime() time.Time {
	return time.Now()
}

func (_ TemplateFunc) ToUpper(s string) string {
	return strings.ToUpper(s)
}

func (_ TemplateFunc) ToLower(s string) string {
	return strings.ToLower(s)
}

func (_ TemplateFunc) Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func (_ TemplateFunc) Replace(s, old, new string, n int) string {
	return strings.Replace(s, old, new, n)
}

func (_ TemplateFunc) TrimSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

func (_ TemplateFunc) LookupIP(data string) []string {
	ips, err := net.LookupIP(data)
	if err != nil {
		return nil
	}
	// "Cast" IPs into strings and sort the array
	ipStrings := make([]string, len(ips))

	for i, ip := range ips {
		ipStrings[i] = ip.String()
	}
	sort.Strings(ipStrings)
	return ipStrings
}

func (_ TemplateFunc) LookupSRV(service, proto, name string) []*net.SRV {
	_, s, err := net.LookupSRV(service, proto, name)
	if err != nil {
		return []*net.SRV{}
	}

	sort.Slice(s, func(i, j int) bool {
		str1 := fmt.Sprintf("%s%d%d%d", s[i].Target, s[i].Port, s[i].Priority, s[i].Weight)
		str2 := fmt.Sprintf("%s%d%d%d", s[j].Target, s[j].Port, s[j].Priority, s[j].Weight)
		return str1 < str2
	})
	return s
}

func (_ TemplateFunc) FileExists(filepath string) bool {
	return utilFileExist(filepath)
}

func (_ TemplateFunc) Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func (_ TemplateFunc) Base64Decode(data string) (string, error) {
	s, err := base64.StdEncoding.DecodeString(data)
	return string(s), err
}

func (_ TemplateFunc) ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// reverse returns the array in reversed order
// works with []string and []KVPair
func (_ TemplateFunc) Reverse(values interface{}) interface{} {
	switch values.(type) {
	case []string:
		v := values.([]string)
		for left, right := 0, len(v)-1; left < right; left, right = left+1, right-1 {
			v[left], v[right] = v[right], v[left]
		}
	case []KVPair:
		v := values.([]KVPair)
		for left, right := 0, len(v)-1; left < right; left, right = left+1, right-1 {
			v[left], v[right] = v[right], v[left]
		}
	}
	return values
}

func (_ TemplateFunc) SortKVByLength(values []KVPair) []KVPair {
	sort.Slice(values, func(i, j int) bool {
		return len(values[i].Key) < len(values[j].Key)
	})
	return values
}

func (_ TemplateFunc) SortByLength(values []string) []string {
	sort.Slice(values, func(i, j int) bool {
		return len(values[i]) < len(values[j])
	})
	return values
}

func (_ TemplateFunc) Add(a, b int) int {
	return a + b
}
func (_ TemplateFunc) Sub(a, b int) int {
	return a - b
}
func (_ TemplateFunc) Div(a, b int) int {
	return a / b
}
func (_ TemplateFunc) Mod(a, b int) int {
	return a % b
}
func (_ TemplateFunc) Mul(a, b int) int {
	return a * b
}

// seq creates a sequence of integers. It's named and used as GNU's seq.
// seq takes the first and the last element as arguments. So Seq(3, 5) will generate [3,4,5]
func (_ TemplateFunc) Seq(first, last int) []int {
	var arr []int
	for i := first; i <= last; i++ {
		arr = append(arr, i)
	}
	return arr
}

func (_ TemplateFunc) Atoi(s string) (int, error) {
	return strconv.Atoi(s)
}

// ----------------------------------------------------------------------------
// END
// ----------------------------------------------------------------------------
