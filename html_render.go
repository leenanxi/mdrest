package mdrest

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type HTMLRenderer struct {
	basePath string
	location string
	blackfriday.Renderer
}

func (renderer *HTMLRenderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	if bytes.HasPrefix(link, []byte(".")) || !bytes.Contains(link, []byte("://")) {
		link = []byte(strings.Replace(AbsPath(renderer.basePath, renderer.location, string(link)), " ", "%20", -1))
	}
	renderer.Renderer.Link(out, link, title, content)
}

func (renderer *HTMLRenderer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	if !(bytes.HasPrefix(link, []byte("http://")) || bytes.HasPrefix(link, []byte("https://"))) {
		link = []byte(AbsPath(renderer.basePath, renderer.location, string(link)))
	}
	out.WriteString(`<div class="image-package"><img src="`)
	out.Write(link)
	out.WriteString(`" alt="`)
	out.Write(alt)
	out.WriteString(`"`)
	if title != nil {
		out.WriteString(` title="`)
		out.Write(title)
		out.WriteString(`"/>`)
		out.WriteString(`<div class="caption">`)
		out.Write(title)
		out.WriteString(`</div>`)
	} else {
		out.WriteString(`/>`)
	}
	out.WriteString(`</div>`)
}

/***
<div class="tabs">
	<input id="515736620" type="radio" name="tab1" track-name="golang" checked="checked"/>
	<label for="515736620">Golang</label>
	<section>1</section>
	<input id="479272831" type="radio" track-name="bash" name="tab1"/>
	<label for="479272831">Bash</label>
	<section>2</section>
	<input id="479272834" type="radio" track-name="test" name="tab1"/>
	<label for="479272834">Test</label>
	<section>3</section>
</div>
*/


/***
<div class="tabs">
	<input id="515736620" type="radio" name="tab1" track-name="golang" checked="checked"/>
	<label for="515736620">Golang</label>
	<section>1</section>
	<input id="479272831" type="radio" track-name="bash" name="tab1"/>
	<label for="479272831">Bash</label>
	<section>2</section>
	<input id="479272834" type="radio" track-name="test" name="tab1"/>
	<label for="479272834">Test</label>
	<section>3</section>
</div>
*/


// ListItem adds task list support to the Blackfriday renderer.
func (renderer *HTMLRenderer) ListItem(out *bytes.Buffer, text []byte, flags int) {
	//if string(text) == "&lt;--" {
	//	text = []byte(`<div class="tabs">`)
	//	renderer.Renderer.ListItem(out, text, flags)
	//	return
	//}
	//if string(text) == "--&gt;" {
	//	text = []byte( `</div>`)
	//	fmt.Println("哈哈", flags)
	//
	//	fmt.Println(out.String())
	//	renderer.Renderer.ListItem(out, text, flags)
	//	return
	//}

	t := strconv.FormatInt(time.Now().UnixNano() / 10, 10)
	id := t[len(t) - 10:]
	str := string(text)
	if strings.HasPrefix(str, "&gt; ") {
		tpl := `<input track-name="%s" id="%s" type="radio" name="%s"/><label for="%s">%s</label>`
		label := strings.TrimLeft(str, "&gt; ")
		tpl = fmt.Sprintf(tpl, id, "tab", url.QueryEscape(strings.ToLower(label)), id, label)
		text = []byte(tpl)
		renderer.Renderer.ListItem(out, text, flags)
	}

	fmt.Println(string("--------"))

	if bytes.HasPrefix(text, []byte(`<p>[ ] `)) {
		right := bytes.Index(text, []byte("</p>"))
		text = []byte(`[ ] ` + string(text[7:right]) + string(text[right+4:]))
	}

	switch {
	case bytes.HasPrefix(text, []byte("[ ] ")):
		text = append([]byte(`<input type="checkbox" disabled />`), text[3:]...)

	case bytes.HasPrefix(text, []byte("[x] ")) || bytes.HasPrefix(text, []byte("[X] ")):
		text = append([]byte(`<input type="checkbox" checked disabled />`), text[3:]...)
	}
	renderer.Renderer.ListItem(out, text, flags)
}
