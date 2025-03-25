package metatags

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(MetaTags{})
	httpcaddyfile.RegisterHandlerDirective("metatags", func(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
		m := new(MetaTags)
		err := m.UnmarshalCaddyfile(h.Dispenser)
		return m, err
	})
}

// MetaTags injects meta tags based on the request path
type MetaTags struct {
	Routes []Route `json:"routes,omitempty"`
	logger *zap.Logger
}

// Route represents a URL pattern and its associated meta tags
type Route struct {
	Pattern  string   `json:"pattern,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
}

// Metadata matches Next.js metadata structure
type Metadata struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Keywords    []string  `json:"keywords,omitempty"`
	OpenGraph   OpenGraph `json:"openGraph,omitempty"`
	Twitter     Twitter   `json:"twitter,omitempty"`
}

// OpenGraph metadata
type OpenGraph struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Images      []string `json:"images,omitempty"`
	Type        string   `json:"type,omitempty"`
	SiteName    string   `json:"siteName,omitempty"`
	Locale      string   `json:"locale,omitempty"`
}

// Twitter metadata
type Twitter struct {
	Card        string   `json:"card,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Images      []string `json:"images,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (MetaTags) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.metatags",
		New: func() caddy.Module { return new(MetaTags) },
	}
}

// Provision sets up the middleware
func (m *MetaTags) Provision(ctx caddy.Context) error {
	m.logger = ctx.Logger(m)
	m.logger.Info("metatags middleware provisioned",
		zap.Int("routes", len(m.Routes)))
	return nil
}

// Validate ensures the middleware is properly configured
func (m *MetaTags) Validate() error {
	return nil
}

// ServeHTTP implements the HTTP handler interface
func (m *MetaTags) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	ctx := r.Context()
	repl, ok := ctx.Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	if !ok {
		return next.ServeHTTP(w, r)
	}
	path := repl.ReplaceAll("{http.request.orig_uri}", "") // https://github.com/caddyserver/caddy/blob/master/modules/caddyhttp/app.go#L65

	if r.URL.Path != "/index.html" {
		return next.ServeHTTP(w, r)
	}

	// https://github.com/caddyserver/caddy/blob/86c620fb4e7bfad5888832c491147af53fd5390a/modules/caddyhttp/templates/templates.go#L393C1-L395C24
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()            // Clear any previous content
	defer bufPool.Put(buf) // Return to pool when done

	// Proper usage: https://github.com/caddyserver/caddy/blob/86c620fb4e7bfad5888832c491147af53fd5390a/modules/caddyhttp/responsewriter.go#L104-L112
	wrec := caddyhttp.NewResponseRecorder(w, buf, func(status int, header http.Header) bool { return strings.Contains(header.Get("Content-Type"), "html") })
	err := next.ServeHTTP(wrec, r)
	if err != nil {
		return err
	}
	if !wrec.Buffered() {
		return nil
	}

	body := buf.Bytes()

	// Look for head tag
	headEndPos := bytes.Index(body, []byte("</head>")) //TODO: make a constant position
	if headEndPos == -1 {
		m.logger.Warn("no </head> tag found", zap.String("path", path))
		return wrec.WriteResponse()
	}

	metaTags := generateMetaTags(m.Routes, path)
	if metaTags == "" {
		m.logger.Info("no meta tags generated", zap.String("path", path))
		return wrec.WriteResponse()
	}

	// Safely create the new body with the meta tags
	newBody := make([]byte, 0, len(body)+len(metaTags))
	newBody = append(newBody, body[:headEndPos]...)
	newBody = append(newBody, []byte(metaTags)...)
	newBody = append(newBody, body[headEndPos:]...)

	// Update the response
	wrec.Header().Set("Content-Length", fmt.Sprintf("%d", len(newBody)))
	wrec.Buffer().Reset()
	_, err = wrec.Buffer().Write(newBody)
	if err != nil {
		m.logger.Error("failed to write modified body", zap.Error(err))
		return err
	}

	return wrec.WriteResponse()
}

