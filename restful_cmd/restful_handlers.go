package main

import (
	"aapi/config"
	gw "aapi/services/user/protos/user/v1"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
}

type proxyHandler struct {
	userApiClient gw.UserClient
	cfg           *config.Config
}

type ProxyHandler interface {
	Enroll(w http.ResponseWriter, r *http.Request)
	Verify(w http.ResponseWriter, r *http.Request)
	Ping(w http.ResponseWriter, r *http.Request)
}

func NewProxyHandler(
	userCnn *grpc.ClientConn,
	cfg *config.Config,
) ProxyHandler {
	return &proxyHandler{
		userApiClient: gw.NewUserClient(userCnn),
		cfg:           cfg,
	}
}

type Image struct {
	Data    string
	ImageId string
}

func (h *proxyHandler) Enroll(w http.ResponseWriter, r *http.Request) {
	log.Println("Enroll invoked")
	request, err := getEnrollRequest(w, r)
	if request == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	response, err := h.userApiClient.Enroll(r.Context(), request)
	log.Default().Printf("Enroll.Response :%v", response)
	log.Default().Printf("Enroll.Err :%v", err)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(response)
}

func (h *proxyHandler) Verify(w http.ResponseWriter, r *http.Request) {
	log.Println("Verify invoked")
	request := getVerifyRequest(w, r)
	if request == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response, err := h.userApiClient.Verify(r.Context(), request)
	log.Default().Printf("Verify.Response :%v", response)
	log.Default().Printf("Verify.Err :%v", err)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(response)
	}
}

func (h *proxyHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong\n"))
}

func getEnrollRequest(w http.ResponseWriter, r *http.Request) (*gw.EnrollRequest, error) {
	err := r.ParseMultipartForm(12 << 21)
	if err != nil {
		return nil, err
	}
	// Get handler for filename, size and headers
	formdata := r.MultipartForm
	request := &gw.EnrollRequest{}
	values := r.MultipartForm.Value
	log.Default().Printf("multiparform values : %v \n", values)
	images := []*gw.Image{}
	files := formdata.File["images"]
	for i := range files { // loop through the files one by one
		file, err := files[i].Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()
		log.Printf("processing file %v: %v\n", i, files[i].Filename)
		log.Printf("File Size: %v\n", files[i].Size)
		log.Printf("MIME Header: %v\n", files[i].Header)

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			log.Default().Printf("Got error when parse image to binary: %v", err)
			return nil, err
		}
		image := gw.Image{
			Data:    buf.Bytes(),
			ImageId: request.UserName,
		}

		images = append(images, &image)
	}

	for k, v := range values {
		if strings.EqualFold(k, "UserId") {
			request.UserId = v[0]
		}
		if strings.EqualFold(k, "ReferenceId") {
			request.ReferenceId = v[0]
		}
		if strings.EqualFold(k, "UserInfo") {
			request.UserInfo = v[0]
		}
		if strings.EqualFold(k, "UserName") {
			request.UserName = v[0]
		}
	}

	if request.UserId == "" {
		request.UserId = strconv.FormatInt(rand.Int63n(100000000), 10)
	}
	if request.ReferenceId == "" {
		request.ReferenceId = uuid.NewString()
	}
	request.UserRoleId = 1
	request.UserGroupIds = []int64{1}
	request.UserGroups = []string{uuid.NewString()}
	request.ExpiryDate = time.Now().AddDate(0, 0, 7*365).UnixMilli()
	request.ActivationDate = time.Now().UTC().UnixMilli()

	request.Images = images
	return request, nil
}

func getVerifyRequest(w http.ResponseWriter, r *http.Request) *gw.VerifyRequest {
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("images")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return nil
	}
	defer file.Close()
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)
	request := &gw.VerifyRequest{}
	values := r.MultipartForm.Value
	log.Default().Printf("multiparform values : %v \n", values)

	images := []*gw.Image{}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Default().Printf("Got error when parse image to binary: %v", err)
		return nil
	}

	image := gw.Image{
		Data:    buf.Bytes(),
		ImageId: uuid.NewString(),
	}

	images = append(images, &image)
	request.Images = images
	return request
}
