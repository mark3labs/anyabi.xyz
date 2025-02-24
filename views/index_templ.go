// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "fmt"

func Index(chainIDs map[string]int64) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\" class=\"dark uk-theme-violet\"><head><meta charset=\"utf-8\"><link rel=\"icon\" href=\"/static/favicon.png\"><meta name=\"viewport\" content=\"width=device-width\"><title>AnyABI - The easiest way to fetch an ABI</title><meta name=\"twitter:title\" content=\"AnyABI - The easiest way to fetch an ABI\"><meta name=\"description\" content=\"Quickly grab the ABI for ANY EVM smart contract on ANY EVM chain.\"><meta name=\"twitter:description\" content=\"Quickly grab the ABI for ANY EVM smart contract on ANY EVM chain.\"><meta name=\"twitter:image\" content=\"https://anyabi.xyz/static/thumbnail.png\"><meta property=\"og:image\" content=\"https://anyabi.xyz/static/thumbnail.png\"><meta name=\"twitter:card\" content=\"summary_large_image\"><meta http-equiv=\"content-security-policy\" content=\"\"><link rel=\"stylesheet\" href=\"https://unpkg.com/franken-ui@2.0.0-internal.45/dist/css/core.min.css\"><link rel=\"stylesheet\" href=\"https://unpkg.com/franken-ui@2.0.0-internal.45/dist/css/utilities.min.css\"><script type=\"module\" src=\"https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.0-beta.7/bundles/datastar.js\"></script><script type=\"module\" src=\"https://unpkg.com/franken-ui@2.0.0-internal.45/dist/js/core.iife.js\"></script><script type=\"module\" src=\"https://unpkg.com/franken-ui@2.0.0-internal.45/dist/js/icon.iife.js\"></script></head><body class=\"bg-background text-foreground\"><nav class=\"w-full border-b border-border bg-background\"><div class=\"uk-container\"><div class=\"flex h-14 items-center\"><a href=\"/\" class=\"uk-h4 m-0 mr-4\"><img src=\"/static/logo.png\" class=\"h-10 inline-block mr-2\"> Any ABI</a></div></div></nav><main class=\"my-10\"><div id=\"abi-explorer\" class=\"flex items-center justify-center uk-container\"><div class=\"border border-border rounded-lg px-8 py-6 w-full h-auto\" data-signals=\"{address:&#39;&#39;, chainId:&#39;1&#39;}\"><h3 class=\"uk-h3 mb-6\">ABI Explorer</h3><div class=\"flex gap-4\"><input type=\"text\" class=\"uk-input flex-1\" placeholder=\"0x...\" data-bind-address> <uk-select id=\"chain-select\" cls-custom=\"button: uk-input-fake w-48; dropdown: w-48\" searchable reactive value=\"1\" data-on-uk-select:input__case.kebab=\"$chainId = evt.detail.value\"><select hidden>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for name, id := range chainIDs {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<option value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(id))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/index.templ`, Line: 59, Col: 33}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "\" data-keywords=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/index.templ`, Line: 60, Col: 31}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/index.templ`, Line: 62, Col: 17}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</select></uk-select><div class=\"flex items-center\"><button class=\"uk-btn uk-btn-primary\" data-attr-disabled=\"!$address || !$chainId\" data-on-click=\"@get(&#39;/get-abi&#39;)\" data-indicator-fetching data-show=\"!$fetching\">Get ABI</button><div data-uk-spinner data-show=\"$fetching\"></div></div></div><div id=\"result\"></div></div></div></main></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