// func (m *MetaTags) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
// 	ctx := r.Context()

// 	if repl, ok := ctx.Value(caddy.ReplacerCtxKey).(*caddy.Replacer); ok {
// 		pathVars := map[string]string{
// 			"path":      repl.ReplaceAll("{http.request.uri.path}", ""),
// 			"uri":       repl.ReplaceAll("{http.request.uri}", ""),
// 			"orig_path": repl.ReplaceAll("{http.request.orig_uri.path}", ""),
// 			"orig_uri":  repl.ReplaceAll("{http.request.orig_uri}", ""),
// 		}
// 		m.logger.Info("passthrough mode - not modifying response", zap.Any("path_vars", pathVars))
// 	}

// 	return next.ServeHTTP(w, r)
// }

// converts metadata to HTML meta tags
func generateMetaTags(routes []Route, path string) string {
	var metadata *Metadata
	for _, route := range routes {
		matched, _ := regexp.MatchString(route.Pattern, path)
		if matched {
			metadata = &route.Metadata
			break
		}
	}

	if metadata == nil {
		return ""
	}

	var sb strings.Builder

	// Basic metadata
	if metadata.Title != "" {
		sb.WriteString(fmt.Sprintf("<title>%s</title>\n", metadata.Title)) // TODO: figure out what to do here
		sb.WriteString(fmt.Sprintf("<meta name=\"title\" content=\"%s\" />\n", metadata.Title))
	}
	if metadata.Description != "" {
		sb.WriteString(fmt.Sprintf("<meta name=\"description\" content=\"%s\" />\n", metadata.Description))
	}
	if len(metadata.Keywords) > 0 {
		sb.WriteString(fmt.Sprintf("<meta name=\"keywords\" content=\"%s\" />\n", strings.Join(metadata.Keywords, ", ")))
	}

	// OpenGraph metadata
	if metadata.OpenGraph.Title != "" {
		sb.WriteString(fmt.Sprintf("<meta property=\"og:title\" content=\"%s\" />\n", metadata.OpenGraph.Title))
	}
	if metadata.OpenGraph.Description != "" {
		sb.WriteString(fmt.Sprintf("<meta property=\"og:description\" content=\"%s\" />\n", metadata.OpenGraph.Description))
	}
	if len(metadata.OpenGraph.Images) > 0 {
		for _, img := range metadata.OpenGraph.Images {
			// Process images with path variables
			processedImg := processPathVariables(img, path)
			sb.WriteString(fmt.Sprintf("<meta property=\"og:image\" content=\"%s\" />\n", processedImg))
		}
	}
	if metadata.OpenGraph.Type != "" {
		sb.WriteString(fmt.Sprintf("<meta property=\"og:type\" content=\"%s\" />\n", metadata.OpenGraph.Type))
	}
	if metadata.OpenGraph.SiteName != "" {
		sb.WriteString(fmt.Sprintf("<meta property=\"og:site_name\" content=\"%s\" />\n", metadata.OpenGraph.SiteName))
	}
	if metadata.OpenGraph.Locale != "" {
		sb.WriteString(fmt.Sprintf("<meta property=\"og:locale\" content=\"%s\" />\n", metadata.OpenGraph.Locale))
	}

	// Twitter metadata
	if metadata.Twitter.Card != "" {
		sb.WriteString(fmt.Sprintf("<meta name=\"twitter:card\" content=\"%s\" />\n", metadata.Twitter.Card))
	}
	if metadata.Twitter.Title != "" {
		sb.WriteString(fmt.Sprintf("<meta name=\"twitter:title\" content=\"%s\" />\n", metadata.Twitter.Title))
	}
	if metadata.Twitter.Description != "" {
		sb.WriteString(fmt.Sprintf("<meta name=\"twitter:description\" content=\"%s\" />\n", metadata.Twitter.Description))
	}
	if len(metadata.Twitter.Images) > 0 {
		for _, img := range metadata.Twitter.Images {
			// Process images with path variables
			processedImg := processPathVariables(img, path)
			sb.WriteString(fmt.Sprintf("<meta name=\"twitter:image\" content=\"%s\" />\n", processedImg))
		}
	}

	return sb.String()
}

