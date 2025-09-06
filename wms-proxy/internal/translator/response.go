package translator

import (
	"fmt"
	"net/http"
)

// TranslateArcGISResponse handles the response from ArcGIS and prepares it for WMS client
func TranslateArcGISResponse(arcgisResp *http.Response, wmsWriter http.ResponseWriter) error {
	// Copy status code
	wmsWriter.WriteHeader(arcgisResp.StatusCode)

	// Copy relevant headers
	copyHeaders(arcgisResp.Header, wmsWriter.Header())

	// For successful image responses, copy the body directly
	if arcgisResp.StatusCode == http.StatusOK {
		// Ensure proper content type for images
		contentType := arcgisResp.Header.Get("Content-Type")
		if contentType == "" {
			// Default to PNG if no content type specified
			wmsWriter.Header().Set("Content-Type", "image/png")
		}

		// Copy response body
		_, err := copyResponseBody(arcgisResp, wmsWriter)
		return err
	}

	// For error responses, we might want to convert to WMS error format
	// For now, just pass through the ArcGIS error
	_, err := copyResponseBody(arcgisResp, wmsWriter)
	return err
}

// copyHeaders copies relevant headers from ArcGIS response to WMS response
func copyHeaders(src http.Header, dst http.Header) {
	// Headers to copy
	headersToCopy := []string{
		"Content-Type",
		"Content-Length",
		"Cache-Control",
		"Expires",
		"Last-Modified",
		"ETag",
	}

	for _, header := range headersToCopy {
		if value := src.Get(header); value != "" {
			dst.Set(header, value)
		}
	}
}

// copyResponseBody copies the response body from source to destination
func copyResponseBody(src *http.Response, dst http.ResponseWriter) (int64, error) {
	defer src.Body.Close()

	// Use io.Copy for efficient streaming
	written, err := copyBody(src.Body, dst)
	if err != nil {
		return written, fmt.Errorf("failed to copy response body: %w", err)
	}

	return written, nil
}

// copyBody is a simple implementation of io.Copy for response body
func copyBody(src interface{ Read([]byte) (int, error) }, dst interface{ Write([]byte) (int, error) }) (int64, error) {
	buf := make([]byte, 32*1024) // 32KB buffer
	var written int64

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = fmt.Errorf("invalid write result")
				}
			}
			written += int64(nw)
			if ew != nil {
				return written, ew
			}
			if nr != nw {
				return written, fmt.Errorf("short write")
			}
		}
		if er != nil {
			if er.Error() == "EOF" {
				break
			}
			return written, er
		}
	}
	return written, nil
}

// GenerateWMSError creates a WMS-compliant error response
func GenerateWMSError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/vnd.ogc.se_xml")
	w.WriteHeader(code)

	errorXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<ServiceExceptionReport version="1.1.1">
  <ServiceException>%s</ServiceException>
</ServiceExceptionReport>`, escapeXML(message))

	w.Write([]byte(errorXML))
}

// escapeXML performs basic XML escaping
func escapeXML(s string) string {
	s = replaceString(s, "&", "&amp;")
	s = replaceString(s, "<", "&lt;")
	s = replaceString(s, ">", "&gt;")
	s = replaceString(s, "\"", "&quot;")
	s = replaceString(s, "'", "&#39;")
	return s
}

// replaceString is a simple string replacement function
func replaceString(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}
