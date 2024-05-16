package main

import (
	"crypto"
	"crypto/sha256"
	"io"
	"log"
	"net/http"
	"time"

	"crypto/rand"

	"cloud.google.com/go/storage"
	salkmssign "github.com/salrashid123/signer/kms"
)

var (
	projectId           = "your-project"
	bucketName          = "your-bucket"
	kmsKeyRing          = "mycacerts"
	kmsKey              = "key1"
	kmsKeyVersion       = "1"
	kmsLocationId       = "us-central1"
	serviceAccountEmail = "kms-svc-account@your-project.iam.gserviceaccount.com"
	keyId               = "ce4ceffd5f9c8b399df9bf7b5c13327dab65f180"
)

func main() {

	r, err := salkmssign.NewKMSCrypto(&salkmssign.KMS{
		ProjectId:  projectId,
		LocationId: kmsLocationId,
		KeyRing:    kmsKeyRing,
		Key:        kmsKey,
		KeyVersion: kmsKeyVersion,
	})

	if err != nil {
		log.Fatal(err)
	}

	bucket := bucketName
	object := "foo.txt"
	keyID := serviceAccountEmail
	expires := time.Now().Add(time.Minute * 10)

	s, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		GoogleAccessID: keyID,
		SignBytes: func(b []byte) ([]byte, error) {
			sum := sha256.Sum256(b)
			return r.Sign(rand.Reader, sum[:], crypto.SHA256)
		},
		Method:  "GET",
		Expires: expires,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(s)

	resp, err := http.Get(s)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	log.Println("SignedURL Response :\n", string(body))
	if err != nil {
		log.Fatal(err)
	}

	// =============================================

}
