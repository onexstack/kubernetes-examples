package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
)

type CustomResource struct {
	Message string `json:"message"`
}

// 嵌入的证书和私钥
const serverCert = `-----BEGIN CERTIFICATE-----
MIIFazCCA1OgAwIBAgIUUSWa2AJ2fboO6VNYOhGadPCP9DAwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yNDA5MDUxNjU0MjFaFw0yNTA5
MDUxNjU0MjFaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggIiMA0GCSqGSIb3DQEB
AQUAA4ICDwAwggIKAoICAQC0zMGYQBUTxT96Wja4wPr1b7hjeCXuo5NyvGDrQwK4
/fF5EWzXYvt/QhE593Ni8C0HNGvEs1g2TPAh7aOGaMx9pA3WYEA/iIY5xM3NC0H1
Q+wNnYKEtRdV6ioFdHrks0QP+ohQMUVNfgIonOer3JK+KIR5nDzMBaMbwINgs5Ty
C74oQOLfB/DCcha4MbKABV/Xt6667gIbIKirUa56iLravUdO0XiV7lsQK+HVQEAH
1rw0AAg/R61ElK41dYqERqljX/lnch5Tk71HM9wG9gh5wkf6/F7CU/CktuT6XepY
1tQs3m/UErRCnVJWd2AOY4hVmhmDyU6iCy2/cmZj86eUhOhXlPHwCZHRw8i6XnqM
S6FVUFyc9xb2zbI1zoRIgZjpZpFDc1pimm+6hWZBUrvtmxT107EBH1AWsICbCpVG
Rfw7aVLfPLZ42/98Sd6PZVMWHvV/Inkaw+s0lA41x7b6P9m+xHPB9Q+WqXvgBIPx
Bh8K2WGlaBh7+a6V/Yf0S+Rs7TfZCHDksxCSIzjKaT6l2/KsoMuO7d7cBwddyMCz
B8CTNT6glgr35QLkrBLLetYwTDQTk/W3IE0O3nxr+t8ZwHbhIbISO8hlHvZ75ODi
AqveMCq5AIybNyjjy4MtHhYyCi+fh2mxBsQ4r0O1kFy3DZ35pUJN6Av9Rcs4WGIY
IQIDAQABo1MwUTAdBgNVHQ4EFgQUEqBH4cytt+3xvxzt4p37MORBuPEwHwYDVR0j
BBgwFoAUEqBH4cytt+3xvxzt4p37MORBuPEwDwYDVR0TAQH/BAUwAwEB/zANBgkq
hkiG9w0BAQsFAAOCAgEAQJGfYzSnD7MODdcSmCMkbqaiF9/5RKd9LqtEpzc6VsXz
GIoaMHLRvhcTgOjBjsDUR+NKWTGUXXWX95ekHK7NsxWxywNdftHBHbfrFxf084oq
6MwmX4l8BJeNFzmuhdQ20irl1Jg4o8trYbLgQRcVLewk9CpXONWueAdssx7vvpce
6GA+p/iJjouQQOlX04PxN9kBu6iSwTt/XoRp2jFayGA5pNBMOiTrQ+3h7qehAeFF
UMgnS0gkAbYd2cNfAzJUyVw3N1YM5eifBXv2wiunbsKL0IRumaJQx0JyDzux5c/A
ipbeoLpUvWtXyJY42pmz03fNtlL9Qi5129DeAN+lITLzq4Heh/9bI3YKNrbtQUhC
8lXilJ0/dxg7b1cdHhp7wGLvHtIn2UWh8jl5ZsYWf7CpY1bHIP8gwgPacazdP5Uv
hrSqlMwltlWz/PcZZcdLbTh3jIANIBU58xcA1Z0Jtjq8OGw351KUN1+NmLYZLxlT
GPiSA+t722o2A06SZkpaWTVGOg+7roD++nRYLj+4zdyoe7IIILAHe51TdlJsf06E
Eg1fzGi1JOx9rwbOX1DENAGX69b7dWZzZYoTWx6qKmAyu1Hvj0bpYdxcrnG9f+AW
JfklZBVoV6Z8gFgAEB1gxoliDs5kLJjcz0UysSiw4LViRGxeFn+tBDfilRTQ6Ao=
-----END CERTIFICATE-----`
const serverKey = `-----BEGIN PRIVATE KEY-----
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
WbBOn53+DyqVbmc8p4RNYKr/7Ltlt8R1
-----END PRIVATE KEY-----`

func main() {
	http.HandleFunc("/api/v1/customresource", func(w http.ResponseWriter, r *http.Request) {
		response := CustomResource{Message: "Hello from Custom Resource API!"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// 创建一个自定义的 TLS 配置
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13, // 设置最小 TLS 版本
	}

	// 加载嵌入的证书和私钥
	cert, err := tls.X509KeyPair([]byte(serverCert), []byte(serverKey))
	if err != nil {
		log.Fatalf("Failed to parse certificate and key: %v", err)
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// 创建 HTTPS 服务器
	srv := &http.Server{
		Addr:      ":8443", // HTTPS 端口
		Handler:   nil,
		TLSConfig: tlsConfig,
	}

	// 启动 HTTPS 服务器
	log.Println("Starting HTTPS server on port 8443...")
	if err := srv.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("ListenAndServeTLS failed: %v", err)
	}
}
