package DouYu

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dop251/goja"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func (l *Link) getRateStream() (gjson.Result, error) {
	re := regexp.MustCompile(`(function ub98484234.*)\s(var.*)`)
	match := re.FindStringSubmatch(l.res)
	if len(match) < 2 {
		return gjson.Result{}, fmt.Errorf("could not find ub98484234 function")
	}
	jsFunc := match[0]

	// Step 2: Modify the JavaScript function to remove eval statement
	replaceRe := regexp.MustCompile(`eval.*;}`)
	jsFunc = replaceRe.ReplaceAllString(jsFunc, "strc;}")

	// Step 3: Compile and execute the modified JavaScript function
	js := fmt.Sprintf("%s\nub98484234()", jsFunc)

	vm := goja.New()
	vmResult, err := vm.RunString(js)
	if err != nil {
		return gjson.Result{}, err
	}
	result := vmResult.String()
	re = regexp.MustCompile(`v=(\d+)`)
	match = re.FindStringSubmatch(result)
	if len(match) < 2 {
		return gjson.Result{}, fmt.Errorf("could not find v parameter")
	}
	v := match[1]

	// Step 5: Generate rb parameter using md5 function
	rbByte := md5.Sum([]byte(fmt.Sprintf("%s%s%s%s", l.rid, l.did, l.t10, v)))
	rb := hex.EncodeToString(rbByte[:])
	// Step 6: Modify JavaScript function to replace return statement with rb parameter
	jsFunc = strings.Replace(result, "return rt;})", "return rt;}", -1)
	jsFunc = strings.Replace(jsFunc, "(function (", "function sign(", -1)
	jsFunc = strings.Replace(jsFunc, "CryptoJS.MD5(cb).toString()", fmt.Sprintf("\"%s\"", rb), -1)
	jsFunc = fmt.Sprintf("%s sign(\"%s\", \"%s\", \"%s\");", jsFunc, l.rid, l.did, l.t10)
	vmSignResult, err := vm.RunString(jsFunc)
	if err != nil {
		return gjson.Result{}, err
	}
	result = vmSignResult.String()
	params := fmt.Sprintf("%s&ver=22107261&rid=%s&rate=-1", result, l.rid)
	resp, err := http.Post("https://m.douyu.com/api/room/ratestream", "application/x-www-form-urlencoded", strings.NewReader(params))
	if err != nil {
		return gjson.Result{}, fmt.Errorf("error making POST request: %s", err)
	}
	defer resp.Body.Close()

	// Step 8: Parse the JSON response and return it
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return gjson.Result{}, err
	}
	return gjson.ParseBytes(body), nil
}
