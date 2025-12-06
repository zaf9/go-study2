package handler

const htmlStyle = `
<style>
body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: #1e1e1e; color: #d4d4d4; padding: 20px; max-width: 900px; margin: 0 auto; }
h1 { color: #569cd6; border-bottom: 1px solid #333; padding-bottom: 10px; }
a { color: #9cdcfe; text-decoration: none; transition: color 0.3s; }
a:hover { color: #4ec9b0; }
ul { list-style-type: none; padding: 0; }
li { background: #252526; margin: 8px 0; padding: 12px; border-radius: 4px; border-left: 4px solid #569cd6; }
li:hover { background: #2d2d30; }
pre { font-family: 'Consolas', 'Courier New', monospace; background: #2d2d30; padding: 15px; border-radius: 4px; overflow-x: auto; font-size: 14px; line-height: 1.5; color: #ce9178; box-shadow: 0 4px 6px rgba(0,0,0,0.3); }
.back-link { display: inline-block; margin-top: 20px; padding: 10px 15px; background: #0e639c; color: white; border-radius: 4px; }
.back-link:hover { background: #1177bb; color: white; }
</style>
`

// getHtmlPage 生成带有样式的 HTML 页面
func getHtmlPage(title, content string) string {
	return "<!DOCTYPE html>\n<html>\n<head><title>" + title + "</title>" + htmlStyle + "</head>\n<body>\n" +
		content +
		"\n</body>\n</html>"
}
