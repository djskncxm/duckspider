package httpio

import (
	"bytes"
	"github.com/andybalholm/cascadia"
	"github.com/antchfx/htmlquery"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"io"
	"mime"
	http_ "net/http"
	"regexp"
)

type Response struct {
	URL     string
	Headers http_.Header
	Body    io.Reader
	Request *Request
	Status  int
	Meta    map[string]interface{}
}

type Selection interface {
	XPath(expr string) Selection
	CSS(selector string) Selection
	Text() (s string, err error)
	JSON(format interface{}) error
	Regex(pattern string) []string
	Get() string
	GetAll() []string
}

type DataSelection struct {
	rootNode *html.Node
	curNodes []*html.Node
	Response Response
	text     *string
}

func BuildDataSelection(
	url string,
	headers http_.Header,
	status int,
	request *Request,
	body io.Reader,
) (*DataSelection, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		URL:     url,
		Headers: headers,
		Status:  status,
		Request: request,
		Body:    bytes.NewReader(bodyBytes),
		Meta:    request.Meat,
	}

	root, err := htmlquery.Parse(bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	text := string(bodyBytes)
	ds := &DataSelection{
		rootNode: root,
		Response: *resp,
		text:     &text,
	}

	return ds, nil
}

func (d *DataSelection) XPath(expr string) Selection {
	if d.curNodes == nil {
		d.curNodes = []*html.Node{d.rootNode}
	}
	var results []*html.Node
	for _, node := range d.curNodes {
		found := htmlquery.Find(node, expr)
		results = append(results, found...)
	}
	return &DataSelection{
		rootNode: d.rootNode,
		curNodes: results,
		Response: d.Response,
		text:     d.text,
	}
}

func (d *DataSelection) CSS(selector string) Selection {
	sel, err := cascadia.Compile(selector)
	if err != nil {
		panic("请传递一个正确的CSS选择器写法 " + err.Error())
	}
	if d.curNodes == nil {
		d.curNodes = []*html.Node{d.rootNode}
	}

	var results []*html.Node
	for _, node := range d.curNodes {
		found := cascadia.QueryAll(node, sel)
		results = append(results, found...)
	}
	return &DataSelection{
		rootNode: d.rootNode,
		curNodes: results,
		Response: d.Response,
		text:     d.text,
	}
}

func (d *DataSelection) Text() (s string, err error) {
	if d.text != nil {
		return *d.text, nil
	}
	contentType := d.Response.Headers.Get("Content-Type")
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		// 解析失败，默认为空 charset
		params = map[string]string{}
	}
	charsetStr := params["charset"]

	reader, err := charset.NewReader(d.Response.Body, charsetStr)
	if err != nil {
		return "", err
	}

	result, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (d *DataSelection) JSON(format interface{}) error {
	return json.Unmarshal([]byte(*d.text), format)
}

func (d *DataSelection) Regex(pattern string) []string {
	if d.text == nil {
		return nil
	}
	re := regexp.MustCompile(pattern)
	return re.FindAllString(*d.text, -1)
}

func (d *DataSelection) Get() string {
	if len(d.curNodes) == 0 {
		return ""
	}
	var buf bytes.Buffer
	err := html.Render(&buf, d.curNodes[0])
	if err != nil {
		return ""
	}
	return buf.String()
}

func (d *DataSelection) GetAll() []string {
	results := make([]string, len(d.curNodes))

	for i, node := range d.curNodes {
		var buf bytes.Buffer
		err := html.Render(&buf, node)
		if err != nil {
			return nil
		}
		results[i] = buf.String()
	}
	return results
}
