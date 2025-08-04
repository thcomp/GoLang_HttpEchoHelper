package HttpEchoHelper

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/thcomp/GoLang_HttpEntityHelper/entity"
)

type EchoHelperContext struct {
	ctx       echo.Context
	helper    *EchoHelper
	entityIns entity.HttpEntity
}

func NewEchoHelperContext(ctx echo.Context, helper *EchoHelper) *EchoHelperContext {
	return &EchoHelperContext{
		ctx:    ctx,
		helper: helper,
	}
}

func (e *EchoHelperContext) Request() *http.Request {
	return e.ctx.Request()
}

func (e *EchoHelperContext) SetRequest(r *http.Request) {
	e.ctx.SetRequest(r)
}

func (e *EchoHelperContext) SetResponse(r *echo.Response) {
	e.ctx.SetResponse(r)
}

func (e *EchoHelperContext) Response() *echo.Response {
	return e.ctx.Response()
}

func (e *EchoHelperContext) IsTLS() bool {
	return e.ctx.IsTLS()
}

func (e *EchoHelperContext) IsWebSocket() bool {
	return e.ctx.IsWebSocket()
}

func (e *EchoHelperContext) Scheme() string {
	return e.ctx.Scheme()
}

func (e *EchoHelperContext) RealIP() string {
	return e.ctx.RealIP()
}

func (e *EchoHelperContext) Path() string {
	return e.ctx.Path()
}

func (e *EchoHelperContext) SetPath(p string) {
	e.ctx.SetPath(p)
}

func (e *EchoHelperContext) Param(name string) string {
	return e.ctx.Param(name)
}

func (e *EchoHelperContext) ParamNames() []string {
	return e.ctx.ParamNames()
}

func (e *EchoHelperContext) SetParamNames(names ...string) {
	e.ctx.SetParamNames(names...)
}

func (e *EchoHelperContext) ParamValues() []string {
	return e.ctx.ParamValues()
}

func (e *EchoHelperContext) SetParamValues(values ...string) {
	e.ctx.SetParamValues(values...)
}

func (e *EchoHelperContext) QueryParam(name string) string {
	return e.ctx.QueryParam(name)
}

func (e *EchoHelperContext) QueryParams() url.Values {
	return e.ctx.QueryParams()
}

func (e *EchoHelperContext) QueryString() string {
	return e.ctx.QueryString()
}

func (e *EchoHelperContext) FormValue(name string) string {
	return e.ctx.FormValue(name)
}

func (e *EchoHelperContext) FormParams() (url.Values, error) {
	return e.ctx.FormParams()
}

func (e *EchoHelperContext) FormFile(name string) (*multipart.FileHeader, error) {
	return e.ctx.FormFile(name)
}

func (e *EchoHelperContext) MultipartForm() (*multipart.Form, error) {
	return e.ctx.MultipartForm()
}

func (e *EchoHelperContext) Cookie(name string) (*http.Cookie, error) {
	return e.ctx.Cookie(name)
}

func (e *EchoHelperContext) SetCookie(cookie *http.Cookie) {
	e.ctx.SetCookie(cookie)
}

func (e *EchoHelperContext) Cookies() []*http.Cookie {
	return e.ctx.Cookies()
}

func (e *EchoHelperContext) Get(key string) interface{} {
	return e.ctx.Get(key)
}

func (e *EchoHelperContext) Set(key string, val interface{}) {
	e.ctx.Set(key, val)
}

func (e *EchoHelperContext) Bind(i interface{}) error {
	return e.ctx.Bind(i)
}

func (e *EchoHelperContext) Validate(i interface{}) error {
	return e.ctx.Validate(i)
}

func (e *EchoHelperContext) Render(code int, name string, data interface{}) error {
	return e.ctx.Render(code, name, data)
}

func (e *EchoHelperContext) HTML(code int, html string) error {
	return e.ctx.HTML(code, html)
}

func (e *EchoHelperContext) HTMLBlob(code int, b []byte) error {
	return e.ctx.HTMLBlob(code, b)
}

func (e *EchoHelperContext) String(code int, s string) error {
	return e.ctx.String(code, s)
}

func (e *EchoHelperContext) JSON(code int, i interface{}) error {
	return e.ctx.JSON(code, i)
}

func (e *EchoHelperContext) JSONPretty(code int, i interface{}, indent string) error {
	return e.ctx.JSONPretty(code, i, indent)
}

func (e *EchoHelperContext) JSONBlob(code int, b []byte) error {
	return e.ctx.JSONBlob(code, b)
}

func (e *EchoHelperContext) JSONP(code int, callback string, i interface{}) error {
	return e.ctx.JSONP(code, callback, i)
}

func (e *EchoHelperContext) JSONPBlob(code int, callback string, b []byte) error {
	return e.ctx.JSONPBlob(code, callback, b)
}

func (e *EchoHelperContext) XML(code int, i interface{}) error {
	return e.ctx.XML(code, i)
}

func (e *EchoHelperContext) XMLPretty(code int, i interface{}, indent string) error {
	return e.ctx.XMLPretty(code, i, indent)
}

func (e *EchoHelperContext) XMLBlob(code int, b []byte) error {
	return e.ctx.XMLBlob(code, b)
}

func (e *EchoHelperContext) Blob(code int, contentType string, b []byte) error {
	return e.ctx.Blob(code, contentType, b)
}

func (e *EchoHelperContext) Stream(code int, contentType string, r io.Reader) error {
	return e.ctx.Stream(code, contentType, r)
}

func (e *EchoHelperContext) File(file string) error {
	return e.ctx.File(file)
}

func (e *EchoHelperContext) Attachment(file string, name string) error {
	return e.ctx.Attachment(file, name)
}

func (e *EchoHelperContext) Inline(file string, name string) error {
	return e.ctx.Inline(file, name)
}

func (e *EchoHelperContext) NoContent(code int) error {
	return e.ctx.NoContent(code)
}

func (e *EchoHelperContext) Redirect(code int, url string) error {
	return e.ctx.Redirect(code, url)
}

func (e *EchoHelperContext) Error(err error) {
	e.ctx.Error(err)
}

func (e *EchoHelperContext) Handler() echo.HandlerFunc {
	return e.ctx.Handler()
}

func (e *EchoHelperContext) SetHandler(h echo.HandlerFunc) {
	e.ctx.SetHandler(h)
}

func (e *EchoHelperContext) Logger() echo.Logger {
	return e.ctx.Logger()
}

func (e *EchoHelperContext) SetLogger(l echo.Logger) {
	e.ctx.SetLogger(l)
}

func (e *EchoHelperContext) Echo() *echo.Echo {
	return e.ctx.Echo()
}

func (e *EchoHelperContext) Reset(r *http.Request, w http.ResponseWriter) {
	e.ctx.Reset(r, w)
}

func (e *EchoHelperContext) EchoHelper(inf ...*EchoHelper) *EchoHelper {
	if len(inf) > 0 && inf[0] != nil {
		e.helper = inf[0]
	}

	return e.helper
}

func (e *EchoHelperContext) HttpEntity(entityIns ...entity.HttpEntity) entity.HttpEntity {
	if len(entityIns) > 0 && entityIns[0] != nil {
		e.entityIns = entityIns[0]
	}

	return e.entityIns
}