// replaces path variables in content strings
func processPathVariables(content string, path string) string {
	// Extract path segments
	segments := strings.Split(strings.Trim(path, "/"), "/")

	// Replace {segment:n} with the corresponding path segment
	re := regexp.MustCompile(`\{segment:(\d+)\}`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		indexStr := re.FindStringSubmatch(match)[1]
		index := 0
		fmt.Sscanf(indexStr, "%d", &index)

		if index < len(segments) {
			return segments[index]
		}
		return ""
	})
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler
func (m *MetaTags) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.NextBlock(0) {
			switch d.Val() {
			case "route":
				if !d.NextArg() {
					return d.ArgErr()
				}

				route := Route{
					Pattern: d.Val(),
					Metadata: Metadata{
						OpenGraph: OpenGraph{},
						Twitter:   Twitter{},
					},
				}

				for d.NextBlock(1) {
					key := d.Val()

					switch key {
					case "title":
						if !d.NextArg() {
							return d.ArgErr()
						}
						route.Metadata.Title = d.Val()

					case "description":
						if !d.NextArg() {
							return d.ArgErr()
						}
						route.Metadata.Description = d.Val()

					case "keywords":
						var keywords []string
						for d.NextArg() {
							keywords = append(keywords, d.Val())
						}
						route.Metadata.Keywords = keywords

					case "openGraph":
						for d.NextBlock(2) {
							ogKey := d.Val()

							switch ogKey {
							case "title":
								if !d.NextArg() {
									return d.ArgErr()
								}
								route.Metadata.OpenGraph.Title = d.Val()

							case "description":
								if !d.NextArg() {
									return d.ArgErr()
								}
								route.Metadata.OpenGraph.Description = d.Val()

							case "images":
								var images []string
								for d.NextArg() {
									images = append(images, d.Val())
								}
								route.Metadata.OpenGraph.Images = images

							case "type":
								if !d.NextArg() {
									return d.ArgErr()
								}
								route.Metadata.OpenGraph.Type = d.Val()

							case "siteName":
								if !d.NextArg() {
									return d.ArgErr()
								}
								route.Metadata.OpenGraph.SiteName = d.Val()

							case "locale":
								if !d.NextArg() {
									return d.ArgErr()
								}
								route.Metadata.OpenGraph.Locale = d.Val()
							}
						}

					case "twitter":
						for d.NextBlock(2) {
							twitterKey := d.Val()

							switch twitterKey {
							case "card":
								if !d.NextArg() {
									return d.ArgErr()
								}
								route.Metadata.Twitter.Card = d.Val()

							case "title":
								if !d.NextArg() {
									return d.ArgErr()
								}
								route.Metadata.Twitter.Title = d.Val()

							case "description":
								if !d.NextArg() {
									return d.ArgErr()
								}
								route.Metadata.Twitter.Description = d.Val()

							case "images":
								var images []string
								for d.NextArg() {
									images = append(images, d.Val())
								}
								route.Metadata.Twitter.Images = images
							}
						}
					}
				}

				m.Routes = append(m.Routes, route)
			}
		}
	}

	return nil
}

// Interface guards
var (
	_ caddy.Provisioner           = (*MetaTags)(nil)
	_ caddy.Validator             = (*MetaTags)(nil)
	_ caddyhttp.MiddlewareHandler = (*MetaTags)(nil)
	_ caddyfile.Unmarshaler       = (*MetaTags)(nil)
)

var bufPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}
