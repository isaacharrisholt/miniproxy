package proxy

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gobwas/glob"
)

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

type proxyRoute struct {
	source     string
	sourceGlob glob.Glob
	target     string
}

type proxy struct {
	routes        []proxyRoute
	targets       map[string]proxyTarget
	defaultTarget string
	port          int
	logger        Logger
}

func NewProxy(s Settings) (proxy, error) {
	p := proxy{
		logger: NewChannelLogger("miniproxy"),
	}
	for target, spec := range s.Targets {
		if spec.Port == 0 {
			return p, fmt.Errorf("no port provided for target: %s", target)
		}
	}
	for source, target := range s.Routes {
		if _, ok := s.Targets[target]; !ok {
			return p, fmt.Errorf("Target not found for route %s: %s", source, target)
		}

		route := proxyRoute{}
		g, err := glob.Compile(source)
		if err != nil {
			return p, err
		}
		route.target = target
		route.source = source
		route.sourceGlob = g
		p.routes = append(p.routes, route)
	}
	p.targets = s.Targets
	p.defaultTarget = s.Default
	if s.Port == 0 {
		p.port = proxyDefaultPort
	} else {
		p.port = s.Port
	}
	return p, nil
}

func (p proxy) handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	p.logger.info(fmt.Sprintf("request: %s %s", r.Method, path))
	p.logger.debug(fmt.Sprintf("p.defaultTarget: %v", p.defaultTarget))

	var targetPort int
	for _, route := range p.routes {
		p.logger.debug(fmt.Sprintf("checking route: %s", route.source))
		if route.sourceGlob.Match(path) {
			target, ok := p.targets[route.target]
			if !ok {
				p.logger.error(fmt.Sprintf("no target found for: %s", route.target))
				respond(w, "Not found", http.StatusNotFound)
				return
			}
			p.logger.debug(fmt.Sprintf("route found: %s -> %d", route.source, target.Port))
			targetPort = target.Port
			break
		}
	}
	if targetPort == 0 {
		if p.defaultTarget == "" {
			p.logger.error(fmt.Sprintf("no target found for: %s", path))
			respond(w, "Not found", http.StatusNotFound)
			return
		}
		target, ok := p.targets[p.defaultTarget]
		if !ok {
			p.logger.error(fmt.Sprintf("no target found for: %s", p.defaultTarget))
			respond(w, "Not found", http.StatusNotFound)
			return
		}
		p.logger.debug(fmt.Sprintf("default target found: %s -> %d", p.defaultTarget, target.Port))
		targetPort = target.Port
	}
	targetUrl := fmt.Sprintf(
		"http://localhost:%d%s",
		targetPort,
		path,
	)
	if r.URL.RawQuery != "" {
		targetUrl += "?" + r.URL.RawQuery
	}
	targetReq, err := http.NewRequest(r.Method, targetUrl, r.Body)
	if err != nil {
		p.logger.error(fmt.Sprintf("error creating request: %s", err))
		respond(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	copyHeaders(targetReq.Header, r.Header)
	targetRes, err := http.DefaultClient.Do(targetReq)
	if err != nil {
		p.logger.error(fmt.Sprintf("error sending request: %s", err))
		respond(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	copyHeaders(w.Header(), targetRes.Header)
	w.WriteHeader(targetRes.StatusCode)
	_, err = io.Copy(w, targetRes.Body)
	if err != nil {
		p.logger.error(fmt.Sprintf("error copying response: %s", err))
		respond(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (p proxy) startServices() error {
	go stdOutReceiver()
	for targetName, target := range p.targets {
		if target.Service.Command != nil {
			startServiceWithChannelLogger(target.Service, targetName)
		}
	}
	return nil
}

func (p proxy) Start() error {
	err := p.startServices()
	if err != nil {
		return err
	}
	http.HandleFunc("/", p.handler)
	p.logger.info(fmt.Sprintf("listening on port %d", p.port))
	return http.ListenAndServe(fmt.Sprintf(":%d", p.port), nil)
}

func respond(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}
